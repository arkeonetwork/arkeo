package utils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func ParseTypedEvent(result tmCoreTypes.ResultEvent, eventType string) (proto.Message, error) {
	var (
		msg proto.Message
	)
	switch v := result.Data.(type) {
	case tmtypes.EventDataTx:
		for _, evt := range v.TxResult.Result.Events {
			if evt.Type == eventType {
				return sdk.ParseTypedEvent(evt)
			}
		}
	case tmtypes.EventDataNewBlock:
		for _, evt := range v.ResultEndBlock.Events {
			if evt.Type == eventType {
				return sdk.ParseTypedEvent(evt)
			}
		}
	}

	return msg, fmt.Errorf("event %s not found", eventType)
}

func MakeResultEvent(sdkEvent sdk.Event, resultTx *tmCoreTypes.ResultTx) tmCoreTypes.ResultEvent {
	evts := make(map[string][]string, len(sdkEvent.Attributes))
	for _, attr := range sdkEvent.Attributes {
		evts[string(attr.Key)] = []string{string(attr.Value)}
	}

	abciEvents := []abciTypes.Event{{
		Type:       sdkEvent.Type,
		Attributes: sdkEvent.Attributes,
	}}
	_ = abciEvents

	query := fmt.Sprintf("tm.event = 'Tx' AND message.action='/%s'", sdkEvent.Type)
	v := tmCoreTypes.ResultEvent{
		Query:  query,
		Events: evts,
	}
	if resultTx != nil {
		v.Data = tmtypes.EventDataTx{
			TxResult: abciTypes.TxResult{
				Height: resultTx.Height,
				Index:  resultTx.Index,
				Tx:     resultTx.Tx,
				Result: resultTx.TxResult,
			},
		}
	} else {
		v.Data = tmtypes.EventDataNewBlock{}
	}
	return v
}
