package pingPong

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/goprotoextra"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Data struct {
	bytesSend         int
	requestId         int64
	ctx               context.Context
	goFunctionCounter GoFunctionCounter.IService
	OutHandler        goCommsDefinitions.IRxNextHandler
	Ticker            *time.Ticker
}

func (self *Data) SendError(err error) {
	if self.OutHandler != nil {
		self.OutHandler.OnError(err)
	}
}

func (self *Data) SendPong(pong *pingpong.Pong) {
	if self.OutHandler != nil {
		self.OutHandler.OnSendData(pong)
	}
}

func (self *Data) GetBytesSend() int {
	return self.bytesSend
}

func (self *Data) SendRws(rws goprotoextra.ReadWriterSize) {
	if self.OutHandler != nil {
		self.OutHandler.OnSendData(rws)
	}
}

func (self *Data) PongReceived(_ *pingpong.Pong) {
}

func NewStackData(
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
) *Data {
	return &Data{
		requestId:         0,
		ctx:               ctx,
		goFunctionCounter: goFunctionCounter,
		Ticker:            time.NewTicker(time.Hour * 24),
	}
}

func (self *Data) Destroy() error {
	return nil
}

func (self *Data) Start(_ context.Context) error {
	self.Ticker.Reset(time.Minute)
	return nil
}

func (self *Data) Stop() error {
	self.Ticker.Stop()
	return nil
}

func (self *Data) SendPing(time2 time.Time) interface{} {
	sendTime := time.Now()
	return &pingpong.Ping{
		RequestId: 0,
		RequestTimeStamp: &timestamppb.Timestamp{
			Seconds: sendTime.Unix(),
			Nanos:   int32(sendTime.Nanosecond()),
		},
	}
}
