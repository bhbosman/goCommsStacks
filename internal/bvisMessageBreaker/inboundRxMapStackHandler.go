package bvisMessageBreaker

import (
	"fmt"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gomessageblock"
	"golang.org/x/net/context"
)

type inboundRxMapStackHandler struct {
	inboundStackHandler
}

func (self *inboundRxMapStackHandler) PublishCounters(counters *model.PublishRxHandlerCounters) {
}

func (self *inboundRxMapStackHandler) EmptyQueue() {
}

func (self *inboundRxMapStackHandler) ClearCounters() {
}

func (self *inboundRxMapStackHandler) MapReadWriterSize(ctx context.Context, unk interface{}) (interface{}, error) {
	return nil, fmt.Errorf("not implemented. use FlatMapHandler")
}

func (self *inboundRxMapStackHandler) ErrorState() error {
	return self.errorState
}

func (self *inboundRxMapStackHandler) FlatMapHandler(ctx context.Context, item interface{}) (RxHandlers.FlatMapHandlerResult, error) {
	switch v := item.(type) {
	case *gomessageblock.ReaderWriter:
		err := self.rw.SetNext(v)
		if err != nil {
			return RxHandlers.FlatMapHandlerResult{}, err
		}
		var result []interface{}
		count := 0
		bytesIn := v.Size()
		bytesOut := 0
		err = self.inboundState(
			func(dataBlock *gomessageblock.ReaderWriter) error {
				result = append(result, dataBlock)
				count++
				bytesOut += dataBlock.Size()
				return nil
			},
		)
		return RxHandlers.FlatMapHandlerResult{
			UseDefaultPath: false,
			Items:          result,
			RwsCount:       count,
			OtherCount:     0,
			BytesIn:        bytesIn,
			BytesOut:       bytesOut,
		}, nil

	default:
		return RxHandlers.FlatMapHandlerResult{
			UseDefaultPath: true,
			Items:          nil,
			RwsCount:       0,
			OtherCount:     1,
		}, nil
	}
}

func NewRxMapStackHandler(
	marker [4]byte,
) (RxHandlers.IRxMapStackHandler, error) {
	result := &inboundRxMapStackHandler{
		inboundStackHandler: inboundStackHandler{
			rw:           gomessageblock.NewReaderWriter(),
			errorState:   nil,
			marker:       marker,
			currentState: nil,
		},
	}
	result.currentState = result.onReadHeader()
	return result, nil
}
