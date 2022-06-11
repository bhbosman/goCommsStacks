package internal

import (
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"github.com/golang/protobuf/ptypes"
)

type InboundStackHandler struct {
	errorState error
	stackData  *StackData
}

func (self *InboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *InboundStackHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *InboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
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
	_ func(size int) error) error {

	if self.errorState != nil {
		return self.errorState
	}
	var tc uint32
	var err error
	tc, err = incoming.ReadTypeCode()
	if err != nil {
		self.errorState = err
		return self.errorState
	}
	switch tc {
	case pingpong.PingTypeCode:
		msg, err := stream.UnMarshal(incoming, nil, nil, nil, nil)
		if err != nil {
			self.errorState = err
			return self.errorState
		}
		switch v := msg.(type) {
		case *pingpong.PingWrapper:
			pong := &pingpong.Pong{
				RequestId:         v.Data.RequestId,
				RequestTimeStamp:  v.Data.RequestTimeStamp,
				ResponseTimeStamp: ptypes.TimestampNow(),
			}
			marshall, err := stream.Marshall(pong)
			if err != nil {
				self.errorState = err
				return self.errorState
			}
			self.stackData.SendPong(marshall)
		}
		return nil

	case pingpong.PongTypeCode:
		msg := &pingpong.Pong{}
		err := stream.UnMarshalMessage(incoming, msg)
		if err != nil {
			self.errorState = err
			return self.errorState
		}
		self.stackData.PongReceived(msg)
	}

	_ = onNext(incoming)
	return nil
}

func (self *InboundStackHandler) OnComplete() {
	self.errorState = RxHandlers.RxHandlerComplete
}

func NewInboundStackHandler(stackData *StackData) *InboundStackHandler {
	return &InboundStackHandler{
		errorState: nil,
		stackData:  stackData,
	}
}
