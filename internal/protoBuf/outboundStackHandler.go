package protoBuf

import (
	"context"
	proto2 "github.com/bhbosman/goCommsStacks/internal/protoBuf/proto"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type OutboundStackHandler struct {
	stackHandler
}

func (self *OutboundStackHandler) FlatMapHandler(_ context.Context, _ interface{}) (RxHandlers.FlatMapHandlerResult, error) {
	return RxHandlers.NewFlatMapHandlerResult(true, nil, 0, 0, 0, 0), nil
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *OutboundStackHandler) MapReadWriterSize(
	ctx context.Context,
	unk interface{},
) (interface{}, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	localMarshall := func(m proto.Message) (interface{}, error) {
		self.addCounter(reflect.TypeOf(m))

		rws, err := stream.Marshall(m)
		if err != nil {
			return nil, err
		}
		return rws, nil
	}
	switch v := unk.(type) {
	case []proto.Message:
		outData := proto2.ProtoBufStackMultiMessage{
			Messages: make([]*proto2.ProtoBufStackMessageInstance, 0, len(v)),
		}
		for _, data := range v {
			self.addCounter(reflect.TypeOf(data))
			dataAsBytes, err := stream.Marshall(data)
			if err != nil {
				return nil, err
			}
			flatten, err := dataAsBytes.Flatten()
			if err != nil {
				return nil, err
			}
			outData.Messages = append(outData.Messages, &proto2.ProtoBufStackMessageInstance{MessageData: flatten})
		}
		return localMarshall(&outData)
	case goprotoextra.IMessageWrapper:
		if protoMessage, ok := v.(proto.Message); ok {
			return localMarshall(protoMessage)
		}
		return v, nil
	case proto.Message:
		return localMarshall(v)
	default:
		return v, nil
	}
}

func NewOutboundStackHandler(logger *zap.Logger) (RxHandlers.IRxMapStackHandler, error) {
	var errList error = nil
	if logger == nil {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}
	if errList != nil {
		return nil, errList
	}
	return &OutboundStackHandler{
		stackHandler: stackHandler{
			errorState: nil,
			logger:     logger,
			counterMap: make(map[reflect.Type]int),
			prefix:     "Outbound",
		},
	}, nil
}
