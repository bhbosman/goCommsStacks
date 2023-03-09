package internal

import (
	"context"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"io"
)

type inboundStackHandler struct {
	errorState error
	stackData  *Data
}

func (self *inboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *inboundStackHandler) SendError(err error) {
	if self.stackData.InboundHandler != nil {
		self.stackData.InboundHandler.OnError(err)
	}
}

func (self *inboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.UpgradedConnection.BytesRead
}

func (self *inboundStackHandler) ReadMessage(_ interface{}) (interface{}, bool, error) {
	return nil, false, nil
}

func (self *inboundStackHandler) Close() error {
	return self.stackData.Close()
}

func (self *inboundStackHandler) SendData(data interface{}) {
	if self.stackData.InboundHandler != nil {
		self.stackData.InboundHandler.OnSendData(data)
	}
}

func (self *inboundStackHandler) SendRws(
	rws goprotoextra.ReadWriterSize,
) {
	_, err := io.Copy(self.stackData.PipeWriteClose, rws)
	if err != nil {
		self.errorState = err
		return
	}
}

func (self *inboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *inboundStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	update func(size int) error,
) error {

	if self.errorState != nil {
		return self.errorState
	}
	n := rws.Size()
	_ = update(n)
	self.SendRws(rws)
	return self.errorState
}

func (self *inboundStackHandler) OnComplete() {
	if self.errorState != nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *inboundStackHandler) Complete() {

}

func NewInboundStackHandler(
	stackData *Data,
	parentContext context.Context,
	parentCancelFunc context.CancelFunc,
) (RxHandlers.IRxNextStackHandler, error) {
	//cancel, cancelFunc := context.WithCancel(parentContext)
	return &inboundStackHandler{
		stackData: stackData,
	}, nil
}
