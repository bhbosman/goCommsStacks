package internal

import (
	"context"
	"github.com/bhbosman/goprotoextra"
)

type InboundStackHandler struct {
	errorState error
}

func (self *InboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *InboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func NewInboundStackHandler() *InboundStackHandler {
	return &InboundStackHandler{}
}

func (self *InboundStackHandler) MapReadWriterSize(ctx context.Context, size goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	return size, nil
}
