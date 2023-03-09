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
	extractValues intf.IConnectionReactorFactoryExtractValues,
	logger *zap.Logger,
	conn net.Conn,
	ctx context.Context,
	cancelFunc context.CancelFunc,
	Url *url.URL,
	goFunctionCounter GoFunctionCounter.IService,
) common2.IStackState {
	return common2.NewStackState(
		goCommsDefinitions.WebSocketStackName,
		true,
		func() (common2.IStackCreateData, error) {
			var err error
			var data *Data
			data, err = NewStackData(
				logger,
				connectionType,
				conn,
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
		func(
			connectionType model.ConnectionType,
			stackData common2.IStackCreateData,
		) error {
			if _, ok := stackData.(*Data); !ok {
				return WrongStackDataError(connectionType, stackData)
			}
			if closer, ok := stackData.(io.Closer); ok {
				return closer.Close()
			}
			return nil
		},
		func(
			conn common2.IInputStreamForStack,
			stackData common2.IStackCreateData,
			ToReactorFunc rxgo.NextFunc,
		) (common2.IInputStreamForStack, error) {
			sd, ok := stackData.(*Data)
			if !ok {
				return nil, WrongStackDataError(connectionType, stackData)
			}
			return sd.OnStart(conn, Url)
		},
		func(stackData interface{}) error {
			_, ok := stackData.(*Data)
			if !ok {
				return WrongStackDataError(connectionType, stackData)
			}
			var err error = nil
			return err
		},
	)
}
