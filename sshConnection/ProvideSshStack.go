package sshConnection

import (
	"github.com/bhbosman/goCommsStacks/sshConnection/common"
	internal2 "github.com/bhbosman/goCommsStacks/sshConnection/internal"
	"github.com/bhbosman/gocommon/model"
	internalComms "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func ProvideSshCommunicationStack() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group: "StackDefinition",
			Target: func(params struct {
				fx.In
				ConnectionCancelFunc model.ConnectionCancelFunc
				Opts                 []rxgo.Option
				Logger               *zap.Logger
				ConnectionType       model.ConnectionType
			}) (*internalComms.StackDefinition, error) {
				var errList error = nil
				if params.ConnectionCancelFunc == nil {
					errList = multierr.Append(errList, goerrors.InvalidParam)
				}

				if params.Logger == nil {
					errList = multierr.Append(errList, goerrors.InvalidParam)
				}

				if errList != nil {
					return nil, errList
				}

				return &internalComms.StackDefinition{
					Name: common.StackName,
					Inbound: internalComms.NewBoundResultImpl(
						internal2.Inbound(
							params.ConnectionType,
							params.ConnectionCancelFunc,
							params.Logger,
							params.Opts...)),
					Outbound: internalComms.NewBoundResultImpl(
						internal2.Outbound(
							params.ConnectionType,
							params.ConnectionCancelFunc,
							params.Logger,
							params.Opts...)),
					StackState: internal2.CreateStackState(params.ConnectionType),
				}, nil
			},
		},
	)
}
