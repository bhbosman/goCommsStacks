package internal

import (
	"context"
	"github.com/bhbosman/goprotoextra"
)

type OutboundStackHandler struct {
	errorState error
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *OutboundStackHandler) ReadMessage(_ interface{}) (interface{}, bool, error) {
	return nil, false, nil
}

func NewOutboundStackHandler() *OutboundStackHandler {
	return &OutboundStackHandler{}
}

func (self *OutboundStackHandler) MapReadWriterSize(
	ctx context.Context,
	rws goprotoextra.ReadWriterSize,
) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	return rws, nil
}
