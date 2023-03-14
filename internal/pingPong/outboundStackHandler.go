package pingPong

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
)

type outboundStackHandler struct {
	errorState error
	stackData  *Data
}

func (self *outboundStackHandler) PublishCounters(counters *model.PublishRxHandlerCounters) {
}

func (self *outboundStackHandler) EmptyQueue() {
}

func (self *outboundStackHandler) ClearCounters() {
}

func (self *outboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *outboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.GetBytesSend()
}

func (self *outboundStackHandler) Close() error {
	return nil
}

func (self *outboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *outboundStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	f func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error,
) error {

	if self.errorState != nil {
		return self.errorState
	}
	return f(rws)

}

func (self *outboundStackHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func NewOutboundStackHandler(stackData *Data) (RxHandlers.IRxNextStackHandler, error) {
	if stackData == nil {
		return nil, goerrors.InvalidParam
	}
	return &outboundStackHandler{
		errorState: nil,
		stackData:  stackData,
	}, nil
}
