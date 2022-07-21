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

func outbound(
	ConnectionCancelFunc model.ConnectionCancelFunc,
	connectionType model.ConnectionType,
	logger *zap.Logger,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	opts ...rxgo.Option,
) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(
					sd common2.IStackCreateData,
					pipeData common2.IPipeCreateData,
					obs rxgo.Observable,
				) (rxgo.Observable, error) {
					if stackData, ok := sd.(*data); ok {
						nextOutboundChannel := make(chan rxgo.Item)
						var err error
						stackData.outboundHandler, err = RxHandlers.All2(
							goCommsDefinitions.TlsStackName,
							model.StreamDirectionUnknown,
							nextOutboundChannel,
							logger,
							ctx,
							true,
						)
						if err != nil {
							return nil, err
						}

						connWrapper, err := common2.NewConnWrapper(
							stackData.Conn,
							stackData.ctx,
							stackData.PipeRead,
							stackData.outboundHandler.OnSendData,
						)
						if err != nil {
							return nil, err
						}

						err = stackData.setConnWrapper(connWrapper)
						if err != nil {
							return nil, err
						}

						outboundStackHandler, err := NewOutboundStackHandler(stackData)
						if err != nil {
							return nil, err
						}

						nextHandler, err := RxHandlers.NewRxNextHandler2(
							goCommsDefinitions.TlsStackName,
							ConnectionCancelFunc,
							outboundStackHandler,
							stackData.outboundHandler,
							logger)
						if err != nil {
							return nil, err
						}

						rxOverride.ForEach2(
							goCommsDefinitions.TlsStackName,
							model.StreamDirectionOutbound,
							obs,
							ctx,
							goFunctionCounter,
							nextHandler,
							opts...)
						nextObs := rxgo.FromChannel(nextOutboundChannel, opts...)
						return nextObs, nil
					}
					return nil, WrongStackDataError(connectionType, sd)
				},
				nil),
			nil
	}
}
