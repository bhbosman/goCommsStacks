package internal

import (
	"fmt"
	"github.com/bhbosman/goCommsStacks/websocket/common"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Outbound(connectionType model.ConnectionType, ConnectionCancelFunc model.ConnectionCancelFunc, logger *zap.Logger, opts ...rxgo.Option) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(stackData common2.IStackCreateData, pipeData common2.IPipeCreateData, obs rxgo.Observable) (string, rxgo.Observable, error) {
					if sd, ok := stackData.(*StackData); ok {
						MultiReceiverChannel := make(chan rxgo.Item)
						var err error
						sd.onMultiOutBoundSendData, err = RxHandlers.CreateSendData(MultiReceiverChannel, logger)
						if err != nil {
							return "", nil, err
						}

						sd.onMultiOutBoundSendError, err = RxHandlers.CreateSendError(MultiReceiverChannel, logger)
						if err != nil {
							return "", nil, err
						}

						sd.onMultiOutBoundComplete, err = RxHandlers.CreateComplete(MultiReceiverChannel, logger)
						if err != nil {
							return "", nil, err
						}

						outboundStackHandler001 := NewOutboundStackHandler001(sd)
						nextHandler001, err := RxHandlers.NewRxNextHandler(
							fmt.Sprintf("%v-01", common.StackName),
							ConnectionCancelFunc,
							outboundStackHandler001,
							sd.onMultiOutBoundSendData,
							sd.onMultiOutBoundSendError,
							sd.onMultiOutBoundComplete,
							logger)
						if err != nil {
							return "", nil, err
						}

						_ = obs.ForEach(nextHandler001.OnSendData, nextHandler001.OnError, nextHandler001.OnComplete, opts...)

						NextOutboundChannelManager := make(chan rxgo.Item)
						sd.onOutBoundSendData, err = RxHandlers.CreateSendData(NextOutboundChannelManager, logger)
						if err != nil {
							return "", nil, err
						}

						sd.onOutBoundSendError, err = RxHandlers.CreateSendError(NextOutboundChannelManager, logger)
						if err != nil {
							return "", nil, err
						}

						sd.onOutBoundComplete, err = RxHandlers.CreateComplete(NextOutboundChannelManager, logger)
						if err != nil {
							return "", nil, err
						}

						connWrapper, err := common2.NewConnWrapper(
							sd.conn,
							sd.ctx,
							sd.pipeRead,
							sd.onOutBoundSendData)
						if err != nil {
							return "", nil, err
						}

						err = sd.setConnWrapper(connWrapper)
						if err != nil {
							return "", nil, err
						}

						outboundStackHandler002, err := NewOutboundStackHandler002(sd)
						if err != nil {
							return "", nil, err
						}

						nextHandler002, err := RxHandlers.NewRxNextHandler(
							fmt.Sprintf("%v-02", common.StackName),
							ConnectionCancelFunc,
							outboundStackHandler002,
							sd.onOutBoundSendData,
							sd.onOutBoundSendError,
							sd.onOutBoundComplete,
							logger)
						if err != nil {
							return "", nil, err
						}

						tempObs := rxgo.FromChannel(MultiReceiverChannel)
						_ = tempObs.ForEach(nextHandler002.OnSendData, nextHandler002.OnError, nextHandler002.OnComplete, opts...)

						return common.StackName, rxgo.FromChannel(NextOutboundChannelManager), nil
					}
					return "", nil, WrongStackDataError(connectionType, stackData)
				},
				nil),
			nil
	}
}
