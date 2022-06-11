package internal

import (
	"context"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type OutboundStackHandler struct {
	errorState error
	logger     *zap.Logger
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *OutboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func NewOutboundStackHandler(logger *zap.Logger) (*OutboundStackHandler, error) {
	var errList error = nil
	if logger == nil {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}
	if errList != nil {
		return nil, errList
	}
	return &OutboundStackHandler{
		logger: logger,
	}, nil
}

func (self *OutboundStackHandler) MapReadWriterSize(ctx context.Context, size goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	return size, nil
}
