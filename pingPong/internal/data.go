package internal

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/goprotoextra"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type data struct {
	bytesSend         int
	requestId         int64
	ctx               context.Context
	goFunctionCounter GoFunctionCounter.IService
	outHandler        goCommsDefinitions.IRxNextHandler
}

func (self *data) SendError(err error) {
	if self.outHandler != nil {
		self.outHandler.OnError(err)
	}
}

func (self *data) SendPong(pong *pingpong.Pong) {
	if self.outHandler != nil {
		self.outHandler.OnSendData(pong)
	}
}

func (self *data) GetBytesSend() int {
	return self.bytesSend
}

func (self *data) SendRws(rws goprotoextra.ReadWriterSize) {
	if self.outHandler != nil {
		self.outHandler.OnSendData(rws)
	}
}

func (self *data) PongReceived(_ *pingpong.Pong) {
}

func NewStackData(
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
) *data {
	return &data{
		requestId:         0,
		ctx:               ctx,
		goFunctionCounter: goFunctionCounter,
	}
}

func (self *data) Destroy() error {
	return nil
}

func (self *data) Ping() {
	sendTime := time.Now()
	msg := &pingpong.Ping{
		RequestId: 0,
		RequestTimeStamp: &timestamppb.Timestamp{
			Seconds: sendTime.Unix(),
			Nanos:   int32(sendTime.Nanosecond()),
		},
	}
	rws, err := stream.Marshall(msg)
	if err != nil {
		return
	}
	self.bytesSend += rws.Size()
	if self.outHandler != nil {
		self.outHandler.OnSendData(rws)
	}
}

func (self *data) Start(_ context.Context) error {
	return self.goFunctionCounter.GoRun("PingPong.data.Start",
		func() {
			//
			ticker := time.NewTicker(time.Second * 1)
			defer func(ticker *time.Ticker) {
				ticker.Stop()
			}(ticker)

			//
		loop:
			for true {
				select {
				case <-self.ctx.Done():
					break loop
				case _, ok := <-ticker.C:
					if ok {
						self.Ping()
					}
					continue loop
				}
			}
		},
	)
}

func (self *data) Stop() error {
	return nil
}
