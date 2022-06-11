package internal

import (
	"github.com/bhbosman/goCommsStacks/pingPong/common"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	common3 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Outbound(ConnectionCancelFunc model.ConnectionCancelFunc, logger *zap.Logger, opts ...rxgo.Option) func() (common3.IStackBoundDefinition, error) {
	return func() (common3.IStackBoundDefinition, error) {
		return common3.NewBoundDefinition(
				func(stackData common3.IStackCreateData, pipeData common3.IPipeCreateData, obs rxgo.Observable) (string, rxgo.Observable, error) {
					if sd, ok := stackData.(*StackData); ok {
						outboundChannel := make(chan rxgo.Item)
						var err error
						onRxSendData, err := RxHandlers.CreateSendData(outboundChannel, logger)
						if err != nil {
							return "", nil, err
						}
						err = sd.SetOnRxSendData(onRxSendData)
						if err != nil {
							return "", nil, err
						}

						onRxSendError, err := RxHandlers.CreateSendError(outboundChannel, logger)
						if err != nil {
							return "", nil, err
						}
						err = sd.setOnRxSendError(onRxSendError)
						if err != nil {
							return "", nil, err
						}

						onRxComplete, err := RxHandlers.CreateComplete(outboundChannel, logger)
						if err != nil {
							return "", nil, err
						}
						err = sd.setOnRxComplete(onRxComplete)
						if err != nil {
							return "", nil, err
						}

						outboundStackHandler, err := NewOutboundStackHandler(sd)
						if err != nil {
							return "", nil, err
						}

						rxNextHandler, err := RxHandlers.NewRxNextHandler(
							common.StackName,
							ConnectionCancelFunc,
							outboundStackHandler,
							sd.onRxSendData,
							sd.onRxSendError,
							sd.onRxComplete,
							logger)
						if err != nil {
							return "", nil, err
						}

						_ = obs.ForEach(
							rxNextHandler.OnSendData,
							rxNextHandler.OnError,
							rxNextHandler.OnComplete, opts...)
						obs := rxgo.FromChannel(outboundChannel, opts...)
						return common.StackName, obs, nil
					}
					return common.StackName, nil, goerrors.InvalidType
				},
				nil),
			nil
	}
}
