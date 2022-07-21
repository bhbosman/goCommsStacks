package internal

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Outbound(
	ConnectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	opts ...rxgo.Option) func() (common.IStackBoundDefinition, error) {
	return func() (common.IStackBoundDefinition, error) {
		return common.NewBoundDefinition(
				func(
					stackData common.IStackCreateData,
					pipeData common.IPipeCreateData,
					obs rxgo.Observable,
				) (rxgo.Observable, error) {
					if sd, ok := stackData.(*data); ok {
						outboundChannel := make(chan rxgo.Item)
						var err error
						sd.outHandler, err = RxHandlers.All2(
							goCommsDefinitions.PingPongStackName,
							model.StreamDirectionUnknown,
							outboundChannel,
							logger,
							ctx,
							true,
						)
						if err != nil {
							return nil, err
						}

						outboundStackHandler, err := NewOutboundStackHandler(sd)
						if err != nil {
							return nil, err
						}

						rxNextHandler, err := RxHandlers.NewRxNextHandler2(
							goCommsDefinitions.PingPongStackName,
							ConnectionCancelFunc,
							outboundStackHandler,
							sd.outHandler,
							logger)
						if err != nil {
							return nil, err
						}

						_ = rxOverride.ForEach2(
							goCommsDefinitions.PingPongStackName,
							model.StreamDirectionOutbound,
							obs,
							ctx,
							goFunctionCounter,
							rxNextHandler,
							opts...)
						resultObs := rxgo.FromChannel(outboundChannel, opts...)
						return resultObs, nil
					}
					return nil, goerrors.InvalidType
				},
				nil),
			nil
	}
}
