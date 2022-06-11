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

func Outbound(connectionType model.ConnectionType, ConnectionCancelFunc model.ConnectionCancelFunc, logger *zap.Logger, opts ...rxgo.Option) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(sd common2.IStackCreateData, pipeData common2.IPipeCreateData, obs rxgo.Observable) (string, rxgo.Observable, error) {
					if pipeData != nil {
						return "", nil, goerrors.InvalidParam
					}
					if stackData, ok := sd.(*StackData); ok {
						nextOutboundChannel := make(chan rxgo.Item)

						createSendData, err := RxHandlers.CreateSendData(nextOutboundChannel, logger)
						if err != nil {
							return "", nil, err
						}
						err = stackData.setOnOutBoundSendData(createSendData)
						if err != nil {
							return "", nil, err
						}

						createSendError, err := RxHandlers.CreateSendError(nextOutboundChannel, logger)
						if err != nil {
							return "", nil, err
						}
						err = stackData.setOnOutBoundSendError(createSendError)
						if err != nil {
							return "", nil, err
						}

						createComplete, err := RxHandlers.CreateComplete(nextOutboundChannel, logger)
						if err != nil {
							return "", nil, err
						}
						err = stackData.setOnOutBoundComplete(createComplete)
						if err != nil {
							return "", nil, err
						}

						connWrapper, err := common2.NewConnWrapper(
							stackData.conn,
							stackData.ctx,
							stackData.PipeRead,
							stackData.onOutBoundSendData)
						if err != nil {
							return "", nil, err
						}

						err = stackData.setConnWrapper(connWrapper)
						if err != nil {
							return "", nil, err
						}

						outboundStackHandler, err := NewOutboundStackHandler(stackData)
						if err != nil {
							return "", nil, err
						}

						nextHandler, err := RxHandlers.NewRxNextHandler(
							common.StackName,
							ConnectionCancelFunc,
							outboundStackHandler,
							stackData.onOutBoundSendData,
							stackData.onOutBoundSendError,
							stackData.onOutBoundComplete,
							logger)
						if err != nil {
							return "", nil, err
						}

						obs.ForEach(nextHandler.OnSendData, nextHandler.OnError, nextHandler.OnComplete, opts...)
						nextObs := rxgo.FromChannel(nextOutboundChannel, opts...)
						return common.StackName, nextObs, nil
					}
					return "", nil, WrongStackDataError(connectionType, sd)
				},
				nil),
			nil
	}
}
