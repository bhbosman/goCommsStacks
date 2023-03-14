package bvisMessageBreaker

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
)

// Deprecated: Use inboundRxMapStackHandler. Trying to work out all map usage
type inboundRxNextStackHandler struct {
	inboundStackHandler
}

func (self *inboundRxNextStackHandler) PublishCounters(counters *model.PublishRxHandlerCounters) {
}

func (self *inboundRxNextStackHandler) EmptyQueue() {
}

func (self *inboundRxNextStackHandler) ClearCounters() {
}

func (self *inboundRxNextStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *inboundRxNextStackHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *inboundRxNextStackHandler) Close() error {
	return nil
}

func (self *inboundRxNextStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *inboundRxNextStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	onNext func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error,
) error {
	if self.errorState != nil {
		return self.errorState
	}
	err := self.rw.SetNext(rws)
	if err != nil {
		return err
	}
	return self.inboundState(
		func(dataBlock *gomessageblock.ReaderWriter) error {
			return onNext(dataBlock)
		},
	)
}

func (self *inboundRxNextStackHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

// Deprecated: Use NewRxMapStackHandler
func NewInboundStackHandler(
	marker [4]byte,
) RxHandlers.IRxNextStackHandler {
	result := &inboundRxNextStackHandler{
		inboundStackHandler: inboundStackHandler{
			rw:           gomessageblock.NewReaderWriter(),
			errorState:   nil,
			marker:       marker,
			currentState: nil,
		},
	}
	result.currentState = result.onReadHeader()
	return result
}
