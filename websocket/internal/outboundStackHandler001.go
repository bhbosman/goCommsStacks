package internal

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler001 struct {
	errorState error
	StackData  *StackData
}

func (self *OutboundStackHandler001) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *OutboundStackHandler001) SendError(err error) {
	if self.StackData.onMultiOutBoundSendError != nil {
		self.StackData.onMultiOutBoundSendError(err)
	}
}

func (self *OutboundStackHandler001) GetAdditionalBytesSend() int {
	return 0
	//return self.StackData.UpgradedConnection.BytesWrite
}

func (self *OutboundStackHandler001) ReadMessage(_ interface{}) error {
	return nil
}

func (self *OutboundStackHandler001) Close() error {
	return self.StackData.Close()
}

func (self *OutboundStackHandler001) SendData(data interface{}) {
	if self.errorState != nil {
		return
	}
	if self.StackData.onMultiOutBoundSendData != nil {
		self.StackData.onMultiOutBoundSendData(data)
	}
}

func NewOutboundStackHandler001(stackData *StackData) *OutboundStackHandler001 {
	return &OutboundStackHandler001{
		errorState: nil,
		StackData:  stackData}
}

func (self *OutboundStackHandler001) OnError(err error) {
	self.errorState = err
}

func (self *OutboundStackHandler001) NextReadWriterSize(
	size goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	sizeUpDate func(size int) error) error {
	if self.errorState != nil {
		return self.errorState
	}
	sizeUpDate(size.Size())
	self.SendData(size)

	return self.errorState
}

func (self *OutboundStackHandler001) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *OutboundStackHandler001) Complete() {

}
