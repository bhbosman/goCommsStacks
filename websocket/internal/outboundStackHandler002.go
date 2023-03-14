package internal

import (
	wsmsg2 "github.com/bhbosman/goCommsStacks/webSocketMessages/wsmsg"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"time"
)

type OutboundStackHandler002 struct {
	errorState error
	stackData  *Data
}

func (self *OutboundStackHandler002) PublishCounters(counters *model.PublishRxHandlerCounters) {
}

func (self *OutboundStackHandler002) EmptyQueue() {
}

func (self *OutboundStackHandler002) ClearCounters() {
}

func (self *OutboundStackHandler002) GetAdditionalBytesIncoming() int {
	return -self.stackData.additionalDataSend
}

func (self *OutboundStackHandler002) GetAdditionalBytesSend() int {
	return self.stackData.connWrapper.BytesWritten
}

func NewOutboundStackHandler002(stackData *Data) (RxHandlers.IRxNextStackHandler, error) {
	return &OutboundStackHandler002{
		stackData: stackData,
	}, nil
}

func (self *OutboundStackHandler002) Close() error {
	return self.stackData.Close()
}

func (self *OutboundStackHandler002) SendData(data interface{}) {
	if self.stackData.OutboundHandler != nil {
		self.stackData.OutboundHandler.OnSendData(data)
	}
}

func (self *OutboundStackHandler002) SendRws(rws goprotoextra.ReadWriterSize) {
	if self.errorState != nil {
		return
	}
	switch v := rws.(type) {
	case *gomessageblock.ReaderWriter:
		messageWrapper, err := stream.UnMarshal(v)
		if err != nil {
			return
		}
		switch vv := messageWrapper.(type) {
		case *wsmsg2.WebSocketMessageWrapper:
			switch vv.Data.OpCode {
			case wsmsg2.WebSocketMessage_OpText:
				err := wsutil.WriteClientMessage(self.stackData.UpgradedConnection, ws.OpText, vv.Data.Message)
				if err != nil {
					self.errorState = err
					return
				}
			case wsmsg2.WebSocketMessage_OpPing:
				binary, _ := time.Now().MarshalBinary()
				err := wsutil.WriteClientMessage(self.stackData.UpgradedConnection, ws.OpPing, binary)
				if err != nil {
					self.errorState = err
					return
				}
			}
		}
	}
}

func (self *OutboundStackHandler002) OnError(err error) {
	self.errorState = err
}

func (self *OutboundStackHandler002) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error,
) error {

	if self.errorState != nil {
		return self.errorState
	}
	self.SendRws(rws)

	return self.errorState
}

func (self *OutboundStackHandler002) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *OutboundStackHandler002) Complete() {

}
