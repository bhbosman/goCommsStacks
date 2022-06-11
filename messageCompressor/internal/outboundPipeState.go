package internal

import (
	"context"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
)

func OutboundPipeState(id string) *common.PipeState {
	return &common.PipeState{
		ID: id,
		Create: func(stackData common.IStackCreateData, ctx context.Context) (interface{}, error) {
			return NewOutboundStackHandler()
		},
		Destroy: func(stackData common.IStackCreateData, pipeData common.IPipeCreateData) error {
			if pd, ok := pipeData.(*OutboundStackHandler); ok {
				return pd.Destroy()
			}
			return goerrors.InvalidType
		},
		Start: func(stackData common.IStackCreateData, pipeData common.IPipeCreateData, ctx context.Context) error {
			if pd, ok := pipeData.(*OutboundStackHandler); ok {
				return pd.Start(ctx)
			}
			return goerrors.InvalidType
		},
		End: func(stackData common.IStackCreateData, pipeData common.IPipeCreateData) error {
			if pd, ok := pipeData.(*OutboundStackHandler); ok {
				return pd.End()
			}
			return goerrors.InvalidType

		},
	}
}
