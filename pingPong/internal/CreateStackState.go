package internal

import (
	"context"
	"github.com/bhbosman/goCommsStacks/pingPong/common"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goerrors"
	"go.uber.org/zap"
	"net"
	"net/url"
)

func CreateStackState() *common2.StackState {
	return &common2.StackState{
		Id:          common.StackName,
		HijackStack: false,
		Create: func(
			logger *zap.Logger,
			connectionType model.ConnectionType,
			conn net.Conn,
			url *url.URL,
			ctx context.Context,
			cancelFunc model.ConnectionCancelFunc,
			cfr intf.IConnectionReactorFactoryExtractValues) (common2.IStackCreateData, error) {
			return NewStackData(ctx), nil
		},
		Destroy: func(connectionType model.ConnectionType, stackData common2.IStackCreateData) error {
			if sd, ok := stackData.(*StackData); ok {
				return sd.Destroy()
			}
			return goerrors.InvalidType
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
				return conn, sd.Start(Ctx)
			}
			return nil, goerrors.InvalidType
		},
		Stop: func(stackData interface{}, endParams common2.StackEndStateParams) error {
			if sd, ok := stackData.(*StackData); ok {
				return sd.Stop()
			}
			return goerrors.InvalidType
		},
	}
}
