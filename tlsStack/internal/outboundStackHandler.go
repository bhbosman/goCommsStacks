package internal

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"io"
	"time"
)

type outboundStackHandler struct {
	errorState error
	stackData  *Data
}

func (self *outboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *outboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.ConnWrapper.BytesWritten
}

func (self *outboundStackHandler) ReadMessage(_ interface{}) (interface{}, bool, error) {
	return nil, false, nil
}

func (self *outboundStackHandler) Close() error {
	return self.stackData.Close()
}

func (self *outboundStackHandler) SendRws(rws goprotoextra.ReadWriterSize) error {
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

func (self *outboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *outboundStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error) error {
	if self.errorState != nil {
		return self.errorState
	}
	return self.SendRws(rws)
}

func (self *outboundStackHandler) OnComplete() {
	if self.errorState != nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func NewOutboundStackHandler(stackData *Data) (RxHandlers.IRxNextStackHandler, error) {
	return &outboundStackHandler{
		stackData: stackData,
	}, nil
}
