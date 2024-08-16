package cosmos

import (
	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32" //nolint:staticcheck
	se "github.com/cosmos/cosmos-sdk/types/errors"
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
	NewRoute                     = baseapp.NewMsgServiceRouter
	NewKVStoreKeys               = storetypes.NewKVStoreKeys
	NewUint                      = sdkmath.NewUint
	ParseUint                    = sdkmath.ParseUint
	NewInt                       = sdkmath.NewInt
	NewIntFromString             = sdkmath.NewIntFromString
	NewDec                       = sdkmath.LegacyNewDec
	ZeroInt                      = sdkmath.ZeroInt
	ZeroUint                     = sdkmath.ZeroUint
	ZeroDec                      = sdkmath.LegacyZeroDec
	OneUint                      = sdkmath.OneUint
	NewInt64Coin                 = sdk.NewInt64Coin
	NewCoin                      = sdk.NewCoin
	NewCoins                     = sdk.NewCoins
	ParseCoin                    = sdk.ParseCoinNormalized
	ParseCoins                   = sdk.ParseCoinsNormalized
	NewDecWithPrec               = sdkmath.LegacyNewDecWithPrec
	NewDecFromBigInt             = sdkmath.LegacyNewDecFromBigInt
	NewIntFromBigInt             = sdkmath.NewIntFromBigInt
	NewUintFromBigInt            = sdkmath.NewUintFromBigInt
	ValAddressFromBech32         = sdk.ValAddressFromBech32
	AccAddressFromBech32         = sdk.AccAddressFromBech32
	VerifyAddressFormat          = sdk.VerifyAddressFormat
	GetFromBech32                = sdk.GetFromBech32
	NewAttribute                 = sdk.NewAttribute
	NewDecFromStr                = sdkmath.LegacyNewDecFromStr
	GetConfig                    = sdk.GetConfig
	NewEvent                     = sdk.NewEvent
	RegisterCodec                = sdk.RegisterLegacyAminoCodec
	NewEventManager              = sdk.NewEventManager
	EventTypeMessage             = sdk.EventTypeMessage
	AttributeKeyModule           = sdk.AttributeKeyModule
	KVStorePrefixIterator        = storetypes.KVStorePrefixIterator
	KVStoreReversePrefixIterator = storetypes.KVStoreReversePrefixIterator
	NewKVStoreKey                = storetypes.NewKVStoreKey
	NewTransientStoreKey         = storetypes.NewTransientStoreKey
	NewContext                   = sdk.NewContext

	// nolint
	GetPubKeyFromBech32 = legacybech32.UnmarshalPubKey
	// nolint
	Bech32ifyPubKey         = legacybech32.MarshalPubKey
	Bech32PubKeyTypeConsPub = legacybech32.ConsPK
	Bech32PubKeyTypeAccPub  = legacybech32.AccPK
	Bech32PubkeyTypeValPK   = legacybech32.ValPK
	Wrapf                   = errors.Wrapf
	MustSortJSON            = sdk.MustSortJSON
	CodeUnauthorized        = uint32(4)
	CodeInsufficientFunds   = uint32(5)
)

type (
	Context = sdk.Context
	// Route      = baseapp.MessageRouter
	Uint       = sdkmath.Uint
	Int        = sdkmath.Int
	Coin       = sdk.Coin
	Coins      = sdk.Coins
	AccAddress = sdk.AccAddress
	ValAddress = sdk.ValAddress
	Attribute  = sdk.Attribute
	Result     = sdk.Result
	Event      = sdk.Event
	Events     = sdk.Events
	Dec        = sdkmath.LegacyDec
	Msg        = sdk.Msg
	Iterator   = storetypes.Iterator
	// Handler    = baseapp.MsgServiceHandler
	// Querier    = func(ctx Context, path []string, req abci.RequestQuery) ([]byte, error)
	TxResponse = sdk.TxResponse
	Account    = sdk.AccountI
)

var _ sdk.Address = AccAddress{}

func ErrUnknownRequest(msg string) error {
	return errors.Wrap(se.ErrUnknownRequest, msg)
}

func ErrInvalidAddress(addr string) error {
	return errors.Wrap(se.ErrInvalidAddress, addr)
}

func ErrInvalidCoins(msg string) error {
	return errors.Wrap(se.ErrInvalidCoins, msg)
}

func ErrUnauthorized(msg string) error {
	return errors.Wrap(se.ErrUnauthorized, msg)
}

func ErrInsufficientCoins(err error, msg string) error {
	return errors.Wrap(multierror.Append(se.ErrInsufficientFunds, err), msg)
}
