package pingPong

import (
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InboundStackHandler struct {
	StackData *Data
}

func (self *InboundStackHandler) ReadMessage(i interface{}) error {
	return nil
}

func (self *InboundStackHandler) MapReadWriterSize(ctx context.Context, unk interface{}) (interface{}, error) {
	switch rws := unk.(type) {
	case goprotoextra.IReadWriterSize:
		return rws, nil
	default:
		return rws, nil
	}
}

func (self *InboundStackHandler) ErrorState() error {
	return nil
}

func (self *InboundStackHandler) FlatMapHandler(_ context.Context, item interface{}) (RxHandlers.FlatMapHandlerResult, error) {
	switch v := item.(type) {
	case *pingpong.PingWrapper:
		pong := &pingpong.Pong{
			RequestId:         v.Data.RequestId,
			RequestTimeStamp:  v.Data.RequestTimeStamp,
			ResponseTimeStamp: timestamppb.Now(),
		}
		self.StackData.SendPong(pong)
		return RxHandlers.FlatMapHandlerResult{}, nil
	case *pingpong.PongWrapper:
		self.StackData.PongReceived(v.Data)
		return RxHandlers.FlatMapHandlerResult{}, nil
	case *pingpong.Ping:
		return RxHandlers.FlatMapHandlerResult{}, nil
	case *pingpong.Pong:
		return RxHandlers.FlatMapHandlerResult{}, nil
	default:
		return RxHandlers.FlatMapHandlerResult{
			UseDefaultPath: true,
		}, nil
	}
}
