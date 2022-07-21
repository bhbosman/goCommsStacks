package protoBuf

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/protoBuf/internal"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideStack() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: "StackDefinition",
			Target: func(
				params struct {
					fx.In
					ConnectionCancelFunc model.ConnectionCancelFunc
					Opts                 []rxgo.Option
					Logger               *zap.Logger
				},
			) (*common.StackDefinition, error) {
				return &common.StackDefinition{
					Name: goCommsDefinitions.ProtoBufferStackName,
					Inbound: common.NewBoundResultImpl(
						internal.Inbound(
							params.ConnectionCancelFunc,
							params.Logger,
							params.Opts...)),
					Outbound: common.NewBoundResultImpl(
						internal.Outbound(
							params.ConnectionCancelFunc,
							params.Logger,
							params.Opts...)),
				}, nil
			},
		},
	)
}
