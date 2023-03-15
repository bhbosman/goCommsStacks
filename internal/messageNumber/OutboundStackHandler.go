package messageNumber

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler struct {
	errorState error
	number     uint64
}

func (self *OutboundStackHandler) PublishCounters(*model.PublishRxHandlerCounters) {
}

func (self *OutboundStackHandler) EmptyQueue() {
}

func (self *OutboundStackHandler) ClearCounters() {
}

func (self *OutboundStackHandler) FlatMapHandler(context.Context, interface{}) (RxHandlers.FlatMapHandlerResult, error) {
	return RxHandlers.NewFlatMapHandlerResult(true, nil, 0, 0, 0, 0), nil
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func NewOutboundStackHandler() (RxHandlers.IRxMapStackHandler, error) {
	return &OutboundStackHandler{}, nil
}

func (self *OutboundStackHandler) MapReadWriterSize(ctx context.Context, unk interface{}) (interface{}, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	switch v := unk.(type) {
	case goprotoextra.IReadWriterSize:
		self.number++
		buffer := [8]byte{}
		binary.LittleEndian.PutUint64(buffer[:], self.number)
		rwWithSeqNumber := gomessageblock.NewReaderWriterBlock(buffer[:])
		err := rwWithSeqNumber.SetNext(v)
		if err != nil {
			self.errorState = err
			return nil, err
		}
		return rwWithSeqNumber, nil
	default:
		return v, nil
	}
}
