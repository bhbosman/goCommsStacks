package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Outbound(
	connectionType model.ConnectionType,
	ConnectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	ctx context.Context,
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
						MultiReceiverChannel := make(chan rxgo.Item)
						var err error
						sd.multiOutBoundHandler, err = RxHandlers.All2(
							goCommsDefinitions.WebSocketStackName,
							model.StreamDirectionUnknown,
							MultiReceiverChannel,
							logger,
							ctx,
							true,
						)

						if err != nil {
							return nil, err
						}

						outboundStackHandler001 := NewOutboundStackHandler001(sd)
						nextHandler001, err := RxHandlers.NewRxNextHandler2(
							fmt.Sprintf("%v-01", goCommsDefinitions.WebSocketStackName),
							ConnectionCancelFunc,
							outboundStackHandler001,
							sd.multiOutBoundHandler,
							logger)
						if err != nil {
							return nil, err
						}

						_ = rxOverride.ForEach2(
							goCommsDefinitions.WebSocketStackName,
							model.StreamDirectionOutbound,
							obs,
							ctx,
							goFunctionCounter,
							nextHandler001,
							opts...)

						NextOutboundChannelManager := make(chan rxgo.Item)
						sd.outboundHandler, err = RxHandlers.All2(
							goCommsDefinitions.WebSocketStackName,
							model.StreamDirectionUnknown,
							NextOutboundChannelManager,
							logger,
							ctx,
							true,
						)
						if err != nil {
							return nil, err
						}

						connWrapper, err := common2.NewConnWrapper(
							sd.conn,
							sd.ctx,
							sd.pipeRead,
							sd.outboundHandler.OnSendData,
							//NextOutboundChannelManager,
						)
						if err != nil {
							return nil, err
						}

						err = sd.setConnWrapper(connWrapper)
						if err != nil {
							return nil, err
						}

						outboundStackHandler002, err := NewOutboundStackHandler002(sd)
						if err != nil {
							return nil, err
						}

						nextHandler002, err := RxHandlers.NewRxNextHandler2(
							fmt.Sprintf("%v-02", goCommsDefinitions.WebSocketStackName),
							ConnectionCancelFunc,
							outboundStackHandler002,
							sd.outboundHandler,
							logger)
						if err != nil {
							return nil, err
						}

						tempObs := rxgo.FromChannel(MultiReceiverChannel)
						_ = rxOverride.ForEach2(
							goCommsDefinitions.WebSocketStackName,
							model.StreamDirectionOutbound,
							tempObs,
							ctx,
							goFunctionCounter,
							nextHandler002,
							opts...,
						)

						return rxgo.FromChannel(NextOutboundChannelManager), nil
					}
					return nil, WrongStackDataError(connectionType, stackData)
				},
				nil),
			nil
	}
}
