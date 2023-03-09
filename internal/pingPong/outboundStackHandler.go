package pingPong

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
)

type outboundStackHandler struct {
	errorState error
	stackData  *Data
}

func (self *outboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *outboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.GetBytesSend()
}

func (self *outboundStackHandler) ReadMessage(_ interface{}) (interface{}, bool, error) {
	return nil, false, nil
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
	_ = f(rws)
	return nil
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
