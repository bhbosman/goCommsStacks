package internal

import (
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"github.com/golang/protobuf/ptypes"
)

type InboundStackHandler struct {
	errorState error
	stackData  *data
}

func (self *InboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *InboundStackHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *InboundStackHandler) ReadMessage(i interface{}) (interface{}, bool, error) {
	switch v := i.(type) {
	case *pingpong.PongWrapper:
		self.stackData.PongReceived(v.Data)
		return nil, false, nil
	case *pingpong.PingWrapper:
		pong := &pingpong.Pong{
			RequestId:         v.Data.RequestId,
			RequestTimeStamp:  v.Data.RequestTimeStamp,
			ResponseTimeStamp: ptypes.TimestampNow(),
		}
		self.stackData.SendPong(pong)
		return nil, false, nil
	default:
		return nil, false, nil
	}
}

func (self *InboundStackHandler) Close() error {
	return nil
}

func (self *InboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *InboundStackHandler) NextReadWriterSize(
	incoming goprotoextra.ReadWriterSize,
	onNext func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error,
) error {
	return nil
}

func (self *InboundStackHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func NewInboundStackHandler(stackData *data) RxHandlers.IRxNextStackHandler {
	return &InboundStackHandler{
		errorState: nil,
		stackData:  stackData,
	}
}
