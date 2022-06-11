package tlsConnection

import (
	"github.com/bhbosman/goCommsStacks/tlsConnection/common"
	internal2 "github.com/bhbosman/goCommsStacks/tlsConnection/internal"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideTlsConnectionStacks() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "StackDefinition",
				Target: func(params struct {
					fx.In
					ConnectionCancelFunc model.ConnectionCancelFunc
					Opts                 []rxgo.Option
					ConnectionType       model.ConnectionType
					Logger               *zap.Logger
				}) (*internalComms.StackDefinition, error) {
					return &internalComms.StackDefinition{
						Name: common.StackName,
						Inbound: internalComms.NewBoundResultImpl(internal2.Inbound(
							params.ConnectionCancelFunc,
							params.ConnectionType,
							params.Logger,
							params.Opts...)),
						Outbound: internalComms.NewBoundResultImpl(internal2.Outbound(
							params.ConnectionCancelFunc,
							params.ConnectionType,
							params.Logger,
							params.Opts...)),
						StackState: internal2.CreateStackState(params.ConnectionType),
					}, nil
				},
			}))
}
