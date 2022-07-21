package internal

import (
	"context"
	"github.com/bhbosman/gocommon/stream"
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

func (self *InboundStackHandler) ReadMessage(i interface{}) (interface{}, bool, error) {
	if self.errorState != nil {
		return nil, false, self.errorState
	}
	if unk, ok := i.(goprotoextra.ReadWriterSize); ok {
		var err error
		var tc uint32
		tc, err = unk.ReadTypeCode()
		if err != nil {
			return nil, false, err
		}
		_, ok = stream.Find(tc)
		if ok {
			var msgWrapper goprotoextra.IMessageWrapper
			msgWrapper, err = stream.UnMarshal(unk)
			if err != nil {
				return nil, false, err
			}
			return msgWrapper, true, nil
		}
	}
	return nil, false, nil
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
