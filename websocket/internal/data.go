package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/webSocketMessages/wsmsg"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goerrors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"io"
	"net"
	"net/url"
	"time"
)

type Data struct {
	Ctx                context.Context
	cancelFunc         context.CancelFunc
	Conn               net.Conn
	connectionType     model.ConnectionType
	UpgradedConnection *common.UpgradedConnectionWrapper
	LastPongReceived   time.Time
	connWrapper        *common.ConnWrapper
	pipeWriteClose     io.WriteCloser
	PipeRead           io.ReadCloser
	socketDataReceived int

	MultiOutBoundHandler     goCommsDefinitions.IRxNextHandler
	ConnectionReactorFactory intf.IConnectionReactorFactoryExtractValues
	OutboundHandler          goCommsDefinitions.IRxNextHandler
	additionalDataSend       int
	Logger                   *zap.Logger
	cfr                      intf.IConnectionReactorFactoryExtractValues
	goFunctionCounter        GoFunctionCounter.IService
	InboundHandler           goCommsDefinitions.IRxNextHandler
}

func (self *Data) Close() error {
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

func (self *Data) sendPing(_ context.Context) {
	msg := &wsmsg.WebSocketMessage{
		OpCode: wsmsg.WebSocketMessage_OpPing,
	}
	marshall, err := stream.Marshall(msg)
	if err != nil {
		return
	}
	self.additionalDataSend += marshall.Size()

	if self.MultiOutBoundHandler != nil {
		self.MultiOutBoundHandler.OnSendData(marshall)
	}
}

func (self *Data) triggerPingLoop(cancelContext context.Context) {
	ticker := time.NewTicker(time.Second * 30)
	defer func() {
		ticker.Stop()
	}()

	for true {
		select {
		case <-cancelContext.Done():
			return
		case <-ticker.C:
			if time.Now().Sub(self.LastPongReceived) > time.Minute*2 {
				self.Logger.Error("Close connection", zap.Error(fmt.Errorf("pong time out")))
				self.cancelFunc()
				return
			}
			self.sendPing(cancelContext)
		}
	}
}

func (self *Data) OnStart(
	inputStreamForStack common.IInputStreamForStack,
	Url *url.URL,
) (common.IInputStreamForStack, error) {
	if conn, ok := inputStreamForStack.(net.Conn); ok {
		if self.Ctx.Err() != nil {
			return nil, self.Ctx.Err()
		}

		// create map and fill it in with some values
		inputValues := make(map[string]interface{})
		inputValues["url"] = Url
		inputValues["localAddr"] = conn.LocalAddr()
		inputValues["remoteAddr"] = conn.RemoteAddr()
		outputValues, err := self.cfr.Values(inputValues)
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
		if self.Ctx.Err() != nil {
			return nil, self.Ctx.Err()
		}
		newConn, _, _, err := dialer.Dial(self.Ctx, Url.String())
		self.UpgradedConnection = common.NewUpgradedConnectionWrapper(newConn)
		// On error exit
		if err != nil {
			return nil, err
		}

		// On error exit
		if self.Ctx.Err() != nil {
			return nil, self.Ctx.Err()
		}

		_ = self.goFunctionCounter.GoRun("PingPong.data.Start.01",
			func() {
				self.connectionLoop(self.Ctx)
			},
		)
		_ = self.goFunctionCounter.GoRun("PingPong.data.Start.02",
			func() {
				self.triggerPingLoop(self.Ctx)
			},
		)

		return nil, self.Ctx.Err()
	}
	return nil, goerrors.InvalidParam
}

func (self *Data) connectionLoop(ctx context.Context) {
	sendMessage := func(message interface{}) {
		if self.InboundHandler != nil {
			self.InboundHandler.OnSendData(message)
		}
	}

	sendMessage(
		&wsmsg.WebSocketMessage{
			OpCode:  wsmsg.WebSocketMessage_OpStartLoop,
			Message: nil,
		},
	)

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
				sendMessage(
					&wsmsg.WebSocketMessage{
						OpCode:  wsmsg.WebSocketMessage_OpContinuation,
						Message: msg.Payload,
					},
				)
			case ws.OpText:
				sendMessage(
					&wsmsg.WebSocketMessage{
						OpCode:  wsmsg.WebSocketMessage_OpText,
						Message: msg.Payload,
					},
				)
			case ws.OpBinary:
				sendMessage(
					&wsmsg.WebSocketMessage{
						OpCode:  wsmsg.WebSocketMessage_OpBinary,
						Message: msg.Payload,
					},
				)
			case ws.OpClose:
				self.Logger.Debug("OpClose received")
				self.Logger.Error(
					"closing connection",
					zap.Error(fmt.Errorf("close notification received from socket")))
				//self.stackCancelFunc(
				//	"close webSocketMessage received",
				//	true,
				//	fmt.Errorf("close notification received from socket"),
				//)
			case ws.OpPing:
				self.Logger.Debug("Ping received")
				errWriteClientMessage := wsutil.WriteClientMessage(self.UpgradedConnection, ws.OpPong, msg.Payload)
				if errWriteClientMessage != nil {
					self.Logger.Error(
						"closing connection",
						zap.Error(errWriteClientMessage),
					)
					self.cancelFunc()
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

func (self *Data) SetConnWrapper(wrapper *common.ConnWrapper) error {
	self.connWrapper = wrapper
	return nil
}
