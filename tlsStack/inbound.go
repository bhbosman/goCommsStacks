package tlsStack

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

func inbound(
	ConnectionCancelFunc model.ConnectionCancelFunc,
	connectionType model.ConnectionType,
	logger *zap.Logger,
	ctx context.Context,
	cancelFunc context.CancelFunc,
	goFunctionCounter GoFunctionCounter.IService,
	opts ...rxgo.Option) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(
					stackData common2.IStackCreateData,
					pipeData common2.IPipeCreateData,
					obs rxgo.Observable,
				) (rxgo.Observable, error) {
					if sd, ok := stackData.(*data); ok {
						NextInBoundChannel := make(chan rxgo.Item)
						var err error
						// do not assign sd.inboundHandler.onInBoundComplete to onComplete as this will do a double close
						// the double close may be handled
						// make sure useCompleteCallback = false
						sd.inboundHandler, err = RxHandlers.All2(
							goCommsDefinitions.TlsStackName,
							model.StreamDirectionUnknown,
							NextInBoundChannel,
							logger,
							ctx,
							false,
						)
						if err != nil {
							return nil, err
						}

						inboundStackHandlerInstance, err := NewInboundStackHandler(
							sd,
							ctx,
							cancelFunc,
						)
						if err != nil {
							return nil, err
						}

						nextHandler, err := RxHandlers.NewRxNextHandler2(
							goCommsDefinitions.TlsStackName,
							ConnectionCancelFunc,
							inboundStackHandlerInstance,
							sd.inboundHandler,
							logger)
						if err != nil {
							return nil, err
						}

						rxOverride.ForEach2(
							goCommsDefinitions.TlsStackName,
							model.StreamDirectionInbound,
							obs,
							ctx,
							goFunctionCounter,
							nextHandler,
							opts...)
						nextObs := rxgo.FromChannel(NextInBoundChannel, opts...)
						return nextObs, nil
					}
					return nil, WrongStackDataError(connectionType, stackData)
				},
				nil),
			nil
	}
}
