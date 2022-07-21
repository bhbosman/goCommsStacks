package tlsStack

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net"
)

func ProvideTlsConnectionStacks() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "StackDefinition",
				Target: func(
					params struct {
						fx.In
						ConnectionCancelFunc model.ConnectionCancelFunc
						Opts                 []rxgo.Option
						ConnectionType       model.ConnectionType
						Logger               *zap.Logger
						Conn                 net.Conn
						Ctx                  context.Context
						CtxCancelFunc        context.CancelFunc
						CancelFunc           model.ConnectionCancelFunc
						GoFunctionCounter    GoFunctionCounter.IService
					},
				) (*internalComms.StackDefinition, error) {
					return &internalComms.StackDefinition{
						Name: goCommsDefinitions.TlsStackName,
						Inbound: internalComms.NewBoundResultImpl(inbound(
							params.ConnectionCancelFunc,
							params.ConnectionType,
							params.Logger,
							params.Ctx,
							params.CtxCancelFunc,
							params.GoFunctionCounter,
							params.Opts...)),
						Outbound: internalComms.NewBoundResultImpl(outbound(
							params.ConnectionCancelFunc,
							params.ConnectionType,
							params.Logger,
							params.Ctx,
							params.GoFunctionCounter,
							params.Opts...)),
						StackState: createStackState(
							params.Logger,
							params.ConnectionType,
							params.Conn,
							params.Ctx,
							params.CtxCancelFunc,
							params.CancelFunc,
							params.GoFunctionCounter,
						),
					}, nil
				},
			},
		),
	)
}
