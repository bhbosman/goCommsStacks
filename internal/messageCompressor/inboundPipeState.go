package messageCompressor

import (
	"context"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
)

func InboundPipeState(id string) *common.PipeState {
	return &common.PipeState{
		ID: id,
		OnCreate: func(stackData common.IStackCreateData, ctx context.Context) (interface{}, error) {
			return NewInboundStackHandler()
		},
		OnDestroy: func(stackData common.IStackCreateData, pipeData common.IPipeCreateData) error {
			if pd, ok := pipeData.(*inboundStackHandler); ok {
				return pd.Destroy()
			}
			return goerrors.InvalidType
		},
		OnStart: func(stackData common.IStackCreateData, pipeData common.IPipeCreateData, ctx context.Context) error {
			if pd, ok := pipeData.(*inboundStackHandler); ok {
				return pd.Start(ctx)
			}
			return goerrors.InvalidType
		},
		OnEnd: func(stackData common.IStackCreateData, pipeData common.IPipeCreateData) error {
			if pd, ok := pipeData.(*inboundStackHandler); ok {
				return pd.End()
			}
			return goerrors.InvalidType
		},
	}
}
