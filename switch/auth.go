package switchd

import (
	"fmt"
	"log"
	"mercury/common"
	"mercury/common/cosmos"
	"mercury/switch/conf"
	"mercury/x/mercury/types"
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
)

// Create a map to hold the rate limiters for each visitor and a mutex.
var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func auth(config conf.Configuration, mem *MemStore, claimStore *ClaimStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		args := r.URL.Query()

		// check if request has auth args
		height, heightOK := args[QueryHeight]
		nonce, nonceOK := args[QueryNonce]
		spender, spenderOK := args[QuerySpender]
		sig, sigOK := args[QuerySig]

		if heightOK && nonceOK && spenderOK && sigOK {
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

			if ok := paidTier(heightInt, nonceInt, spender[0], sig[0], config, mem, claimStore, w, r); !ok {
				return
			}
		} else {
			if ok := freeTier(config, mem, w, r); !ok {
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func freeTier(config conf.Configuration, mem *MemStore, w http.ResponseWriter, r *http.Request) bool {
	// Get the IP address for the current user.
	key, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}

	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(config.FreeTierRateLimitDuration), config.FreeTierRateLimit)
		visitors[key] = limiter
	}

	if !limiter.Allow() {
		log.Printf("Entity %s has been rate limited :(", key)
		http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
		return false
	}

	return true
}

func paidTier(height, nonce int64, spender, sig string, config conf.Configuration, mem *MemStore, claimStore *ClaimStore, w http.ResponseWriter, r *http.Request) bool {
	log.Println("Checking Authorization...")

	creator, err := config.ProviderPubKey.GetMyAddress()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}

	spenderPK, err := common.NewPubKey(spender)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Invalid spender pubkey", http.StatusBadRequest)
		return false
	}

	// TODO: fetch chain from URL
	chain := "btc-mainnet"

	msg := types.NewMsgClaimContractIncome(creator.String(), config.ProviderPubKey, chain, spenderPK, nonce, height, sig)
	if err := msg.ValidateBasic(); err != nil {
		log.Print(err.Error())
		http.Error(w, "Invalid spender pubkey", http.StatusBadRequest)
		return false
	}

	contractKey := mem.Key(config.ProviderPubKey.String(), chain, spender)
	contract, err := mem.Get(contractKey)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}

	if contract.IsClose(mem.GetHeight()) {
		http.Error(w, "open a contract", http.StatusPaymentRequired)
		return false
	}

	key := fmt.Sprintf("%s-%s", chain, spender)
	chainInt, err := common.NewChain("btc-mainnet")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}

	claim := NewClaim(config.ProviderPubKey, chainInt, spenderPK, nonce, height, sig)
	if claimStore.Has(key) {
		var err error
		claim, err = claimStore.Get(key)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return false
		}
		if claim.Height == height && contract.Height == height && claim.Nonce >= nonce {
			log.Print(err.Error())
			http.Error(w, fmt.Sprintf("bad nonce (%d/%d)", nonce, claim.Nonce), http.StatusBadRequest)
			return false
		}
	}

	// check if we've exceed the total number of pay-as-you-go queries
	if contract.Type == types.ContractType_PayAsYouGo {
		if contract.Deposit.LT(cosmos.NewInt(nonce * contract.Rate)) {
			http.Error(w, "open a contract", http.StatusPaymentRequired)
			return false
		}
	}

	mu.Lock()
	defer mu.Unlock()

	limitTokens := config.SubTierRateLimit
	limitDuration := config.SubTierRateLimitDuration
	if contract.Type == types.ContractType_PayAsYouGo {
		limitTokens = config.AsGoTierRateLimit
		limitDuration = config.AsGoTierRateLimitDuration
	}

	limiter, exists := visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(limitDuration), limitTokens)
		visitors[key] = limiter
	}

	if !limiter.Allow() {
		log.Printf("Entity %s has been rate limited :(", key)
		http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
		return false
	}

	claim.Nonce = nonce
	claim.Height = height
	claim.Signature = sig
	claim.Claimed = false
	if err := claimStore.Set(claim); err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}
	contract.Nonce = nonce
	contract.Height = height
	mem.Put(contractKey, contract)

	return true
}
