package internal

import (
	"context"
	"github.com/bhbosman/goCommsStacks/sshConnection/common"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/zap"
	"net"
	"net/url"
	"reflect"
)

func CreateStackState(connectionType model.ConnectionType) *common2.StackState {
	return &common2.StackState{
		Id:          common.StackName,
		HijackStack: true,
		Create: func(
			logger *zap.Logger,
			connectionType model.ConnectionType,
			conn net.Conn,
			url *url.URL,
			ctx context.Context,
			cancelFunc model.ConnectionCancelFunc,
			cfr intf.IConnectionReactorFactoryExtractValues) (common2.IStackCreateData, error) {
			return NewStackData(connectionType, conn, ctx, logger), nil
		},
		Destroy: func(connectionType model.ConnectionType, stackData common2.IStackCreateData) error {
			if closer, ok := stackData.(*StackData); ok {
				return closer.Close()
			}
			return WrongStackDataError(connectionType, stackData)
		},
		Start: func(
			connectionType model.ConnectionType,
			inputStreamForStack common2.IInputStreamForStack,
			stackData common2.IStackCreateData,
			Url *url.URL,
			Ctx context.Context,
			CancelCtxFunc context.CancelFunc,
			ConnectionCancelFunc model.ConnectionCancelFunc,
			ConnectionReactorFactory intf.IConnectionReactorFactoryExtractValues) (common2.IInputStreamForStack, error) {
			if stackDataInstance, ok := stackData.(*StackData); ok {
				return stackDataInstance.Start(Ctx)
			}
			return nil, WrongStackDataError(connectionType, stackData)
		},
		Stop: func(stackData interface{}, endParams common2.StackEndStateParams) error {
			if stop, ok := stackData.(*StackData); ok {
				return stop.Stop()
			}
			return WrongStackDataError(connectionType, stackData)
		},
	}
}

func WrongStackDataError(connectionType model.ConnectionType, stackData interface{}) error {
	return common2.NewWrongStackDataType(
		common.StackName,
		connectionType,
		reflect.TypeOf((*StackData)(nil)),
		reflect.TypeOf(stackData))
}
