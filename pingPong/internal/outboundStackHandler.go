package internal

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler struct {
	errorState error
	stackData  *data
}

func (self *OutboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *OutboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.GetBytesSend()
}

func (self *OutboundStackHandler) ReadMessage(_ interface{}) (interface{}, bool, error) {
	return nil, false, nil
}

func (self *OutboundStackHandler) Close() error {
	return nil
}

func (self *OutboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *OutboundStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	f func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error,
) error {

	if self.errorState != nil {
		return self.errorState
	}
	_ = f(rws)
	return nil
}

func (self *OutboundStackHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func NewOutboundStackHandler(stackData *data) (RxHandlers.IRxNextStackHandler, error) {
	if stackData == nil {
		return nil, goerrors.InvalidParam
	}
	return &OutboundStackHandler{
		errorState: nil,
		stackData:  stackData,
	}, nil
}
