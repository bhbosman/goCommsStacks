package internal

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"io"
	"time"
)

type OutboundStackHandler struct {
	errorState error
	stackData  *StackData
}

func (self *OutboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *OutboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.ConnWrapper.BytesWritten
}

func (self *OutboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func (self *OutboundStackHandler) Close() error {
	return self.stackData.Close()
}

func (self *OutboundStackHandler) SendRws(rws goprotoextra.ReadWriterSize) error {
	if self.errorState != nil {
		return self.errorState
	}
	if self.stackData.UpgradedConnection == nil {
		time.Sleep(time.Second)
	}
	_, err := io.Copy(self.stackData.UpgradedConnection, rws)
	if err != nil {
		self.errorState = err
	}
	return self.errorState
}

func (self *OutboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *OutboundStackHandler) NextReadWriterSize(
	size goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(size int) error) error {
	if self.errorState != nil {
		return self.errorState
	}
	return self.SendRws(size)
}

func (self *OutboundStackHandler) OnComplete() {
	if self.errorState != nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func NewOutboundStackHandler(stackData *StackData) (*OutboundStackHandler, error) {
	return &OutboundStackHandler{
		stackData: stackData,
	}, nil
}
