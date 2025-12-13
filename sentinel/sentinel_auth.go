package sentinel

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"golang.org/x/crypto/sha3"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	QueryArkAuth  = "arkauth"
	QueryContract = "arkcontract"
	ServiceHeader = "arkservice"
)

// Create a map to hold the rate limiters for each visitor and a mutex.
var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

type ContractAuth struct {
	ContractId uint64
	Timestamp  int64
	Signature  []byte
	ChainId    string
}

type ArkAuth struct {
	ContractId uint64
	Spender    common.PubKey
	Nonce      int64
	Signature  []byte
	ChainId    string
}

// String implement fmt.Stringer
func (aa ArkAuth) String() string {
	return GenerateArkAuthString(aa.ContractId, aa.Nonce, aa.Signature, aa.ChainId)
}

func GenerateArkAuthString(contractId uint64, nonce int64, signature []byte, chainId string) string {
	return fmt.Sprintf("%s:%s", GenerateMessageToSign(contractId, nonce, chainId), hex.EncodeToString(signature))
}

func GenerateMessageToSign(contractId uint64, nonce int64, chainId string) string {
	return fmt.Sprintf("%d:%d:", contractId, nonce)
}

func parseContractAuth(raw string) (ContractAuth, error) {
	var auth ContractAuth
	var err error

	parts := strings.SplitN(raw, ":", 3)

	if len(parts) > 0 {
		auth.ContractId, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return auth, err
		}
	}

	if len(parts) > 1 {
		auth.Timestamp, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return auth, err
		}
	}

	if len(parts) > 2 {
		auth.Signature, err = hex.DecodeString(parts[2])
		if err != nil {
			return auth, err
		}
	}
	return auth, nil
}

func parseArkAuth(raw string, configChainId string) (ArkAuth, error) {

	var aa ArkAuth
	var err error

	aa.ChainId = configChainId

	parts := strings.SplitN(raw, ":", 4)

	if len(parts) == 1 {
		// Only contractId provided (for open contracts)

		aa.ContractId, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return aa, err
		}

	} else if len(parts) == 3 {
		// Format: contractId:nonce:signature

		aa.ContractId, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return aa, err
		}

		aa.Nonce, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return aa, err
		}

		aa.Signature, err = hex.DecodeString(parts[2])
		if err != nil {
			return aa, err
		}

	} else if len(parts) == 4 {
		// Format: contractId:pubKey:nonce:signature

		aa.ContractId, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return aa, err
		}

		pubKey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, parts[1])
		if err != nil {
			return aa, err
		}
		aa.Spender, err = common.NewPubKeyFromCrypto(pubKey)
		if err != nil {
			return aa, err
		}

		aa.Nonce, err = strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return aa, err
		}

		aa.Signature, err = hex.DecodeString(parts[3])
		if err != nil {
			return aa, err
		}

	} else {
		return aa, fmt.Errorf("invalid arkauth format")
	}

	return aa, nil

}

func (aa ArkAuth) Validate(provider common.PubKey) error {
	creator, err := provider.GetMyAddress()
	if err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}
	msg := types.NewMsgClaimContractIncome(creator, aa.ContractId, aa.Nonce, aa.Signature)
	err = msg.ValidateBasic()
	return err
}

func (auth ContractAuth) Validate(lastTimestamp int64, client common.PubKey) error {
	if auth.ContractId == 0 {
		return fmt.Errorf("contract id cannot be zero")
	}
	if auth.Timestamp <= lastTimestamp {
		return fmt.Errorf("timestamp must be larger than %d", lastTimestamp)
	}

	pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, client.String())
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%d:%d:%s", auth.ContractId, auth.Timestamp, auth.ChainId)

	// --- DEBUG PRINTS ---
	fmt.Printf("DEBUG: ContractId: %d\n", auth.ContractId)
	fmt.Printf("DEBUG: Timestamp:  %d\n", auth.Timestamp)
	fmt.Printf("DEBUG: ChainId:    '%s'\n", auth.ChainId)
	fmt.Printf("DEBUG: Message:    '%s'\n", msg)
	fmt.Printf("DEBUG: Client PubKey (bech32): %s\n", client.String())
	fmt.Printf("DEBUG: Client PubKey (hex):    %x\n", pk.Bytes())
	fmt.Printf("DEBUG: Signature (hex):        %x\n", auth.Signature)
	// --- END DEBUG PRINTS ---

	if !pk.VerifySignature([]byte(msg), auth.Signature) {
		fmt.Println("DEBUG: Signature verification FAILED")
		return fmt.Errorf("invalid signature")
	}

	fmt.Println("DEBUG: Signature verification PASSED")
	return nil
}

func (auth ContractAuth) String() string {
	sig := hex.EncodeToString(auth.Signature)
	return fmt.Sprintf("Contract Id: %d, Timestamp: %d, Signature: %s", auth.ContractId, auth.Timestamp, sig)
}

func (p Proxy) fetchArkAuth(r *http.Request) (aa ArkAuth, err error) {
	rawHeader := r.Header.Get(QueryArkAuth)
	if len(rawHeader) > 0 {
		aa, err = parseArkAuth(rawHeader, p.Config.SourceChain)
		if err != nil {
			return aa, err
		}
		return aa, nil
	}
	args := r.URL.Query()
	raw, aaOK := args[QueryArkAuth]
	if aaOK {
		aa, err = parseArkAuth(raw[0], p.Config.SourceChain)
		if err != nil {
			return aa, err
		}
	}
	return aa, nil
}

func (p Proxy) auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusOK)
			return
		}

		aa, err := p.fetchArkAuth(r)
		remoteAddr := p.getRemoteAddr(r)
		if err != nil {
			// Attempt to fetch contract id from query
			args := r.URL.Query()
			contractIdStr := args.Get("contract_id")
			// Try "arkauth" param (could just be a contractId)
			if contractIdStr == "" {
				arkauthParam := args.Get(QueryArkAuth)
				parts := strings.SplitN(arkauthParam, ":", 2)
				if len(parts) > 0 {
					contractIdStr = parts[0]
				}
			}
			contractId, _ := strconv.ParseUint(contractIdStr, 10, 64)
			contract, cErr := p.MemStore.Get(strconv.FormatUint(contractId, 10))
			whitelisted := false
			if cErr == nil && !contract.Client.IsEmpty() {
				conf, _ := p.ContractConfigStore.Get(contract.Id)
				addr := remoteAddr
				if strings.Contains(addr, ":") {
					addr, _, _ = strings.Cut(addr, ":")
				}
				for _, ip := range conf.WhitelistIPAddresses {
					if strings.EqualFold(addr, ip) {
						whitelisted = true
						break
					}
				}
			}
			if cErr == nil && contract.IsOpenAuthorization() && whitelisted {
				w.Header().Set("tier", "paid")
				next.ServeHTTP(w, r)
				return
			}
			// Otherwise, as before:
			p.logger.Error("failed to parse ark auth", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p.logger.Info("DEBUG: ArkAuth parsed",
			"raw", r.URL.Query().Get(QueryArkAuth),
			"arkAuth", aa.String(),
		)

		var contract types.Contract
		if aa.ContractId > 0 {
			contract, err = p.MemStore.Get(strconv.FormatUint(aa.ContractId, 10))
			if err != nil {
				p.logger.Error("failed to fetch contract", "error", err)
			}

			// Do not serve expired subscription contracts
			if contract.Id != 0 && contract.Type == types.ContractType_SUBSCRIPTION && contract.IsExpired(p.MemStore.GetHeight()) {
				http.Error(w, "subscription contract expired", http.StatusPaymentRequired)
				return
			}
		}

		// Always require a non-empty signature for PAY_AS_YOU_GO contracts
		if contract.Type == types.ContractType_PAY_AS_YOU_GO && len(aa.Signature) == 0 {
			p.logger.Error("Missing signature for PAY_AS_YOU_GO contract", "contract_id", contract.Id)
			http.Error(w, "Signature required for pay-as-you-go contracts", http.StatusUnauthorized)
			return
		}

		// fetch conf and check whitelist/IP/rate
		whitelisted := false
		if !contract.Client.IsEmpty() {
			conf, err := p.ContractConfigStore.Get(contract.Id)
			if err != nil {
				p.logger.Error("failed to fetch contract", "error", err)
			} else {
				w = p.enableCORS(w, conf.CORs)

				// Check IP whitelist
				addr := remoteAddr
				if strings.Contains(addr, ":") {
					addr, _, _ = strings.Cut(addr, ":")
				}
				for _, ip := range conf.WhitelistIPAddresses {
					if strings.EqualFold(addr, ip) {
						whitelisted = true
						break
					}
				}

				if len(conf.WhitelistIPAddresses) > 0 && !whitelisted {
					p.logger.Info("DEBUG: IP not in contract whitelist, falling through to free tier", "addr", addr)
					// Do not return; fall through to free tier logic
				}

			}
		}

		// treat open contracts + whitelisted IPs as paid
		if err == nil && (aa.Validate(p.Config.ProviderPubKey) == nil || (contract.IsOpenAuthorization() && whitelisted)) {
			w.Header().Set("tier", "paid")

			// Determine service name from header or URL path.
			rawHeaderService := r.Header.Get(ServiceHeader)
			serviceName := rawHeaderService
			if serviceName == "" {
				parts := strings.Split(r.URL.Path, "/")
				if len(parts) > 1 {
					serviceName = parts[1]
				}
			}

			// Log the raw header and derived service name.
			p.logger.Info("DEBUG: service header and path",
				"header_service", rawHeaderService,
				"url_path", r.URL.Path,
				"derived_service_name", serviceName,
			)

			// Try to resolve via dynamic registry; fallback to legacy parse for logging only.
			p.serviceMu.RLock()
			reqServiceID, ok := p.serviceIDs[strings.ToLower(serviceName)]
			p.serviceMu.RUnlock()

			if ok {
				p.logger.Info("DEBUG: service match check",
					"contract_id", contract.Id,
					"contract_service_enum", contract.Service,
					"request_service_id", reqServiceID,
				)

				if int32(reqServiceID) != int32(contract.Service) {
					p.logger.Error("Service match failed",
						"serviceName", serviceName,
						"contract_id", contract.Id,
						"contract_service_enum", contract.Service,
						"request_service_id", reqServiceID,
					)
					http.Error(w, "Service mismatch", http.StatusUnauthorized)
					return
				}
			} else {
				// Legacy logging fallback
				ser, serr := common.NewService(serviceName)
				p.logger.Info("DEBUG: service not in registry; legacy parse",
					"serviceName", serviceName,
					"contract_id", contract.Id,
					"contract_service_enum", contract.Service,
					"parsed_service_enum", ser,
					"new_service_err", serr,
				)
				// allow if registry doesn’t know it (dynamic addition)
			}

			httpCode, tierErr := p.paidTier(aa, remoteAddr)
			if tierErr == nil {
				next.ServeHTTP(w, r)
				return
			}
			p.logger.Error("DEBUG: paidTier failed", "error", tierErr, "http_code", httpCode)
			http.Error(w, tierErr.Error(), httpCode)
			return
		}

		// If the contract is present and type is PAY_AS_YOU_GO, do not fall through to the free tier
		if contract.Id != 0 && contract.Type == types.ContractType_PAY_AS_YOU_GO {
			http.Error(w, "Pay-as-you-go contracts do not fall through to free tier.", http.StatusUnauthorized)
			return
		}

		w.Header().Set("tier", "free")
		httpCode, err := p.freeTier(remoteAddr)
		if err != nil {
			http.Error(w, err.Error(), httpCode)
			return
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
	// Extract IP from "IP:port"
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip, _, _ = strings.Cut(ip, ":")
	}
	return ip
}

func (p Proxy) isRateLimited(contractId uint64, key string, limitTokens int, windowSeconds int) bool {
	mu.Lock()
	defer mu.Unlock()

	key = fmt.Sprintf("%d-%s", contractId, key)
	limiter, exists := visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(float64(limitTokens)/float64(windowSeconds)), limitTokens)
		visitors[key] = limiter
	}

	allowed := limiter.Allow()
	p.logger.Info("DEBUG: Rate limiting check",
		"key", key,
		"limit_per_sec", limiter.Limit(),
		"burst", limiter.Burst(),
		"allowed", allowed,
	)
	if allowed {
		p.logger.Debug("DEBUG: Rate limit result", "status", "allowed", "key", key)
	} else {
		p.logger.Debug("DEBUG: Rate limit result", "status", "rate limited", "key", key)
	}
	return !allowed
}

func (p Proxy) freeTier(remoteAddr string) (int, error) {
	if ok := p.isRateLimited(0, remoteAddr, p.Config.FreeTierRateLimit, 60); ok {
		return http.StatusTooManyRequests, fmt.Errorf("free client is rate limited (%s)", http.StatusText(429))
	}

	return http.StatusOK, nil
}

func (p Proxy) paidTier(aa ArkAuth, remoteAddr string) (code int, err error) {

	// Fetch contract by ID; error if not found or datastore issue.
	key := strconv.FormatUint(aa.ContractId, 10)
	contract, err := p.MemStore.Get(key)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err)
	}

	// Ensure spender (client) is recorded in the claim even when arkauth is 3-part.
	if aa.Spender.IsEmpty() {
		aa.Spender = contract.Client
		p.logger.Debug("paidTier: inferred spender from contract client", "spender", aa.Spender.String())
	}

	// Check if the contract has expired (based on current chain height).
	// If expired, require the client to open a new contract before continuing.
	if contract.IsExpired(p.MemStore.GetHeight()) {
		return http.StatusPaymentRequired, fmt.Errorf("open a contract")
	}

	// check if we've exceeded the total number of pay-as-you-go queries
	if contract.IsPayAsYouGo() {
		if contract.Deposit.IsNil() || contract.Deposit.LT(cosmos.NewInt(aa.Nonce*contract.Rate.Amount.Int64())) {
			return http.StatusPaymentRequired, fmt.Errorf("contract spent")
		}
	}

	// Enforce per-contract paid tier rate limiting.
	if ok := p.isRateLimited(contract.Id, remoteAddr, int(contract.QueriesPerMinute), 60); ok {
		return http.StatusTooManyRequests, fmt.Errorf("paid client is rate limited (%s)", http.StatusText(429))
	}

	// For open authorization (subscription) contracts, skip PAYG nonce/signature tracking.
	// Open contracts do not require per-request client signatures or nonce/accounting.
	if contract.IsOpenAuthorization() {
		p.logger.Debug("paidTier: open authorization contract; skipping claim enqueue",
			"contract_id", contract.Id,
			"nonce", aa.Nonce,
			"service", contract.Service.String(),
			"spender", aa.Spender.String(),
		)
		return http.StatusOK, nil
	}

	// Optional self-verify so only claimable entries are stored.
	// Preferred: chain-style SHA-256("<cid>:<nonce>:")
	// Compat: raw preimage, Keccak(preimage), and EIP-191 personal_sign over preimage.
	{
		pre := fmt.Sprintf("%d:%d:", aa.ContractId, aa.Nonce)
		digest := sha256.Sum256([]byte(pre))

		pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, aa.Spender.String())
		if err != nil {
			return http.StatusUnauthorized, fmt.Errorf("invalid client pubkey: %w", err)
		}

		// 1) chain preferred: SHA-256(preimage)
		ok := pk.VerifySignature(digest[:], aa.Signature)

		// 2) compat: raw preimage
		if !ok {
			ok = pk.VerifySignature([]byte(pre), aa.Signature)
		}

		// 3) compat: keccak256(preimage)
		if !ok {
			k := sha3.NewLegacyKeccak256()
			k.Write([]byte(pre))
			ok = pk.VerifySignature(k.Sum(nil), aa.Signature)
		}

		// 4) compat: EIP-191 personal_sign (Ethereum prefix)
		if !ok {
			prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(pre))
			k := sha3.NewLegacyKeccak256()
			k.Write([]byte(prefix))
			k.Write([]byte(pre))
			ok = pk.VerifySignature(k.Sum(nil), aa.Signature)
		}

		if !ok {
			return http.StatusUnauthorized, fmt.Errorf("invalid signature for client")
		}
	}

	// Create or update the claim for this contract request:
	// - Build a new claim using the provided contract ID, spender, nonce, and signature.
	// - If a claim already exists for this contract, fetch it to check for replay or out-of-order requests.
	// - If the incoming nonce is not strictly greater than the stored nonce, reject the request (prevent replay or duplicate).
	sig := hex.EncodeToString(aa.Signature)
	claim := NewClaim(aa.ContractId, aa.Spender, aa.Nonce, sig)
	if p.ClaimStore.Has(key) {
		var err error
		claim, err = p.ClaimStore.Get(key)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err)
		}
		if claim.Nonce >= aa.Nonce {
			return http.StatusBadRequest, fmt.Errorf("bad nonce (%d/%d)", aa.Nonce, claim.Nonce)
		}
	}

	// Update claim and contract state for Pay-As-You-Go contracts.
	//
	// - Stores the new nonce and signature in the claim store (for settlement and anti-replay).
	// - Marks the claim as unclaimed (pending batch settlement).
	// - Updates the contract's nonce in memory to track usage progression.
	// - Calculates and logs usage stats (used/remaining deposit, per-query cost).
	claim.Provider = p.Config.ProviderPubKey
	claim.Spender = aa.Spender
	claim.Nonce = aa.Nonce
	claim.Signature = sig
	claim.Claimed = false
	if err := p.ClaimStore.Set(claim); err != nil {
		p.logger.Error("paidTier: failed to persist claim",
			"contract_id", claim.ContractId,
			"nonce", claim.Nonce,
			"spender", claim.Spender.String(),
			"error", err,
		)
		return http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err)
	}
	p.logger.Info("paidTier: claim stored",
		"contract_id", claim.ContractId,
		"nonce", claim.Nonce,
		"spender", claim.Spender.String(),
		"service", contract.Service.String(),
	)
	contract.Nonce = aa.Nonce
	p.MemStore.Put(contract)

	used := contract.Nonce * contract.Rate.Amount.Int64()
	remaining := contract.Deposit.Int64() - used

	p.logger.Debug("Contract Usage: ",
		"contract_id", contract.Id,
		"nonce", contract.Nonce,
		"deposit", contract.Deposit.Int64(),
		"used", used,
		"remaining", remaining,
		"cost_per_query", contract.Rate.Amount.Int64(),
		"denom", contract.Rate.Denom,
	)

	return http.StatusOK, nil
}

func (p Proxy) enableCORS(w http.ResponseWriter, cors CORs) http.ResponseWriter {
	if len(cors.AllowOrigins) > 0 {
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(cors.AllowOrigins, ", "))
	}
	if len(cors.AllowMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(cors.AllowMethods, ", "))
	}
	if len(cors.AllowHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.AllowHeaders, ", "))
	}
	return w
}
