package tlsStack

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/tlsStack/internal"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net"
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
						ConnectionType       model.ConnectionType
						Logger               *zap.Logger
						Conn                 net.Conn `name:"PrimaryConnection"`
						Ctx                  context.Context
						CtxCancelFunc        context.CancelFunc
						CancelFunc           model.ConnectionCancelFunc
						GoFunctionCounter    GoFunctionCounter.IService
					},
				) (common.IStackDefinition, error) {
					return common.NewStackDefinition(
						goCommsDefinitions.TlsStackName,
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(
										stackData common.IStackCreateData,
										pipeData common.IPipeCreateData,
										obs gocommon.IObservable,
									) (gocommon.IObservable, error) {
										if sd, ok := stackData.(*internal.Data); ok {
											nextChannel := make(chan rxgo.Item)
											var err error
											// do not assign sd.inboundHandler.onInBoundComplete to onComplete as this will do a double close
											// the double close may be handled
											// make sure useCompleteCallback = false
											sd.InboundHandler, err = RxHandlers.All2(
												goCommsDefinitions.TlsStackName,
												model.StreamDirectionUnknown,
												nextChannel,
												params.Logger,
												params.Ctx,
												false,
											)
											if err != nil {
												return nil, err
											}

											inboundStackHandlerInstance, err := internal.NewInboundStackHandler(
												sd,
												params.Ctx,
												params.CtxCancelFunc,
											)
											if err != nil {
												return nil, err
											}

											nextHandler, err := RxHandlers.NewRxNextHandler2(
												goCommsDefinitions.TlsStackName,
												params.ConnectionCancelFunc,
												inboundStackHandlerInstance,
												sd.InboundHandler,
												params.Logger)
											if err != nil {
												return nil, err
											}

											rxOverride.ForEach2(
												goCommsDefinitions.TlsStackName,
												model.StreamDirectionInbound,
												obs,
												params.Ctx,
												params.GoFunctionCounter,
												nextHandler,
												params.Opts...)
											nextObs := rxgo.FromChannel(nextChannel, params.Opts...)
											return nextObs, nil
										}
										return nil, internal.WrongStackDataError(params.ConnectionType, stackData)
									},
									nil),
								nil
						},
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(
										sd common.IStackCreateData,
										pipeData common.IPipeCreateData,
										obs gocommon.IObservable,
									) (gocommon.IObservable, error) {
										if stackData, ok := sd.(*internal.Data); ok {
											nextOutboundChannel := make(chan rxgo.Item)
											var err error
											stackData.OutboundHandler, err = RxHandlers.All2(
												goCommsDefinitions.TlsStackName,
												model.StreamDirectionUnknown,
												nextOutboundChannel,
												params.Logger,
												params.Ctx,
												true,
											)
											if err != nil {
												return nil, err
											}

											connWrapper, err := common.NewConnWrapper(
												stackData.Conn,
												stackData.Ctx,
												stackData.PipeRead,
												stackData.OutboundHandler.OnSendData,
											)
											if err != nil {
												return nil, err
											}

											err = stackData.SetConnWrapper(connWrapper)
											if err != nil {
												return nil, err
											}

											outboundStackHandler, err := internal.NewOutboundStackHandler(stackData)
											if err != nil {
												return nil, err
											}

											nextHandler, err := RxHandlers.NewRxNextHandler2(
												goCommsDefinitions.TlsStackName,
												params.ConnectionCancelFunc,
												outboundStackHandler,
												stackData.OutboundHandler,
												params.Logger)
											if err != nil {
												return nil, err
											}

											rxOverride.ForEach2(
												goCommsDefinitions.TlsStackName,
												model.StreamDirectionOutbound,
												obs,
												params.Ctx,
												params.GoFunctionCounter,
												nextHandler,
												params.Opts...)
											nextObs := rxgo.FromChannel(nextOutboundChannel, params.Opts...)
											return nextObs, nil
										}
										return nil, internal.WrongStackDataError(params.ConnectionType, sd)
									},
									nil),
								nil
						},
						internal.CreateStackState(
							params.Logger,
							params.ConnectionType,
							params.Conn,
							params.Ctx,
							params.CtxCancelFunc,
							params.CancelFunc,
							params.GoFunctionCounter,
						))

				},
			},
		),
	)
}
