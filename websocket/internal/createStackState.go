package internal

import (
	"context"
	"github.com/bhbosman/goCommsStacks/websocket/common"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/zap"
	"io"
	"net"
	"net/url"
)

func CreateStackState(connectionType model.ConnectionType, stackCancelFunc model.ConnectionCancelFunc) *common2.StackState {
	return &common2.StackState{
		Id:          common.StackName,
		HijackStack: true,
		Create: func(
			logger *zap.Logger,
			connectionType model.ConnectionType,
			conn net.Conn, Url *url.URL,
			ctx context.Context,
			cancelFunc model.ConnectionCancelFunc,
			cfr intf.IConnectionReactorFactoryExtractValues) (common2.IStackCreateData, error) {
			var err error
			var data *StackData
			data, err = NewStackData(
				logger,
				connectionType,
				conn,
				stackCancelFunc,
				ctx)
			if err != nil {
				return nil, err
			}
			return data, err
		},
		Destroy: func(connectionType model.ConnectionType, stackData common2.IStackCreateData) error {
			if _, ok := stackData.(*StackData); !ok {
				return WrongStackDataError(connectionType, stackData)
			}
			if closer, ok := stackData.(io.Closer); ok {
				return closer.Close()
			}
			return nil
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
			sd, ok := stackData.(*StackData)
			if !ok {
				return nil, WrongStackDataError(connectionType, stackData)
			}
			return sd.OnStart(conn, Url, Ctx, ConnectionReactorFactory)
		},
		Stop: func(stackData interface{}, endParams common2.StackEndStateParams) error {
			_, ok := stackData.(*StackData)
			if !ok {
				return WrongStackDataError(connectionType, stackData)
			}
			var err error = nil
			return err
		},
	}
}
