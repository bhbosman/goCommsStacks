package internal

import (
	"github.com/bhbosman/goCommsStacks/sshConnection/common"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Inbound(connectionType model.ConnectionType, ConnectionCancelFunc model.ConnectionCancelFunc, logger *zap.Logger, opts ...rxgo.Option) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(stackData common2.IStackCreateData, pipeData common2.IPipeCreateData, obs rxgo.Observable) (string, rxgo.Observable, error) {
					if pipeData != nil {
						return "", nil, goerrors.InvalidParam
					}
					if sd, ok := stackData.(*StackData); ok {
						NextInBoundChannel := make(chan rxgo.Item)

						createSendData, err := RxHandlers.CreateSendData(NextInBoundChannel, logger)
						if err != nil {
							return "", nil, err
						}
						err = sd.setOnInBoundSendData(createSendData)
						if err != nil {
							return "", nil, err
						}

						createSendError, err := RxHandlers.CreateSendError(NextInBoundChannel, logger)
						if err != nil {
							return "", nil, err
						}
						err = sd.setOnInBoundSendError(createSendError)
						if err != nil {
							return "", nil, err
						}

						createComplete, err := RxHandlers.CreateComplete(NextInBoundChannel, logger)
						if err != nil {
							return "", nil, err
						}
						err = sd.setOnInBoundComplete(createComplete)
						if err != nil {
							return "", nil, err
						}

						inboundStackHandler, err := NewInboundStackHandler(sd)
						if err != nil {
							return "", nil, err
						}

						nextHandler, err := RxHandlers.NewRxNextHandler(
							common.StackName,
							ConnectionCancelFunc,
							inboundStackHandler,
							sd.onInBoundSendData,
							sd.onInBoundSendError,
							sd.onInBoundComplete,
							logger)
						if err != nil {
							return "", nil, err
						}

						obs.ForEach(
							nextHandler.OnSendData,
							nextHandler.OnError,
							nextHandler.OnComplete,
							opts...)
						nextObs := rxgo.FromChannel(NextInBoundChannel, opts...)
						return common.StackName, nextObs, nil
					}
					return "", nil, WrongStackDataError(connectionType, stackData)
				},
				nil),
			nil
	}
}
