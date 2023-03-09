package websocket

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/websocket/internal"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/url"
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
						ConnectionId         string `name:"ConnectionId"`
						Opts                 []rxgo.Option
						ConnectionType       model.ConnectionType
						Logger               *zap.Logger
						ExtractValues        intf.IConnectionReactorFactoryExtractValues `optional:"true"`
						ConnectionReactor    intf.IConnectionReactor
						Conn                 net.Conn `name:"PrimaryConnection"`
						Ctx                  context.Context
						CancelFunc           context.CancelFunc
						Url                  *url.URL
						GoFunctionCounter    GoFunctionCounter.IService
					},
				) (common.IStackDefinition, error) {
					if params.ExtractValues == nil {
						params.ExtractValues = intf.NewDefaultConnectionReactorFactoryExtractValues()
					}
					return common.NewStackDefinition(
						goCommsDefinitions.WebSocketStackName,
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(
										stackData common.IStackCreateData,
										pipeData common.IPipeCreateData,
										obs gocommon.IObservable,
									) (gocommon.IObservable, error) {
										if sd, ok := stackData.(*internal.Data); ok {
											nextInBoundChannelManager := make(chan rxgo.Item)
											var err error
											sd.InboundHandler, err = RxHandlers.All2(
												goCommsDefinitions.WebSocketStackName,
												model.StreamDirectionUnknown,
												nextInBoundChannelManager,
												params.Logger,
												params.Ctx,
												true,
											)
											if err != nil {
												return nil, err
											}

											inboundStateParams, err := internal.NewInboundStackHandler(sd)
											if err != nil {
												return nil, err
											}

											nextHandler, err := RxHandlers.NewRxNextHandler2(
												goCommsDefinitions.WebSocketStackName,
												params.ConnectionCancelFunc,
												inboundStateParams,
												sd.InboundHandler,
												params.Logger)
											if err != nil {
												return nil, err
											}

											_ = rxOverride.ForEach2(
												goCommsDefinitions.WebSocketStackName,
												model.StreamDirectionInbound,
												obs,
												params.Ctx,
												params.GoFunctionCounter,
												nextHandler,
												params.Opts...)
											return rxgo.FromChannel(nextInBoundChannelManager), nil
										}
										return nil, internal.WrongStackDataError(params.ConnectionType, stackData)
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
										if sd, ok := stackData.(*internal.Data); ok {
											MultiReceiverChannel := make(chan rxgo.Item)
											var err error
											sd.MultiOutBoundHandler, err = RxHandlers.All2(
												goCommsDefinitions.WebSocketStackName,
												model.StreamDirectionUnknown,
												MultiReceiverChannel,
												params.Logger,
												params.Ctx,
												true,
											)

											if err != nil {
												return nil, err
											}

											outboundStackHandler001 := internal.NewOutboundStackHandler001(sd)
											nextHandler001, err := RxHandlers.NewRxNextHandler2(
												fmt.Sprintf("%v-01", goCommsDefinitions.WebSocketStackName),
												params.ConnectionCancelFunc,
												outboundStackHandler001,
												sd.MultiOutBoundHandler,
												params.Logger)
											if err != nil {
												return nil, err
											}

											_ = rxOverride.ForEach2(
												goCommsDefinitions.WebSocketStackName,
												model.StreamDirectionOutbound,
												obs,
												params.Ctx,
												params.GoFunctionCounter,
												nextHandler001,
												params.Opts...)

											NextOutboundChannelManager := make(chan rxgo.Item)
											sd.OutboundHandler, err = RxHandlers.All2(
												goCommsDefinitions.WebSocketStackName,
												model.StreamDirectionUnknown,
												NextOutboundChannelManager,
												params.Logger,
												params.Ctx,
												true,
											)
											if err != nil {
												return nil, err
											}

											connWrapper, err := common.NewConnWrapper(
												sd.Conn,
												sd.Ctx,
												sd.PipeRead,
												sd.OutboundHandler.OnSendData,
												//NextOutboundChannelManager,
											)
											if err != nil {
												return nil, err
											}

											err = sd.SetConnWrapper(connWrapper)
											if err != nil {
												return nil, err
											}

											outboundStackHandler002, err := internal.NewOutboundStackHandler002(sd)
											if err != nil {
												return nil, err
											}

											nextHandler002, err := RxHandlers.NewRxNextHandler2(
												fmt.Sprintf("%v-02", goCommsDefinitions.WebSocketStackName),
												params.ConnectionCancelFunc,
												outboundStackHandler002,
												sd.OutboundHandler,
												params.Logger)
											if err != nil {
												return nil, err
											}

											tempObs := rxgo.FromChannel(MultiReceiverChannel)
											_ = rxOverride.ForEach2(
												goCommsDefinitions.WebSocketStackName,
												model.StreamDirectionOutbound,
												tempObs,
												params.Ctx,
												params.GoFunctionCounter,
												nextHandler002,
												params.Opts...,
											)

											return rxgo.FromChannel(NextOutboundChannelManager), nil
										}
										return nil, internal.WrongStackDataError(params.ConnectionType, stackData)
									},
									nil),
								nil
						},
						internal.CreateStackState(
							params.ConnectionType,
							params.ExtractValues,
							params.Logger,
							params.Conn,
							params.Ctx,
							params.CancelFunc,
							params.Url,
							params.GoFunctionCounter,
						))
				},
			},
		),
	)
}
