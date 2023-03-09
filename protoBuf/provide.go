package protoBuf

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/internal/protoBuf"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func Provide() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: "StackDefinition",
			Target: func(
				params struct {
					fx.In
					Context              context.Context
					ConnectionCancelFunc model.ConnectionCancelFunc
					Opts                 []rxgo.Option
					Logger               *zap.Logger
				},
			) (common.IStackDefinition, error) {
				return common.NewStackDefinition(
					goCommsDefinitions.ProtoBufferStackName,
					func() (common.IStackBoundFactory, error) {
						return common.NewStackBoundDefinition(
								func(
									stackData common.IStackCreateData,
									pipeData common.IPipeCreateData,
									obs gocommon.IObservable,
								) (gocommon.IObservable, error) {
									if pipeData != nil {
										return nil, goerrors.InvalidParam
									}
									if stackData != nil {
										return nil, goerrors.InvalidParam
									}
									var errorList error
									stackHandler, err := protoBuf.NewInboundStackHandler(params.Logger)
									if err != nil {
										errorList = multierr.Append(errorList, err)
									}
									mapHandler, err := RxHandlers.NewRxMapHandler(
										goCommsDefinitions.ProtoBufferStackName,
										params.ConnectionCancelFunc,
										params.Logger,
										stackHandler,
									)
									if err != nil {
										errorList = multierr.Append(errorList, err)
									}
									if errorList != nil {
										return nil, errorList
									}
									return obs.FlatMap(mapHandler.FlatMapHandler(params.Context), params.Opts...), nil
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
									if pipeData != nil {
										return nil, goerrors.InvalidParam
									}
									if stackData != nil {
										return nil, goerrors.InvalidParam
									}
									var errorList error
									stackHandler, err := protoBuf.NewOutboundStackHandler(params.Logger)
									if err != nil {
										errorList = multierr.Append(errorList, err)
									}
									mapHandler, err := RxHandlers.NewRxMapHandler(
										goCommsDefinitions.ProtoBufferStackName,
										params.ConnectionCancelFunc,
										params.Logger,
										stackHandler,
									)
									if err != nil {
										errorList = multierr.Append(errorList, err)
									}
									if errorList != nil {
										return nil, errorList
									}
									return obs.Map(
										mapHandler.Handler,
										params.Opts...,
									), nil
								},
								nil),
							nil
					},
					nil,
				)
			},
		},
	)
}
