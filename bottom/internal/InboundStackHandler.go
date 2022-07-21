package internal

import (
	"context"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type inboundStackHandler struct {
	errorState error
	logger     *zap.Logger
}

func (self *inboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *inboundStackHandler) ReadMessage(_ interface{}) (interface{}, bool, error) {
	if self.errorState != nil {
		return nil, false, self.errorState
	}
	return nil, false, nil
}

func NewInboundStackHandler(logger *zap.Logger) (RxHandlers.IRxMapStackHandler, error) {
	var err error = nil
	if logger == nil {
		err = multierr.Append(err, goerrors.InvalidParam)
	}
	if err != nil {
		return nil, err
	}
	return &inboundStackHandler{
		errorState: nil,
		logger:     logger,
	}, nil
}

func (self *inboundStackHandler) MapReadWriterSize(
	ctx context.Context,
	rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	return rws, nil
}
