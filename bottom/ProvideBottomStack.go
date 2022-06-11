package bottom

import (
	"github.com/bhbosman/goCommsStacks/bottom/common"
	internal2 "github.com/bhbosman/goCommsStacks/bottom/internal"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideBottomStack() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: "StackDefinition",
			Target: func(params struct {
				fx.In
				ConnectionCancelFunc model.ConnectionCancelFunc
				Opts                 []rxgo.Option
				Logger               *zap.Logger
			}) (*internalComms.StackDefinition, error) {
				return &internalComms.StackDefinition{
					Name: common.StackName,
					Inbound: internalComms.NewBoundResultImpl(
						internal2.Inbound(
							params.ConnectionCancelFunc,
							params.Logger,
							params.Opts...)),
					Outbound: internalComms.NewBoundResultImpl(
						internal2.Outbound(
							params.ConnectionCancelFunc,
							params.Logger,
							params.Opts...)),
				}, nil
			},
		},
	)
}
