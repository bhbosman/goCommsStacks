package internal

import (
	"context"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type InboundStackHandler struct {
	errorState error
	logger     *zap.Logger
}

func (self *InboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *InboundStackHandler) ReadMessage(_ interface{}) error {
	if self.errorState != nil {
		return self.errorState
	}
	return nil
}

func NewInboundStackHandler(logger *zap.Logger) (*InboundStackHandler, error) {
	var err error = nil
	if logger == nil {
		err = multierr.Append(err, goerrors.InvalidParam)
	}
	if err != nil {
		return nil, err
	}
	return &InboundStackHandler{
		errorState: nil,
		logger:     logger,
	}, nil
}

func (self *InboundStackHandler) MapReadWriterSize(
	ctx context.Context,
	size goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	return size, nil
}
