package tlsStack

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"net"
)

func createStackState(
	logger *zap.Logger,
	connectionType model.ConnectionType,
	conn net.Conn,
	ctx context.Context,
	ctxCancelFunc context.CancelFunc,
	cancelFunc model.ConnectionCancelFunc,
	goFunctionCounter GoFunctionCounter.IService,
) *common2.StackState {
	return &common2.StackState{
		Id:          goCommsDefinitions.TlsStackName,
		HijackStack: true,
		Create: func() (common2.IStackCreateData, error) {
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
		Destroy: func(connectionType model.ConnectionType, stackData common2.IStackCreateData) error {
			if sd, ok := stackData.(*data); ok {
				return sd.Close()
			}
			return WrongStackDataError(connectionType, stackData)
		},
		Start: func(
			conn common2.IInputStreamForStack,
			stackData common2.IStackCreateData,
			ToReactorFunc rxgo.NextFunc,
		) (common2.IInputStreamForStack, error) {
			if sd, ok := stackData.(*data); ok {
				return sd.Start()
			}
			return nil, WrongStackDataError(connectionType, stackData)
		},
		Stop: func(stackData interface{}, endParams common2.StackEndStateParams) error {
			if sd, ok := stackData.(*data); ok {
				return sd.Close()
			}
			return WrongStackDataError(connectionType, stackData)
		},
	}
}
