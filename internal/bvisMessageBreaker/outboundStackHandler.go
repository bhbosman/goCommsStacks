package bvisMessageBreaker

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler struct {
	marker         [4]byte
	markerAsUInt32 uint32
	errorState     error
}

func (self *OutboundStackHandler) Close() error {
	return nil
}

func (self *OutboundStackHandler) PublishCounters(counters *model.PublishRxHandlerCounters) {
}

func (self *OutboundStackHandler) EmptyQueue() {
}

func (self *OutboundStackHandler) ClearCounters() {
}

func (self *OutboundStackHandler) FlatMapHandler(ctx context.Context, item interface{}) (RxHandlers.FlatMapHandlerResult, error) {
	return RxHandlers.NewFlatMapHandlerResult(true, nil, 0, 0, 0, 0), nil
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func NewOutboundStackHandler(marker [4]byte) (RxHandlers.IRxMapStackHandler, error) {
	return &OutboundStackHandler{
		marker:         marker,
		markerAsUInt32: binary.LittleEndian.Uint32(marker[:]),
	}, nil
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
		block := [8]byte{}
		binary.LittleEndian.PutUint32(block[0:4], self.markerAsUInt32)
		binary.LittleEndian.PutUint32(block[4:8], uint32(v.Size()))
		result := gomessageblock.NewReaderWriterBlock(block[:])
		err := result.SetNext(v)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return v, nil
	}

}
