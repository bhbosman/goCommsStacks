package messageNumber

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"io"
)

type InboundStackHandler struct {
	errorState error
	number     uint64
}

func (self *InboundStackHandler) PublishCounters(counters *model.PublishRxHandlerCounters) {
}

func (self *InboundStackHandler) EmptyQueue() {
}

func (self *InboundStackHandler) ClearCounters() {
}

func (self *InboundStackHandler) FlatMapHandler(_ context.Context, _ interface{}) (RxHandlers.FlatMapHandlerResult, error) {
	return RxHandlers.NewFlatMapHandlerResult(true, nil, 0, 0, 0, 0), nil
}

func (self *InboundStackHandler) ErrorState() error {
	return self.errorState
}

func NewInboundStackHandler() (RxHandlers.IRxMapStackHandler, error) {
	return &InboundStackHandler{
		errorState: nil,
		number:     0,
	}, nil
}

func (self *InboundStackHandler) MapReadWriterSize(ctx context.Context, unk interface{}) (interface{}, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}

	switch v := unk.(type) {
	case goprotoextra.IReadWriterSize:
		if self.errorState != nil {
			return nil, self.errorState
		}
		buffer := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
		n, err := v.Read(buffer[:])
		if err != nil {
			self.errorState = err
			return nil, self.errorState
		}
		if n != 8 {
			self.errorState = io.ErrUnexpectedEOF
			return nil, self.errorState
		}
		newNumber := binary.LittleEndian.Uint64(buffer[:])
		self.number++
		if newNumber != self.number {
			self.errorState = goerrors.InvalidSequenceNumber
			return nil, self.errorState
		}
		return v, nil
	default:
		return v, nil
	}

}
