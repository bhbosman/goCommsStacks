package internal

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Inbound(
	ConnectionCancelFunc model.ConnectionCancelFunc,
	marker [4]byte,
	logger *zap.Logger,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	opts ...rxgo.Option,
) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(
					stackData common2.IStackCreateData,
					pipeData common2.IPipeCreateData,
					obs rxgo.Observable,
				) (rxgo.Observable, error) {
					channelManager := make(chan rxgo.Item)
					inboundStateParams := NewInboundStackHandler(marker)

					handler, err := RxHandlers.All2(
						goCommsDefinitions.MessageBreakerStackName,
						model.StreamDirectionUnknown,
						channelManager,
						logger,
						ctx,
						true,
					)
					if err != nil {
						return nil, err
					}
					nextHandler, err := RxHandlers.NewRxNextHandler2(
						goCommsDefinitions.MessageBreakerStackName,
						ConnectionCancelFunc,
						inboundStateParams,
						handler,
						logger,
					)
					if err != nil {
						return nil, err
					}

					_ = rxOverride.ForEach2(
						goCommsDefinitions.MessageBreakerStackName,
						model.StreamDirectionInbound,
						obs,
						ctx,
						goFunctionCounter,
						nextHandler,
						opts...,
					)
					return rxgo.FromChannel(
						channelManager,
						opts...,
					), nil
				}, nil),
			nil
	}
}
