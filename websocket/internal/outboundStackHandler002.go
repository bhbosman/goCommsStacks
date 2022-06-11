package internal

import (
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/RxHandlers"
	wsmsg2 "github.com/bhbosman/gocomms/common/webSocketMessages/wsmsg"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"time"
)

type OutboundStackHandler002 struct {
	errorState error
	stackData  *StackData
}

func (self *OutboundStackHandler002) GetAdditionalBytesIncoming() int {
	return -self.stackData.additionalDataSend
}

func (self *OutboundStackHandler002) GetAdditionalBytesSend() int {
	return self.stackData.connWrapper.BytesWritten
}

func (self *OutboundStackHandler002) ReadMessage(_ interface{}) error {
	return nil
}

func NewOutboundStackHandler002(stackData *StackData) (*OutboundStackHandler002, error) {
	return &OutboundStackHandler002{
		stackData: stackData,
	}, nil
}

func (self *OutboundStackHandler002) Close() error {
	return self.stackData.Close()
}

func (self *OutboundStackHandler002) SendData(data interface{}) {
	if self.stackData.onOutBoundSendData != nil {
		self.stackData.onOutBoundSendData(data)
	}
}

func (self *OutboundStackHandler002) SendRws(size goprotoextra.ReadWriterSize) {
	if self.errorState != nil {
		return
	}
	switch v := size.(type) {
	case *gomessageblock.ReaderWriter:
		messageWrapper, err := stream.UnMarshal(v, nil, nil, nil, nil)
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
	size goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(size int) error) error {

	if self.errorState != nil {
		return self.errorState
	}
	self.SendRws(size)

	return self.errorState
}

func (self *OutboundStackHandler002) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *OutboundStackHandler002) Complete() {

}
