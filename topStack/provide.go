package topStack

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/model"
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
						goCommsDefinitions.TopStackName,
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
								func(
									stackData common.IStackCreateData,
									pipeData common.IPipeCreateData,
									obs gocommon.IObservable,
								) (gocommon.IObservable, error) {
									return obs, nil
								},
								nil), nil
						},
						func() (common.IStackBoundFactory, error) {
							return common.NewStackBoundDefinition(
								func(
									stackData common.IStackCreateData,
									pipeData common.IPipeCreateData,
									obs gocommon.IObservable,
								) (gocommon.IObservable, error) {
									return obs, nil
								},
								nil,
							), nil
						},
						nil,
					)
				},
			},
		),
	)
}
