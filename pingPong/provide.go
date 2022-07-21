package pingPong

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	internal2 "github.com/bhbosman/goCommsStacks/pingPong/internal"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvidePingPongStacks() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "StackDefinition",
				Target: func(
					params struct {
						fx.In
						ConnectionCancelFunc model.ConnectionCancelFunc
						Opts                 []rxgo.Option
						Logger               *zap.Logger
						Ctx                  context.Context
						GoFunctionCounter    GoFunctionCounter.IService
					},
				) (*internalComms.StackDefinition, error) {
					return &internalComms.StackDefinition{
						Name: goCommsDefinitions.PingPongStackName,
						Inbound: internalComms.NewBoundResultImpl(
							internal2.Inbound(
								params.ConnectionCancelFunc,
								params.Logger,
								params.Ctx,
								params.GoFunctionCounter,
								params.Opts...,
							),
						),
						Outbound: internalComms.NewBoundResultImpl(
							internal2.Outbound(
								params.ConnectionCancelFunc,
								params.Logger,
								params.Ctx,
								params.GoFunctionCounter,
								params.Opts...,
							),
						),
						StackState: internal2.CreateStackState(
							params.Ctx,
							params.GoFunctionCounter,
						),
					}, nil
				},
			},
		),
	)
}
