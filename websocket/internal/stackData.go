package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/common/webSocketMessages/wsmsg"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goerrors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"io"
	"net"
	"net/url"
	"time"
)

type StackData struct {
	ctx                      context.Context
	conn                     net.Conn
	connectionType           model.ConnectionType
	UpgradedConnection       *common.UpgradedConnectionWrapper
	stackCancelFunc          model.ConnectionCancelFunc
	LastPongReceived         time.Time
	connWrapper              *common.ConnWrapper
	pipeWriteClose           io.WriteCloser
	pipeRead                 io.ReadCloser
	socketDataReceived       int
	onInBoundSendData        rxgo.NextFunc
	onInBoundSendError       rxgo.ErrFunc
	onInBoundComplete        rxgo.CompletedFunc
	onMultiOutBoundSendData  rxgo.NextFunc
	onMultiOutBoundSendError rxgo.ErrFunc
	onMultiOutBoundComplete  rxgo.CompletedFunc

	onOutBoundSendData  rxgo.NextFunc
	onOutBoundSendError rxgo.ErrFunc
	onOutBoundComplete  rxgo.CompletedFunc
	additionalDataSend  int
	Logger              *zap.Logger
}

func (self *StackData) Close() error {
	var err error = nil
	if self.UpgradedConnection != nil {
		err = multierr.Append(err, self.UpgradedConnection.Close())
	}
	if self.pipeWriteClose != nil {
		err = multierr.Append(err, self.pipeWriteClose.Close())
	}
	if self.connWrapper != nil {
		err = multierr.Append(err, self.connWrapper.Close())
	}
	return err
}

func (self *StackData) sendPing(_ context.Context) {
	msg := &wsmsg.WebSocketMessage{
		OpCode: wsmsg.WebSocketMessage_OpPing,
	}
	marshall, err := stream.Marshall(msg)
	if err != nil {
		return
	}
	self.additionalDataSend += marshall.Size()

	if self.onMultiOutBoundSendData != nil {
		self.onMultiOutBoundSendData(marshall)
	}
}

func (self *StackData) triggerPingLoop(cancelContext context.Context) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for true {
		select {
		case <-cancelContext.Done():
			return
		case <-ticker.C:
			if time.Now().Sub(self.LastPongReceived) > time.Minute*2 {
				self.stackCancelFunc("pong time out", true, goerrors.TimeOut)
				return
			}
			self.sendPing(cancelContext)
		}
	}
}

func (self *StackData) OnStart(
	inputStreamForStack common.IInputStreamForStack,
	Url *url.URL,
	Ctx context.Context,
	ConnectionReactorFactory intf.IConnectionReactorFactoryExtractValues) (common.IInputStreamForStack, error) {
	if conn, ok := inputStreamForStack.(net.Conn); ok {
		if Ctx.Err() != nil {
			return nil, Ctx.Err()
		}

		// create map and fill it in with some values
		inputValues := make(map[string]interface{})
		inputValues["url"] = Url
		inputValues["localAddr"] = conn.LocalAddr()
		inputValues["remoteAddr"] = conn.RemoteAddr()
		outputValues, err := ConnectionReactorFactory.Values(inputValues)
		if err != nil {
			return nil, err
		}

		// get header information from outputValues
		header := make(ws.HandshakeHeaderHTTP)
		if additionalHeaderInformation, ok := outputValues["connectionHeader"]; ok {
			if connectionHeader, isMap := additionalHeaderInformation.(map[string][]string); isMap {
				for k, v := range connectionHeader {
					header[k] = v
				}
			}
		}

		// build websocket dialer that will be used
		dialer := ws.Dialer{
			Header: header,
			NetDial: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				return self.connWrapper, nil
			},
		}
		// On error exit
		if Ctx.Err() != nil {
			return nil, Ctx.Err()
		}
		newConn, _, _, err := dialer.Dial(Ctx, Url.String())
		self.UpgradedConnection = common.NewUpgradedConnectionWrapper(newConn)
		// On error exit
		if err != nil {
			return nil, err
		}

		// On error exit
		if Ctx.Err() != nil {
			return nil, Ctx.Err()
		}

		go self.connectionLoop(Ctx)
		go self.triggerPingLoop(Ctx)
		return nil, Ctx.Err()
	}
	return nil, goerrors.InvalidParam
}

func (self *StackData) connectionLoop(ctx context.Context) {
	sendMessage := func(message *wsmsg.WebSocketMessage) {
		stm, err := stream.Marshall(message)
		if err != nil {
			//return
		}
		if self.onInBoundSendData != nil {
			self.socketDataReceived += stm.Size()
			self.onInBoundSendData(stm)
		}
	}
	message := wsmsg.WebSocketMessage{
		OpCode:  wsmsg.WebSocketMessage_OpStartLoop,
		Message: nil,
	}
	sendMessage(&message)

	for {
		messages, err := wsutil.ReadServerMessage(self.UpgradedConnection, nil)
		if err != nil {
			return
		}
		if ctx.Err() != nil {
			return
		}
		for _, msg := range messages {
			if ctx.Err() != nil {
				return
			}
			switch msg.OpCode {
			case ws.OpContinuation:
				self.Logger.Debug("OpContinuation received")
				webSocketMessage := wsmsg.WebSocketMessage{
					OpCode:  wsmsg.WebSocketMessage_OpContinuation,
					Message: msg.Payload,
				}
				sendMessage(&webSocketMessage)
			case ws.OpText:
				self.Logger.Debug("OpText received")
				webSocketMessage := wsmsg.WebSocketMessage{
					OpCode:  wsmsg.WebSocketMessage_OpText,
					Message: msg.Payload,
				}
				sendMessage(&webSocketMessage)
			case ws.OpBinary:
				self.Logger.Debug("OpContinuation received")
				webSocketMessage := wsmsg.WebSocketMessage{
					OpCode:  wsmsg.WebSocketMessage_OpBinary,
					Message: msg.Payload,
				}
				sendMessage(&webSocketMessage)
			case ws.OpClose:
				self.Logger.Debug("OpClose received")
				self.stackCancelFunc(
					"close webSocketMessage received",
					true,
					fmt.Errorf("close notification received from socket"))
			case ws.OpPing:
				self.Logger.Debug("Ping received")
				errWriteClientMessage := wsutil.WriteClientMessage(self.UpgradedConnection, ws.OpPong, msg.Payload)
				if errWriteClientMessage != nil {
					self.stackCancelFunc("creating pong payload", true, err)
					return
				}
				continue
			case ws.OpPong:
				now := time.Now()
				_ = self.LastPongReceived.UnmarshalBinary(msg.Payload)
				self.Logger.Debug("Pong received",
					zap.Time("send", self.LastPongReceived),
					zap.Time("Received", now),
					zap.Duration("RoundTrip", now.Sub(self.LastPongReceived)))
				continue
			default:
				continue
			}
			if ctx.Err() != nil {
				return
			}
		}
	}
}

func (self *StackData) setOnInBoundSendData(onData rxgo.NextFunc) error {
	if onData == nil {
		return goerrors.InvalidParam
	}
	self.onInBoundSendData = onData
	return nil
}

func (self *StackData) setOnInBoundSendError(onSendError rxgo.ErrFunc) error {
	if onSendError == nil {
		return goerrors.InvalidParam
	}
	self.onInBoundSendError = onSendError
	return nil
}

func (self *StackData) SetOnInBoundComplete(onComplete rxgo.CompletedFunc) error {
	if onComplete == nil {
		return goerrors.InvalidParam
	}
	self.onInBoundComplete = onComplete
	return nil
}

func (self *StackData) setConnWrapper(wrapper *common.ConnWrapper) error {
	self.connWrapper = wrapper
	return nil
}
