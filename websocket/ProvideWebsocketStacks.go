package websocket

import (
	"github.com/bhbosman/goCommsStacks/websocket/common"
	internal2 "github.com/bhbosman/goCommsStacks/websocket/internal"
	"github.com/bhbosman/gocommon/Services/IConnectionManager"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideWebsocketStacks() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "StackDefinition",
				Target: func(params struct {
					fx.In
					Conn                 IConnectionManager.IPublishConnectionInformation
					ConnectionCancelFunc model.ConnectionCancelFunc
					ConnectionId         string `name:"ConnectionId"`
					Opts                 []rxgo.Option
					ConnectionType       model.ConnectionType
					Logger               *zap.Logger
				}) (*internalComms.StackDefinition, error) {
					return &internalComms.StackDefinition{
						Name:       common.StackName,
						Inbound:    internalComms.NewBoundResultImpl(internal2.Inbound(params.ConnectionType, params.ConnectionCancelFunc, params.Logger, params.Opts...)),
						Outbound:   internalComms.NewBoundResultImpl(internal2.Outbound(params.ConnectionType, params.ConnectionCancelFunc, params.Logger, params.Opts...)),
						StackState: internal2.CreateStackState(params.ConnectionType, params.ConnectionCancelFunc),
					}, nil
				},
			}))
}
