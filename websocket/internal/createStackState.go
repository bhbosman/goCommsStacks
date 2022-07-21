package internal

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"io"
	"net"
	"net/url"
)

func CreateStackState(
	connectionType model.ConnectionType,
	stackCancelFunc model.ConnectionCancelFunc,
	extractValues intf.IConnectionReactorFactoryExtractValues,
	logger *zap.Logger,
	conn net.Conn,
	ctx context.Context,
	cancelFunc context.CancelFunc,
	Url *url.URL,
	goFunctionCounter GoFunctionCounter.IService,
) *common2.StackState {
	return &common2.StackState{
		Id:          goCommsDefinitions.WebSocketStackName,
		HijackStack: true,
		Create: func() (common2.IStackCreateData, error) {
			var err error
			var data *data
			data, err = NewStackData(
				logger,
				connectionType,
				conn,
				//stackCancelFunc,
				ctx,
				cancelFunc,
				extractValues,
				goFunctionCounter,
			)
			if err != nil {
				return nil, err
			}
			return data, err
		},
		Destroy: func(
			connectionType model.ConnectionType,
			stackData common2.IStackCreateData,
		) error {
			if _, ok := stackData.(*data); !ok {
				return WrongStackDataError(connectionType, stackData)
			}
			if closer, ok := stackData.(io.Closer); ok {
				return closer.Close()
			}
			return nil
		},
		Start: func(
			conn common2.IInputStreamForStack,
			stackData common2.IStackCreateData,
			ToReactorFunc rxgo.NextFunc,
		) (common2.IInputStreamForStack, error) {
			sd, ok := stackData.(*data)
			if !ok {
				return nil, WrongStackDataError(connectionType, stackData)
			}
			return sd.OnStart(conn, Url)
		},
		Stop: func(stackData interface{}, endParams common2.StackEndStateParams) error {
			_, ok := stackData.(*data)
			if !ok {
				return WrongStackDataError(connectionType, stackData)
			}
			var err error = nil
			return err
		},
	}
}
