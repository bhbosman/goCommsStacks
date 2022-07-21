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
	connectionType model.ConnectionType,
	ConnectionCancelFunc model.ConnectionCancelFunc,
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
					if sd, ok := stackData.(*data); ok {
						nextInBoundChannelManager := make(chan rxgo.Item)
						var err error
						sd.inboundHandler, err = RxHandlers.All2(
							goCommsDefinitions.WebSocketStackName,
							model.StreamDirectionUnknown,
							nextInBoundChannelManager,
							logger,
							ctx,
							true,
						)
						if err != nil {
							return nil, err
						}

						inboundStateParams, err := NewInboundStackHandler(sd)
						if err != nil {
							return nil, err
						}

						nextHandler, err := RxHandlers.NewRxNextHandler2(
							goCommsDefinitions.WebSocketStackName,
							ConnectionCancelFunc,
							inboundStateParams,
							sd.inboundHandler,
							logger)
						if err != nil {
							return nil, err
						}

						_ = rxOverride.ForEach2(
							goCommsDefinitions.WebSocketStackName,
							model.StreamDirectionInbound,
							obs,
							ctx,
							goFunctionCounter,
							nextHandler,
							opts...)
						return rxgo.FromChannel(nextInBoundChannelManager), nil
					}
					return nil, WrongStackDataError(connectionType, stackData)
				},
				nil),
			nil
	}
}
