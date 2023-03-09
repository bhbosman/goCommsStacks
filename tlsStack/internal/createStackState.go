package internal

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"net"
)

func CreateStackState(
	logger *zap.Logger,
	connectionType model.ConnectionType,
	conn net.Conn,
	ctx context.Context,
	ctxCancelFunc context.CancelFunc,
	cancelFunc model.ConnectionCancelFunc,
	goFunctionCounter GoFunctionCounter.IService,
) common.IStackState {
	return common.NewStackState(
		goCommsDefinitions.TlsStackName,
		true,
		func() (common.IStackCreateData, error) {
			return NewStackData(
				connectionType,
				conn,
				ctx,
				ctxCancelFunc,
				cancelFunc,
				logger,
				goFunctionCounter,
			)
		},
		func(connectionType model.ConnectionType, stackData common.IStackCreateData) error {
			if sd, ok := stackData.(*Data); ok {
				return sd.Close()
			}
			return WrongStackDataError(connectionType, stackData)
		},
		func(
			conn common.IInputStreamForStack,
			stackData common.IStackCreateData,
			ToReactorFunc rxgo.NextFunc,
		) (common.IInputStreamForStack, error) {
			if sd, ok := stackData.(*Data); ok {
				return sd.Start()
			}
			return nil, WrongStackDataError(connectionType, stackData)
		},
		func(stackData interface{}) error {
			if sd, ok := stackData.(*Data); ok {
				return sd.Close()
			}
			return WrongStackDataError(connectionType, stackData)
		},
	)
}
