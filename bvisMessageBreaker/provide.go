package bvisMessageBreaker

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/internal/bvisMessageBreaker"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide() fx.Option {
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
				) (common.IStackDefinition, error) {
					marker := [4]byte{'B', 'V', 'I', 'S'}
					return common.NewStackDefinition(
						goCommsDefinitions.MessageBreakerStackName,
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(stackData common.IStackCreateData, pipeData common.IPipeCreateData, obs gocommon.IObservable) (gocommon.IObservable, error) {
										inboundStackHandler, err := bvisMessageBreaker.NewRxMapStackHandler(marker)
										if err != nil {
											return nil, err
										}

										mapHandler, err := RxHandlers.NewRxMapHandler(
											goCommsDefinitions.MessageBreakerStackName,
											params.ConnectionCancelFunc,
											params.Logger,
											inboundStackHandler,
										)
										if err != nil {
											return nil, err
										}

										return obs.FlatMap(mapHandler.FlatMapHandler(params.Ctx), params.Opts...), nil

									}, nil),
								nil
						},
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(stackData common.IStackCreateData, pipeData common.IPipeCreateData, obs gocommon.IObservable) (gocommon.IObservable, error) {
										outboundStackHandler, err := bvisMessageBreaker.NewOutboundStackHandler(marker)
										if err != nil {
											return nil, err
										}

										mapHandler, err := RxHandlers.NewRxMapHandler(
											goCommsDefinitions.MessageBreakerStackName,
											params.ConnectionCancelFunc,
											params.Logger,
											outboundStackHandler,
										)
										if err != nil {
											return nil, err
										}

										return obs.Map(mapHandler.Handler, params.Opts...), nil
									},
									nil,
								),
								nil
						},
						nil,
					)
				},
			},
		),
	)
}
