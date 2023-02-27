package sentinel

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"golang.org/x/time/rate"
)

const (
	QueryArkAuth = "arkauth"
)

// Create a map to hold the rate limiters for each visitor and a mutex.
var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

type ArkAuth struct {
	Provider  common.PubKey
	Chain     common.Chain
	Spender   common.PubKey
	Height    int64
	Nonce     int64
	Signature []byte
}

func parseArkAuth(raw string) (ArkAuth, error) {
	var aa ArkAuth
	var err error

	parts := strings.SplitN(raw, ":", 6)
	if len(parts) != 6 {
		return aa, fmt.Errorf("Not properly formatted ark-auth string: %s\n", raw)
	}
	aa.Provider, err = common.NewPubKey(parts[0])
	if err != nil {
		return aa, err
	}
	aa.Chain, err = common.NewChain(parts[1])
	if err != nil {
		return aa, err
	}
	aa.Spender, err = common.NewPubKey(parts[2])
	if err != nil {
		return aa, err
	}
	aa.Height, err = strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return aa, err
	}
	aa.Nonce, err = strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return aa, err
	}
	aa.Signature, err = hex.DecodeString(parts[5])
	if err != nil {
		return aa, err
	}
	return aa, nil
}

func (aa ArkAuth) Validate(provider common.PubKey) error {
	creator, err := provider.GetMyAddress()
	if err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}
	if !provider.Equals(aa.Provider) {
		return fmt.Errorf("provider pubkey does not match provider")
	}
	msg := types.NewMsgClaimContractIncome(creator.String(), aa.Provider, aa.Chain.String(), aa.Spender, aa.Nonce, aa.Height, aa.Signature)
	return msg.ValidateBasic()
}

func (p Proxy) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var aa ArkAuth

		args := r.URL.Query()
		raw, aaOK := args[QueryArkAuth]
		if aaOK {
			aa, err = parseArkAuth(raw[0])
		}
		remoteAddr := p.getRemoteAddr(r)
		if err != nil || aa.Validate(p.Config.ProviderPubKey) == nil {
			p.logger.Info("serving paid requests", "remote-addr", remoteAddr)
			httpCode, err := p.paidTier(aa, remoteAddr)
			if err != nil {
				p.logger.Error("failed to serve paid tier request", "error", err)
				http.Error(w, err.Error(), httpCode)
				return
			}
		} else {
			p.logger.Info("serving free tier requests", "remote-addr", remoteAddr)
			httpCode, err := p.freeTier(remoteAddr)
			if err != nil {
				p.logger.Error("failed to serve free tier request", "error", err)
				http.Error(w, err.Error(), httpCode)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

const (
	forwardHeaderName = `X-Forwarded-For`
	xRealIPName       = `X-Real-Ip`
)

func (p Proxy) getRemoteAddr(r *http.Request) string {
	realIP := r.Header.Get(xRealIPName)
	if realIP != "" {
		return realIP
	}
	forwardIP := r.Header.Get(forwardHeaderName)
	if forwardIP != "" {
		return forwardIP
	}
	return r.RemoteAddr
}

func (p Proxy) freeTier(remoteAddr string) (int, error) {
	// Get the IP address for the current user.
	key, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err)
	}

	if ok := p.isRateLimited(key, -1); ok {
		return http.StatusTooManyRequests, fmt.Errorf(http.StatusText(429))
	}

	return http.StatusOK, nil
}

func (p Proxy) isRateLimited(key string, contractType types.ContractType) bool {
	mu.Lock()
	defer mu.Unlock()

	var limitTokens int
	var limitDuration time.Duration
	switch contractType {
	case types.ContractType_Subscription:
		limitTokens = p.Config.SubTierRateLimit
		limitDuration = p.Config.SubTierRateLimitDuration
	case types.ContractType_PayAsYouGo:
		limitTokens = p.Config.AsGoTierRateLimit
		limitDuration = p.Config.AsGoTierRateLimitDuration
	default:
		limitTokens = p.Config.FreeTierRateLimit
		limitDuration = p.Config.FreeTierRateLimitDuration
	}

	limiter, exists := visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(limitDuration), limitTokens)
		visitors[key] = limiter
	}

	return !limiter.Allow()
}

func (p Proxy) paidTier(aa ArkAuth, remoteAddr string) (code int, err error) {
	key := fmt.Sprintf("%d-%s", aa.Chain, aa.Spender)
	contractKey := p.MemStore.Key(aa.Provider.String(), aa.Chain.String(), aa.Spender.String())
	contract, err := p.MemStore.Get(contractKey)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err)
	}

	if contract.IsClose(p.MemStore.GetHeight()) {
		return http.StatusPaymentRequired, fmt.Errorf("open a contract")
	}

	sig := hex.EncodeToString(aa.Signature)
	claim := NewClaim(aa.Provider, aa.Chain, aa.Spender, aa.Nonce, aa.Height, sig)
	if p.ClaimStore.Has(key) {
		var err error
		claim, err = p.ClaimStore.Get(key)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err)
		}
		if claim.Height == aa.Height && contract.Height == aa.Height && claim.Nonce >= aa.Nonce {
			return http.StatusBadRequest, fmt.Errorf("bad nonce (%d/%d)", aa.Nonce, claim.Nonce)
		}
	}

	// check if we've exceed the total number of pay-as-you-go queries
	if contract.Type == types.ContractType_PayAsYouGo {
		if contract.Deposit.IsNil() || contract.Deposit.LT(cosmos.NewInt(aa.Nonce*contract.Rate)) {
			return http.StatusPaymentRequired, fmt.Errorf("contract spent")
		}
	}

	if ok := p.isRateLimited(key, contract.Type); ok {
		return http.StatusTooManyRequests, fmt.Errorf("client is ratelimited," + http.StatusText(429))
	}

	claim.Nonce = aa.Nonce
	claim.Height = aa.Height
	claim.Signature = sig
	claim.Claimed = false
	if err := p.ClaimStore.Set(claim); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err)
	}
	contract.Nonce = aa.Nonce
	contract.Height = aa.Height
	p.MemStore.Put(contractKey, contract)
	return http.StatusOK, nil
}
