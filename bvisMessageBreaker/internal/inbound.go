package internal

import (
	"github.com/bhbosman/goCommsStacks/bvisMessageBreaker/common"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Inbound(ConnectionCancelFunc model.ConnectionCancelFunc, marker [4]byte, logger *zap.Logger, opts ...rxgo.Option) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(stackData common2.IStackCreateData, pipeData common2.IPipeCreateData, obs rxgo.Observable) (string, rxgo.Observable, error) {
					channelManager := make(chan rxgo.Item)
					inboundStateParams := NewInboundStackHandler(marker)
					createComplete, err := RxHandlers.CreateComplete(channelManager, logger)
					if err != nil {
						return "", nil, err
					}
					createSendError, err := RxHandlers.CreateSendError(channelManager, logger)
					if err != nil {
						return "", nil, err
					}
					createSendData, err := RxHandlers.CreateSendData(channelManager, logger)
					if err != nil {
						return "", nil, err
					}
					nextHandler, err := RxHandlers.NewRxNextHandler(
						common.StackName,
						ConnectionCancelFunc,
						inboundStateParams,
						createSendData,
						createSendError,
						createComplete,
						logger)
					if err != nil {
						return "", nil, err
					}

					_ = obs.ForEach(
						nextHandler.OnSendData,
						nextHandler.OnError,
						nextHandler.OnComplete,
						opts...)
					return common.StackName, rxgo.FromChannel(channelManager, opts...), nil
				}, nil),
			nil
	}
}
