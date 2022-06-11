package internal

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"io"
)

type InboundStackHandler struct {
	errorState error
	stackData  *StackData
}

func (self *InboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *InboundStackHandler) SendError(err error) {
	if self.stackData.onInBoundSendError != nil {
		self.stackData.onInBoundSendError(err)
	}
}

func (self *InboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.socketDataReceived
}

func (self *InboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func (self *InboundStackHandler) sendRws(rws goprotoextra.ReadWriterSize) {
	_, err := io.Copy(self.stackData.pipeWriteClose, rws)
	if err != nil {
		self.errorState = err
	}
}

func (self *InboundStackHandler) Close() error {
	return self.stackData.Close()
}

func (self *InboundStackHandler) SendData(data interface{}) {
	if self.stackData.onInBoundSendData != nil {
		self.stackData.onInBoundSendData(data)
	}
}

func (self *InboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *InboundStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(size int) error) error {

	if self.errorState != nil {
		return self.errorState
	}
	self.sendRws(rws)
	return self.errorState
}

func (self *InboundStackHandler) OnComplete() {
	self.errorState = RxHandlers.RxHandlerComplete
}

func (self *InboundStackHandler) Complete() {

}

func NewInboundStackHandler(stackData *StackData) (*InboundStackHandler, error) {
	return &InboundStackHandler{
		errorState: nil,
		stackData:  stackData,
	}, nil
}
