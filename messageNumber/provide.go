package messageNumber

import (
	"github.com/bhbosman/goCommsDefinitions"
	messageNumber2 "github.com/bhbosman/goCommsStacks/internal/messageNumber"
	"github.com/bhbosman/gocommon"
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
					},
				) (common.IStackDefinition, error) {
					return common.NewStackDefinition(
						goCommsDefinitions.MessageNumberStackName,
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
								func(_ common.IStackCreateData, _ common.IPipeCreateData, obs gocommon.IObservable) (gocommon.IObservable, error) {
									stackHandler, err := messageNumber2.NewInboundStackHandler()
									if err != nil {
										return nil, err
									}

									mapHandler, err := RxHandlers.NewRxMapHandler(
										goCommsDefinitions.MessageNumberStackName,
										params.ConnectionCancelFunc,
										params.Logger,
										stackHandler,
									)
									if err != nil {
										return nil, err
									}

									return obs.Map(mapHandler.Handler, params.Opts...), nil
								}, nil), nil
						},
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
									func(_ common.IStackCreateData, _ common.IPipeCreateData, obs gocommon.IObservable) (gocommon.IObservable, error) {
										stackHandler, err := messageNumber2.NewOutboundStackHandler()
										if err != nil {
											return nil, err
										}

										mapHandler, err := RxHandlers.NewRxMapHandler(
											goCommsDefinitions.MessageNumberStackName,
											params.ConnectionCancelFunc,
											params.Logger,
											stackHandler,
										)
										if err != nil {
											return nil, err
										}

										return obs.Map(mapHandler.Handler, params.Opts...), nil
									},
									nil),
								nil
						},
						nil,
					)
				},
			},
		),
	)
}
