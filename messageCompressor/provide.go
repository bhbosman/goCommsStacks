package messageCompressor

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/internal/messageCompressor"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "StackDefinition",
				Target: func(params struct {
					fx.In
					ConnectionCancelFunc model.ConnectionCancelFunc
					Opts                 []rxgo.Option
					Logger               *zap.Logger
				},
				) (common.IStackDefinition, error) {
					return common.NewStackDefinition(
						goCommsDefinitions.CompressionStackName,
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(_ common.IStackCreateData, pipeData common.IPipeCreateData, obs gocommon.IObservable) (gocommon.IObservable, error) {
										if pd, ok := pipeData.(RxHandlers.IRxMapStackHandler); ok {
											mapHandler, err := RxHandlers.NewRxMapHandler(
												goCommsDefinitions.CompressionStackName,
												params.ConnectionCancelFunc,
												params.Logger,
												pd,
											)
											if err != nil {
												return nil, err
											}

											return obs.Map(mapHandler.Handler, params.Opts...), nil
										}
										return nil, goerrors.InvalidType
									},
									messageCompressor.InboundPipeState(goCommsDefinitions.CompressionStackName)),
								nil

						},
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(_ common.IStackCreateData, pipeData common.IPipeCreateData, obs gocommon.IObservable) (gocommon.IObservable, error) {
										if pd, ok := pipeData.(RxHandlers.IRxMapStackHandler); ok {
											mapHandler, err := RxHandlers.NewRxMapHandler(
												goCommsDefinitions.CompressionStackName,
												params.ConnectionCancelFunc,
												params.Logger,
												pd,
											)
											if err != nil {
												return nil, err
											}

											return obs.Map(mapHandler.Handler, params.Opts...), nil
										}
										return nil, goerrors.InvalidType
									},
									messageCompressor.OutboundPipeState(goCommsDefinitions.CompressionStackName)),
								nil
						},
						nil,
					)
				},
			},
		),
	)
}
