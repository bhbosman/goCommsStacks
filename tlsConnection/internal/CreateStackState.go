package internal

import (
	"context"
	"github.com/bhbosman/goCommsStacks/tlsConnection/common"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/zap"
	"net"
	"net/url"
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
			return NewStackData(connectionType, conn, ctx, logger)
		},
		Destroy: func(connectionType model.ConnectionType, stackData common2.IStackCreateData) error {
			if sd, ok := stackData.(*StackData); ok {
				return sd.Close()
			}
			return WrongStackDataError(connectionType, stackData)
		},
		Start: func(
			connectionType model.ConnectionType,
			conn common2.IInputStreamForStack,
			stackData common2.IStackCreateData,
			Url *url.URL,
			Ctx context.Context,
			CancelCtxFunc context.CancelFunc,
			ConnectionCancelFunc model.ConnectionCancelFunc,
			ConnectionReactorFactory intf.IConnectionReactorFactoryExtractValues) (common2.IInputStreamForStack, error) {
			if sd, ok := stackData.(*StackData); ok {
				return sd.Start(connectionType, ConnectionCancelFunc, Ctx, CancelCtxFunc)
			}
			return nil, WrongStackDataError(connectionType, stackData)
		},
		Stop: func(stackData interface{}, endParams common2.StackEndStateParams) error {
			if sd, ok := stackData.(*StackData); ok {
				return sd.Close()
			}
			return WrongStackDataError(connectionType, stackData)
		},
	}
}
