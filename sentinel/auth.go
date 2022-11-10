package sentinel

import (
	"arkeo/common"
	"arkeo/common/cosmos"
	"arkeo/x/arkeo/types"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"golang.org/x/time/rate"
)

const (
	QueryHeight  = "__cheight"
	QueryNonce   = "__nonce"
	QuerySpender = "__spender"
	QuerySig     = "__signature"
	QueryChain   = "__chain"
)

// Create a map to hold the rate limiters for each visitor and a mutex.
var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func (p Proxy) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		args := r.URL.Query()

		// check if request has auth args
		height, heightOK := args[QueryHeight]
		nonce, nonceOK := args[QueryNonce]
		spender, spenderOK := args[QuerySpender]
		chain, chainOK := args[QueryChain]
		sig, sigOK := args[QuerySig]

		if heightOK && nonceOK && spenderOK && chainOK && sigOK {
			nonceInt, err := strconv.ParseInt(nonce[0], 10, 64)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, "Invalid nonce", http.StatusBadRequest)
				return
			}
			heightInt, err := strconv.ParseInt(height[0], 10, 64)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, "Invalid nonce", http.StatusBadRequest)
				return
			}

			httpCode, err := p.paidTier(heightInt, nonceInt, chain[0], spender[0], sig[0])
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), httpCode)
				return
			}
		} else {
			httpCode, err := p.freeTier(r.RemoteAddr)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), httpCode)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (p Proxy) freeTier(remoteAddr string) (int, error) {
	// Get the IP address for the current user.
	key, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Internal Server Error: %s", err)
	}

	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(p.Config.FreeTierRateLimitDuration), p.Config.FreeTierRateLimit)
		visitors[key] = limiter
	}

	if !limiter.Allow() {
		return http.StatusTooManyRequests, fmt.Errorf(http.StatusText(429))
	}

	return http.StatusOK, nil
}

func (p Proxy) paidTier(height, nonce int64, chain, spender, signature string) (int, error) {
	creator, err := p.Config.ProviderPubKey.GetMyAddress()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Internal Server Error: %s", err)
	}

	spenderPK, err := common.NewPubKey(spender)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("Invalid spender pubkey")
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("unable to decode signature")
	}

	msg := types.NewMsgClaimContractIncome(creator.String(), p.Config.ProviderPubKey, chain, spenderPK, nonce, height, sig)
	if err := msg.ValidateBasic(); err != nil {
		return http.StatusBadRequest, fmt.Errorf("bad claim: %s", err)
	}

	contractKey := p.MemStore.Key(p.Config.ProviderPubKey.String(), chain, spender)
	contract, err := p.MemStore.Get(contractKey)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Internal Server Error: %s", err)
	}

	if contract.IsClose(p.MemStore.GetHeight()) {
		return http.StatusPaymentRequired, fmt.Errorf("open a contract")
	}

	chainInt, err := common.NewChain(chain)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Internal Server Error: %s", err)
	}

	key := fmt.Sprintf("%d-%s", chainInt, spender)

	claim := NewClaim(p.Config.ProviderPubKey, chainInt, spenderPK, nonce, height, signature)
	if p.ClaimStore.Has(key) {
		var err error
		claim, err = p.ClaimStore.Get(key)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Internal Server Error: %s", err)
		}
		if claim.Height == height && contract.Height == height && claim.Nonce >= nonce {
			return http.StatusBadRequest, fmt.Errorf("bad nonce (%d/%d)", nonce, claim.Nonce)
		}
	}

	// check if we've exceed the total number of pay-as-you-go queries
	if contract.Type == types.ContractType_PayAsYouGo {
		if contract.Deposit.LT(cosmos.NewInt(nonce * contract.Rate)) {
			return http.StatusPaymentRequired, fmt.Errorf("open a contract")
		}
	}

	mu.Lock()
	defer mu.Unlock()

	limitTokens := p.Config.SubTierRateLimit
	limitDuration := p.Config.SubTierRateLimitDuration
	if contract.Type == types.ContractType_PayAsYouGo {
		limitTokens = p.Config.AsGoTierRateLimit
		limitDuration = p.Config.AsGoTierRateLimitDuration
	}

	limiter, exists := visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(limitDuration), limitTokens)
		visitors[key] = limiter
	}

	if !limiter.Allow() {
		return http.StatusTooManyRequests, fmt.Errorf(http.StatusText(429))
	}

	claim.Nonce = nonce
	claim.Height = height
	claim.Signature = signature
	claim.Claimed = false
	if err := p.ClaimStore.Set(claim); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Internal Server Error: %s", err)
	}
	contract.Nonce = nonce
	contract.Height = height
	p.MemStore.Put(contractKey, contract)

	return http.StatusOK, nil
}
