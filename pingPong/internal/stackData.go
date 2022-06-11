package internal

import (
	"context"
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type StackData struct {
	bytesSend     int
	requestId     int64
	ctx           context.Context
	onRxSendData  func(data interface{})
	onRxSendError func(err error)
	onRxComplete  func()
}

func (self *StackData) SendError(err error) {
	if self.onRxSendError != nil {
		self.onRxSendError(err)
	}
}

func (self *StackData) SendPong(rws goprotoextra.ReadWriterSize) {
	if self.onRxSendData != nil {
		self.bytesSend += rws.Size()
		self.onRxSendData(rws)
	}
}

func (self *StackData) GetBytesSend() int {
	return self.bytesSend
}

func (self *StackData) SendRws(rws goprotoextra.ReadWriterSize) {
	if self.onRxSendData != nil {
		self.onRxSendData(rws)
	}
}

func (self *StackData) PongReceived(_ *pingpong.Pong) {
	//receiveTime := time.Now()
	//sendTime := time.Unix(v.RequestTimeStamp.Seconds, int64(v.RequestTimeStamp.Nanos))
	//deltaTime := receiveTime.Sub(sendTime)
}

//func (self *StackData) OutboundChannel() *ChannelManager.ChannelManager {
//	return self.outboundChannel
//}

func NewStackData(ctx context.Context) *StackData {
	return &StackData{
		requestId: 0,
		ctx:       ctx,
	}
}

func (self *StackData) Destroy() error {
	return nil
}

func (self *StackData) Ping() {
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
	if self.onRxSendData != nil {
		self.onRxSendData(rws)
	}
}

func (self *StackData) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 1)
		defer func(ticker *time.Ticker) {
			ticker.Stop()
		}(ticker)
	loop:
		for true {
			select {
			case <-ctx.Done():
				break loop
			case _, ok := <-ticker.C:
				if ok {
					self.Ping()
				}
				continue loop
			}
		}
	}(self.ctx)

	return nil
}

func (self *StackData) Stop() error {
	return nil
}

func (self *StackData) SetOnRxSendData(data rxgo.NextFunc) error {
	if data == nil {
		return goerrors.InvalidParam
	}
	self.onRxSendData = data
	return nil
}

func (self *StackData) setOnRxSendError(sendError rxgo.ErrFunc) error {
	if sendError == nil {
		return goerrors.InvalidParam
	}
	self.onRxSendError = sendError
	return nil
}

func (self *StackData) setOnRxComplete(complete rxgo.CompletedFunc) error {
	if complete == nil {
		return goerrors.InvalidParam
	}
	self.onRxComplete = complete
	return nil
}
