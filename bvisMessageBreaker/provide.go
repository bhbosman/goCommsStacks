package bvisMessageBreaker

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	internal2 "github.com/bhbosman/goCommsStacks/bvisMessageBreaker/internal"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
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
				) (*internalComms.StackDefinition, error) {
					marker := [4]byte{'B', 'V', 'I', 'S'}
					return &internalComms.StackDefinition{
						Name: goCommsDefinitions.MessageBreakerStackName,
						Inbound: internalComms.NewBoundResultImpl(
							internal2.Inbound(
								params.ConnectionCancelFunc,
								marker,
								params.Logger,
								params.Ctx,
								params.GoFunctionCounter,
								params.Opts...,
							),
						),
						Outbound: internalComms.NewBoundResultImpl(
							internal2.Outbound(
								params.ConnectionCancelFunc,
								marker,
								params.Logger,
								params.Opts...,
							),
						),
					}, nil
				},
			},
		),
	)
}
