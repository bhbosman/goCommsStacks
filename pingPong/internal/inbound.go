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

func Inbound(
	ConnectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	opts ...rxgo.Option,
) func() (common.IStackBoundDefinition, error) {
	return func() (common.IStackBoundDefinition, error) {
		return common.NewBoundDefinition(
				func(
					stackData common.IStackCreateData,
					pipeData common.IPipeCreateData,
					obs rxgo.Observable,
				) (rxgo.Observable, error) {
					if stackDataActual, ok := stackData.(*data); ok {
						inBoundChannel := make(chan rxgo.Item)
						inboundStackHandler := NewInboundStackHandler(stackDataActual)

						handler, err := RxHandlers.All2(
							goCommsDefinitions.PingPongStackName,
							model.StreamDirectionUnknown,
							inBoundChannel,
							logger,
							ctx,
							true,
						)
						if err != nil {
							return nil, err
						}

						rxNextHandler, err := RxHandlers.NewRxNextHandler2(
							goCommsDefinitions.PingPongStackName,
							ConnectionCancelFunc,
							inboundStackHandler,
							handler,
							logger)
						if err != nil {
							return nil, err
						}

						_ = rxOverride.ForEach2(
							goCommsDefinitions.PingPongStackName,
							model.StreamDirectionInbound,
							obs,
							ctx,
							goFunctionCounter,
							rxNextHandler,
							opts...)

						resultObs := rxgo.FromChannel(inBoundChannel, opts...)
						return resultObs, nil
					}
					return nil, goerrors.InvalidType
				},
				nil),
			nil
	}
}
