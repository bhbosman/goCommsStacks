package internal

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
)

func CreateStackState(
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
) *common2.StackState {
	return &common2.StackState{
		Id:          goCommsDefinitions.PingPongStackName,
		HijackStack: false,
		Create: func() (common2.IStackCreateData, error) {
			return NewStackData(
				ctx,
				goFunctionCounter,
			), nil
		},
		Destroy: func(connectionType model.ConnectionType, stackData common2.IStackCreateData) error {
			if sd, ok := stackData.(*data); ok {
				return sd.Destroy()
			}
			return goerrors.InvalidType
		},
		Start: func(
			conn common2.IInputStreamForStack,
			stackData common2.IStackCreateData,
			ToReactorFunc rxgo.NextFunc,
		) (common2.IInputStreamForStack, error) {
			if sd, ok := stackData.(*data); ok {
				return conn, sd.Start(ctx)
			}
			return nil, goerrors.InvalidType
		},
		Stop: func(stackData interface{}, endParams common2.StackEndStateParams) error {
			if sd, ok := stackData.(*data); ok {
				return sd.Stop()
			}
			return goerrors.InvalidType
		},
	}
}
