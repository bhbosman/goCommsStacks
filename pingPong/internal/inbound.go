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

func Inbound(ConnectionCancelFunc model.ConnectionCancelFunc, logger *zap.Logger, opts ...rxgo.Option) func() (common3.IStackBoundDefinition, error) {
	return func() (common3.IStackBoundDefinition, error) {
		return common3.NewBoundDefinition(
				func(stackData common3.IStackCreateData, pipeData common3.IPipeCreateData, obs rxgo.Observable) (string, rxgo.Observable, error) {
					if stackDataActual, ok := stackData.(*StackData); ok {
						inBoundChannel := make(chan rxgo.Item)
						inboundStackHandler := NewInboundStackHandler(stackDataActual)

						createComplete, err := RxHandlers.CreateComplete(inBoundChannel, logger)
						if err != nil {
							return "", nil, err
						}

						createSendError, err := RxHandlers.CreateSendError(inBoundChannel, logger)
						if err != nil {
							return "", nil, err
						}

						createSendData, err := RxHandlers.CreateSendData(inBoundChannel, logger)
						if err != nil {
							return "", nil, err
						}

						rxNextHandler, err := RxHandlers.NewRxNextHandler(
							common.StackName,
							ConnectionCancelFunc,
							inboundStackHandler,
							createSendData,
							createSendError,
							createComplete,
							logger)
						if err != nil {
							return "", nil, err
						}

						_ = obs.ForEach(
							rxNextHandler.OnSendData,
							rxNextHandler.OnError,
							rxNextHandler.OnComplete,
							opts...)

						obs := rxgo.FromChannel(inBoundChannel, opts...)
						return common.StackName, obs, nil
					}
					return common.StackName, nil, goerrors.InvalidType
				},
				nil),
			nil
	}
}
