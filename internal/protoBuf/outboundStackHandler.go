package protoBuf

import (
	"context"
	"fmt"
	proto2 "github.com/bhbosman/goCommsStacks/internal/protoBuf/proto"
	model2 "github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strconv"
)

type OutboundStackHandler struct {
	errorState error
	logger     *zap.Logger
	counterMap map[reflect.Type]int
}

func (self *OutboundStackHandler) FlatMapHandler(_ context.Context, _ interface{}) (RxHandlers.FlatMapHandlerResult, error) {
	return RxHandlers.NewFlatMapHandlerResult(true, nil, 0, 0, 0, 0), nil
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *OutboundStackHandler) ReadMessage(msg interface{}) error {
	switch v := msg.(type) {
	case *model2.PublishRxHandlerCounters:
		for r, i := range self.counterMap {
			v.AddMapData(fmt.Sprintf("ProtoBuf Outbound %v", r.String()), strconv.Itoa(i))
		}
		return nil
	default:
		return nil
	}
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

func (self *OutboundStackHandler) addCounter(of reflect.Type) {
	counter, ok := self.counterMap[of]
	if ok {
		self.counterMap[of] = counter + 1
	} else {
		self.counterMap[of] = 1
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
		logger:     logger,
		counterMap: make(map[reflect.Type]int),
	}, nil
}
