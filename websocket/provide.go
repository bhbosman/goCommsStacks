package websocket

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	internal2 "github.com/bhbosman/goCommsStacks/websocket/internal"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/url"
)

func ProvideWebsocketStacks() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "StackDefinition",
				Target: func(
					params struct {
						fx.In
						ConnectionCancelFunc model.ConnectionCancelFunc
						ConnectionId         string `name:"ConnectionId"`
						Opts                 []rxgo.Option
						ConnectionType       model.ConnectionType
						Logger               *zap.Logger
						ExtractValues        intf.IConnectionReactorFactoryExtractValues `optional:"true"`
						ConnectionReactor    intf.IConnectionReactor
						Conn                 net.Conn
						Ctx                  context.Context
						CancelFunc           context.CancelFunc
						Url                  *url.URL
						GoFunctionCounter    GoFunctionCounter.IService
					},
				) (*internalComms.StackDefinition, error) {
					if params.ExtractValues == nil {
						params.ExtractValues = intf.NewDefaultConnectionReactorFactoryExtractValues()
					}
					return &internalComms.StackDefinition{
						Name: goCommsDefinitions.WebSocketStackName,
						Inbound: internalComms.NewBoundResultImpl(
							internal2.Inbound(
								params.ConnectionType,
								params.ConnectionCancelFunc,
								params.Logger,
								params.Ctx,
								params.GoFunctionCounter,
								params.Opts...,
							),
						),
						Outbound: internalComms.NewBoundResultImpl(
							internal2.Outbound(
								params.ConnectionType,
								params.ConnectionCancelFunc,
								params.Logger,
								params.Ctx,
								params.GoFunctionCounter,
								params.Opts...,
							),
						),
						StackState: internal2.CreateStackState(
							params.ConnectionType,
							params.ConnectionCancelFunc,
							params.ExtractValues,
							params.Logger,
							params.Conn,
							params.Ctx,
							params.CancelFunc,
							params.Url,
							params.GoFunctionCounter,
						),
					}, nil
				},
			},
		),
	)
}
