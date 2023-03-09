package pingPong

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/internal/pingPong"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
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
					return common.NewStackDefinition(
						goCommsDefinitions.PingPongStackName,
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(
										stackData common.IStackCreateData,
										pipeData common.IPipeCreateData,
										obs gocommon.IObservable,
									) (gocommon.IObservable, error) {
										if stackDataActual, ok := stackData.(*pingPong.Data); ok {
											handler, err := RxHandlers.NewRxMapHandler(
												goCommsDefinitions.PingPongStackName,
												params.ConnectionCancelFunc,
												params.Logger,
												&pingPong.InboundStackHandler{
													StackData: stackDataActual,
												},
											)
											if err != nil {
												return nil, err
											}
											return obs.FlatMap(handler.FlatMapHandler(params.Ctx), params.Opts...), nil
										}
										return nil, goerrors.InvalidType
									},
									nil),
								nil
						},
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(
										stackData common.IStackCreateData,
										pipeData common.IPipeCreateData,
										obs gocommon.IObservable,
									) (gocommon.IObservable, error) {
										if sd, ok := stackData.(*pingPong.Data); ok {
											outboundChannel := make(chan rxgo.Item)
											var err error
											sd.OutHandler, err = RxHandlers.All2(
												goCommsDefinitions.PingPongStackName,
												model.StreamDirectionUnknown,
												outboundChannel,
												params.Logger,
												params.Ctx,
												true,
											)
											if err != nil {
												return nil, err
											}

											outboundStackHandlerInstance, err := pingPong.NewOutboundStackHandler(sd)
											if err != nil {
												return nil, err
											}

											rxNextHandler, err := RxHandlers.NewRxNextHandler2(
												goCommsDefinitions.PingPongStackName,
												params.ConnectionCancelFunc,
												outboundStackHandlerInstance,
												sd.OutHandler,
												params.Logger)
											if err != nil {
												return nil, err
											}

											_ = rxOverride.ForEachWithTimeChannel(
												goCommsDefinitions.PingPongStackName,
												model.StreamDirectionOutbound,
												obs,
												params.Ctx,
												params.GoFunctionCounter,
												rxNextHandler,
												sd.Ticker.C,
												func(time2 time.Time) interface{} {
													rxNextHandler.OtherMessageCountOut++
													return sd.SendPing(time2)
												},
												params.Opts...)
											resultObs := rxgo.FromChannel(outboundChannel, params.Opts...)
											return resultObs, nil
										}
										return nil, goerrors.InvalidType
									},
									nil),
								nil
						},
						pingPong.CreateStackState(
							params.Ctx,
							params.GoFunctionCounter,
						),
					)
				},
			},
		),
	)
}
