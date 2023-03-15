package internal

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"io"
)

type inboundStackHandler struct {
	errorState error
	stackData  *Data
}

func (self *inboundStackHandler) EmptyQueue() {
}

func (self *inboundStackHandler) ClearCounters() {
}

func (self *inboundStackHandler) PublishCounters(*model.PublishRxHandlerCounters) {
}

func (self *inboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *inboundStackHandler) GetAdditionalBytesSend() int {
	return self.stackData.UpgradedConnection.BytesRead
}

func (self *inboundStackHandler) Close() error {
	return self.stackData.Close()
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

func NewInboundStackHandler(stackData *Data) (RxHandlers.IRxNextStackHandler, error) {
	return &inboundStackHandler{
		stackData: stackData,
	}, nil
}
