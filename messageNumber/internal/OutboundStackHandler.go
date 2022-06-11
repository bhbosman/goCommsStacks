package internal

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler struct {
	errorState error
	number     uint64
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *OutboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func NewOutboundStackHandler() (*OutboundStackHandler, error) {
	return &OutboundStackHandler{}, nil
}

func (self *OutboundStackHandler) MapReadWriterSize(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	self.number++
	buffer := [8]byte{}
	binary.LittleEndian.PutUint64(buffer[:], self.number)
	rwWithSeqNumber := gomessageblock.NewReaderWriterBlock(buffer[:])
	err := rwWithSeqNumber.SetNext(rws)
	if err != nil {
		self.errorState = err
		return nil, err
	}
	return rwWithSeqNumber, nil
}
