package internal

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler struct {
	errorState error
	stackData  *StackData
}

func (self *OutboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *OutboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.GetBytesSend()
}

func (self *OutboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func (self *OutboundStackHandler) Close() error {
	return nil
}

func (self *OutboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *OutboundStackHandler) NextReadWriterSize(
	size goprotoextra.ReadWriterSize,
	f func(rws goprotoextra.ReadWriterSize) error,
	_ func(size int) error) error {

	if self.errorState != nil {
		return self.errorState
	}
	_ = f(size)
	return nil
}

func (self *OutboundStackHandler) OnComplete() {
	self.errorState = RxHandlers.RxHandlerComplete
}

func NewOutboundStackHandler(stackData *StackData) (*OutboundStackHandler, error) {
	if stackData == nil {
		return nil, goerrors.InvalidParam
	}
	return &OutboundStackHandler{
		errorState: nil,
		stackData:  stackData,
	}, nil
}
