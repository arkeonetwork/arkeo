package cosmos

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32"
	se "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/hashicorp/go-multierror"
)

const (
	DefaultCoinDecimals = 8

	EnvSignerName     = "SIGNER_NAME"
	EnvSignerPassword = "SIGNER_PASSWD"
	EnvChainHome      = "CHAIN_HOME_FOLDER"
)

var (
	KeyringServiceName           = sdk.KeyringServiceName
	NewRoute                     = sdk.NewRoute
	NewKVStoreKeys               = sdk.NewKVStoreKeys
	NewUint                      = sdk.NewUint
	ParseUint                    = sdkmath.ParseUint
	NewInt                       = sdk.NewInt
	NewIntFromString             = sdk.NewIntFromString
	NewDec                       = sdk.NewDec
	ZeroInt                      = sdk.ZeroInt
	ZeroUint                     = sdkmath.ZeroUint
	ZeroDec                      = sdk.ZeroDec
	OneUint                      = sdkmath.OneUint
	NewCoin                      = sdk.NewCoin
	NewCoins                     = sdk.NewCoins
	ParseCoins                   = sdk.ParseCoinsNormalized
	NewDecWithPrec               = sdk.NewDecWithPrec
	NewDecFromBigInt             = sdk.NewDecFromBigInt
	NewIntFromBigInt             = sdk.NewIntFromBigInt
	NewUintFromBigInt            = sdkmath.NewUintFromBigInt
	AccAddressFromBech32         = sdk.AccAddressFromBech32
	VerifyAddressFormat          = sdk.VerifyAddressFormat
	GetFromBech32                = sdk.GetFromBech32
	NewAttribute                 = sdk.NewAttribute
	NewDecFromStr                = sdk.NewDecFromStr
	GetConfig                    = sdk.GetConfig
	NewEvent                     = sdk.NewEvent
	RegisterCodec                = sdk.RegisterLegacyAminoCodec
	NewEventManager              = sdk.NewEventManager
	EventTypeMessage             = sdk.EventTypeMessage
	AttributeKeyModule           = sdk.AttributeKeyModule
	KVStorePrefixIterator        = sdk.KVStorePrefixIterator
	KVStoreReversePrefixIterator = sdk.KVStoreReversePrefixIterator
	NewKVStoreKey                = sdk.NewKVStoreKey
	NewTransientStoreKey         = sdk.NewTransientStoreKey
	NewContext                   = sdk.NewContext

	// nolint
	GetPubKeyFromBech32 = legacybech32.UnmarshalPubKey
	// nolint
	Bech32ifyPubKey         = legacybech32.MarshalPubKey
	Bech32PubKeyTypeConsPub = legacybech32.ConsPK
	Bech32PubKeyTypeAccPub  = legacybech32.AccPK
	Wrapf                   = se.Wrapf
	MustSortJSON            = sdk.MustSortJSON
	CodeUnauthorized        = uint32(4)
	CodeInsufficientFunds   = uint32(5)
)

type (
	Context    = sdk.Context
	Route      = sdk.Route
	Uint       = sdk.Uint
	Int        = sdk.Int
	Coin       = sdk.Coin
	Coins      = sdk.Coins
	AccAddress = sdk.AccAddress
	Attribute  = sdk.Attribute
	Result     = sdk.Result
	Event      = sdk.Event
	Events     = sdk.Events
	Dec        = sdk.Dec
	Msg        = sdk.Msg
	Iterator   = sdk.Iterator
	Handler    = sdk.Handler
	Querier    = sdk.Querier
	TxResponse = sdk.TxResponse
	Account    = authtypes.AccountI
)

var _ sdk.Address = AccAddress{}

func ErrUnknownRequest(msg string) error {
	return se.Wrap(se.ErrUnknownRequest, msg)
}

func ErrInvalidAddress(addr string) error {
	return se.Wrap(se.ErrInvalidAddress, addr)
}

func ErrInvalidCoins(msg string) error {
	return se.Wrap(se.ErrInvalidCoins, msg)
}

func ErrUnauthorized(msg string) error {
	return se.Wrap(se.ErrUnauthorized, msg)
}

func ErrInsufficientCoins(err error, msg string) error {
	return se.Wrap(multierror.Append(se.ErrInsufficientFunds, err), msg)
}
