package internal

import (
	"github.com/bhbosman/goCommsStacks/websocket/common"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Inbound(connectionType model.ConnectionType, ConnectionCancelFunc model.ConnectionCancelFunc, logger *zap.Logger, opts ...rxgo.Option) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(stackData common2.IStackCreateData, pipeData common2.IPipeCreateData, obs rxgo.Observable) (string, rxgo.Observable, error) {
					if sd, ok := stackData.(*StackData); ok {
						nextInBoundChannelManager := make(chan rxgo.Item)

						createSendData, err := RxHandlers.CreateSendData(nextInBoundChannelManager, logger)
						if err != nil {
							return "", nil, err
						}

						err = sd.setOnInBoundSendData(createSendData)
						if err != nil {
							return "", nil, err
						}

						createSendError, err := RxHandlers.CreateSendError(nextInBoundChannelManager, logger)
						if err != nil {
							return "", nil, err
						}

						err = sd.setOnInBoundSendError(createSendError)
						if err != nil {
							return "", nil, err
						}

						createComplete, err := RxHandlers.CreateComplete(nextInBoundChannelManager, logger)
						if err != nil {
							return "", nil, err
						}

						err = sd.SetOnInBoundComplete(createComplete)
						if err != nil {
							return "", nil, err
						}

						inboundStateParams, err := NewInboundStackHandler(sd)
						if err != nil {
							return "", nil, err
						}

						nextHandler, err := RxHandlers.NewRxNextHandler(
							common.StackName,
							ConnectionCancelFunc,
							inboundStateParams,
							sd.onInBoundSendData,
							sd.onInBoundSendError,
							sd.onInBoundComplete,
							logger)
						if err != nil {
							return "", nil, err
						}

						_ = obs.ForEach(nextHandler.OnSendData, nextHandler.OnError, nextHandler.OnComplete, opts...)
						return common.StackName, rxgo.FromChannel(nextInBoundChannelManager), nil
					}
					return "", nil, WrongStackDataError(connectionType, stackData)
				},
				nil),
			nil
	}
}
