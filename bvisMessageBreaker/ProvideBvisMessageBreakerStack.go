package bvisMessageBreaker

import (
	"github.com/bhbosman/goCommsStacks/bvisMessageBreaker/common"
	internal2 "github.com/bhbosman/goCommsStacks/bvisMessageBreaker/internal"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideBvisMessageBreakerStack() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotated{
			Group: "StackDefinition",
			Target: func(params struct {
				fx.In
				ConnectionCancelFunc model.ConnectionCancelFunc
				Opts                 []rxgo.Option
				Logger               *zap.Logger
			}) (*internalComms.StackDefinition, error) {
				marker := [4]byte{'B', 'V', 'I', 'S'}
				return &internalComms.StackDefinition{
					Name:     common.StackName,
					Inbound:  internalComms.NewBoundResultImpl(internal2.Inbound(params.ConnectionCancelFunc, marker, params.Logger, params.Opts...)),
					Outbound: internalComms.NewBoundResultImpl(internal2.Outbound(params.ConnectionCancelFunc, marker, params.Logger, params.Opts...)),
				}, nil
			},
		}))
}
