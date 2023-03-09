package protoBuf

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsStacks/internal/protoBuf/proto"
	model2 "github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"reflect"
	"strconv"
)

type InboundStackHandler struct {
	errorState error
	logger     *zap.Logger
	counterMap map[reflect.Type]int
}

func (self *InboundStackHandler) FlatMapHandler(ctx context.Context, item interface{}) (RxHandlers.FlatMapHandlerResult, error) {
	switch v := item.(type) {
	case goprotoextra.ReadWriterSize:
		bytesIn := v.Size()
		tc, err := v.ReadTypeCode()
		if err != nil {
			return RxHandlers.FlatMapHandlerResult{}, err
		}
		_, ok := stream.Find(tc)
		if ok {
			var msgWrapper goprotoextra.IMessageWrapper
			msgWrapper, err = stream.UnMarshal(v)
			if err != nil {
				return RxHandlers.FlatMapHandlerResult{}, err
			}
			self.addCounter(reflect.TypeOf(msgWrapper))
			switch outData := msgWrapper.Message().(type) {
			case *proto.ProtoBufStackMultiMessage:
				rwsTemp := gomessageblock.NewReaderWriter()
				outMessages := make([]interface{}, 0, len(outData.Messages))
				for _, message := range outData.Messages {
					_, _ = rwsTemp.Write(message.MessageData)
					marshal, err := stream.UnMarshal(rwsTemp)
					if err != nil {
						return RxHandlers.FlatMapHandlerResult{}, err
					}
					self.addCounter(reflect.TypeOf(marshal))
					outMessages = append(outMessages, marshal)
				}
				return RxHandlers.NewFlatMapHandlerResult(
					false,
					outMessages,
					len(outMessages),
					0, bytesIn, 0), nil
			default:
				return RxHandlers.FlatMapHandlerResult{
					UseDefaultPath: false,
					Items:          []interface{}{msgWrapper},
					RwsCount:       0,
					OtherCount:     1,
					BytesIn:        bytesIn,
					BytesOut:       0,
				}, err

			}
		} else {
			self.addCounter(reflect.TypeOf(item))
			return RxHandlers.NewFlatMapHandlerResult(true, nil, 0, 0, 0, 0), nil
		}
	default:
		return RxHandlers.NewFlatMapHandlerResult(true, nil, 0, 0, 0, 0), nil
	}
}

func (self *InboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *InboundStackHandler) ReadMessage(i interface{}) error {
	if self.errorState != nil {
		return self.errorState
	}
	switch v := i.(type) {
	case *model2.PublishRxHandlerCounters:
		for r, i := range self.counterMap {
			v.AddMapData(fmt.Sprintf("ProtoBuf Inbound %v", r.String()), strconv.Itoa(i))
		}
		return nil
	}

	return nil
}

func (self *InboundStackHandler) MapReadWriterSize(
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
	switch rws := unk.(type) {
	case goprotoextra.IReadWriterSize:
		return rws, nil
	default:
		return rws, nil
	}
}

func (self *InboundStackHandler) addCounter(of reflect.Type) {
	counter, ok := self.counterMap[of]
	if ok {
		self.counterMap[of] = counter + 1
	} else {
		self.counterMap[of] = 1
	}
}

func NewInboundStackHandler(logger *zap.Logger) (RxHandlers.IRxMapStackHandler, error) {
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
		counterMap: make(map[reflect.Type]int),
	}, nil
}
