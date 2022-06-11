package internal

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
)

type InboundStackHandler struct {
	errorState error
	number     uint64
}

func (self *InboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *InboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func NewInboundStackHandler() (*InboundStackHandler, error) {
	return &InboundStackHandler{
		errorState: nil,
		number:     0,
	}, nil
}

func (self *InboundStackHandler) MapReadWriterSize(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}

	if self.errorState != nil {
		return nil, self.errorState
	}
	buffer := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
	_, err := rws.Read(buffer[:])
	if err != nil {
		self.errorState = err
		return nil, self.errorState
	}
	newNumber := binary.LittleEndian.Uint64(buffer[:])
	self.number++
	if newNumber != self.number {
		self.errorState = goerrors.InvalidSequenceNumber
		return nil, self.errorState
	}
	return rws, nil
}
