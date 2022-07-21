package internal

import (
	"context"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type OutboundStackHandler struct {
	errorState error
	logger     *zap.Logger
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *OutboundStackHandler) ReadMessage(msg interface{}) (interface{}, bool, error) {
	if unk, ok := msg.(goprotoextra.IMessageWrapper); ok {
		msg = unk.Message()
	}
	if unk, ok := msg.(proto.Message); ok {
		rws, err := stream.Marshall(unk)
		if err != nil {
			return nil, false, err
		}
		return rws, true, nil
	}
	return nil, false, nil
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
