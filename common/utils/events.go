package utils

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func ParseTypedEvent(result tmCoreTypes.ResultEvent, eventType string) (proto.Message, error) {
	var (
		msg         proto.Message
		eventDataTx tmtypes.EventDataTx
		ok          bool
	)
	if eventDataTx, ok = result.Data.(tmtypes.EventDataTx); !ok {
		return msg, fmt.Errorf("failed cast %T to EventDataTx", result.Data)
	}

	for _, evt := range eventDataTx.TxResult.Result.Events {
		if evt.Type == eventType {
			return sdk.ParseTypedEvent(evt)
		}
	}

	return msg, fmt.Errorf("event %s not found", eventType)
}
