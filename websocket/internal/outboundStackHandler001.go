package internal

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler001 struct {
	errorState error
	StackData  *data
}

func (self *OutboundStackHandler001) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *OutboundStackHandler001) SendError(err error) {
	if self.StackData.multiOutBoundHandler != nil {
		self.StackData.multiOutBoundHandler.OnError(err)
	}
}

func (self *OutboundStackHandler001) GetAdditionalBytesSend() int {
	return 0
	//return self.data.UpgradedConnection.BytesWrite
}

func (self *OutboundStackHandler001) ReadMessage(_ interface{}) (interface{}, bool, error) {
	return nil, false, nil
}

func (self *OutboundStackHandler001) Close() error {
	return self.StackData.Close()
}

func (self *OutboundStackHandler001) SendData(data interface{}) {
	if self.errorState != nil {
		return
	}
	if self.StackData.multiOutBoundHandler != nil {
		self.StackData.multiOutBoundHandler.OnSendData(data)
	}
}

func NewOutboundStackHandler001(stackData *data) RxHandlers.IRxNextStackHandler {
	return &OutboundStackHandler001{
		errorState: nil,
		StackData:  stackData}
}

func (self *OutboundStackHandler001) OnError(err error) {
	self.errorState = err
}

func (self *OutboundStackHandler001) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	sizeUpDate func(size int) error) error {
	if self.errorState != nil {
		return self.errorState
	}
	sizeUpDate(rws.Size())
	self.SendData(rws)

	return self.errorState
}

func (self *OutboundStackHandler001) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *OutboundStackHandler001) Complete() {

}
