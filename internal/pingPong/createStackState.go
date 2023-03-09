package pingPong

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
) common2.IStackState {
	return common2.NewStackState(
		goCommsDefinitions.PingPongStackName,
		false,
		func() (common2.IStackCreateData, error) {
			return NewStackData(
				ctx,
				goFunctionCounter,
			), nil
		},
		func(connectionType model.ConnectionType, stackData common2.IStackCreateData) error {
			if sd, ok := stackData.(*Data); ok {
				return sd.Destroy()
			}
			return goerrors.InvalidType
		},
		func(
			conn common2.IInputStreamForStack,
			stackData common2.IStackCreateData,
			ToReactorFunc rxgo.NextFunc,
		) (common2.IInputStreamForStack, error) {
			if sd, ok := stackData.(*Data); ok {
				return conn, sd.Start(ctx)
			}
			return nil, goerrors.InvalidType
		},
		func(stackData interface{}) error {
			if sd, ok := stackData.(*Data); ok {
				return sd.Stop()
			}
			return goerrors.InvalidType
		},
	)
}
