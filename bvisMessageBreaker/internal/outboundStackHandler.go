package internal

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler struct {
	marker         [4]byte
	markerAsUInt32 uint32
	errorState     error
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *OutboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func NewOutboundStackHandler(marker [4]byte) (*OutboundStackHandler, error) {
	return &OutboundStackHandler{
		marker:         marker,
		markerAsUInt32: binary.LittleEndian.Uint32(marker[:]),
	}, nil
}

func (self *OutboundStackHandler) MapReadWriterSize(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	block := [8]byte{}
	binary.LittleEndian.PutUint32(block[0:4], self.markerAsUInt32)
	binary.LittleEndian.PutUint32(block[4:8], uint32(rws.Size()))
	result := gomessageblock.NewReaderWriterBlock(block[:])
	err := result.SetNext(rws)
	if err != nil {
		return nil, err
	}
	return result, nil
}
