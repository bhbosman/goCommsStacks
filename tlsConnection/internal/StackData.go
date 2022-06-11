package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/bhbosman/goCommsStacks/tlsConnection/common"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"io"
	"net"
	"reflect"
)

type StackData struct {
	connectionType      model.ConnectionType
	Conn                net.Conn
	ConnWrapper         *common2.ConnWrapper
	PipeWriteClose      io.WriteCloser
	PipeRead            io.ReadCloser
	ctx                 context.Context
	UpgradedConnection  *common2.UpgradedConnectionWrapper
	RxNextHandler       io.Closer
	onInBoundSendData   rxgo.NextFunc
	onInBoundSendError  rxgo.ErrFunc
	onInBoundComplete   rxgo.CompletedFunc
	onOutBoundSendData  rxgo.NextFunc
	onOutBoundSendError rxgo.ErrFunc
	onOutBoundComplete  rxgo.CompletedFunc
	Logger              *zap.Logger
}

func NewStackData(
	connectionType model.ConnectionType,
	Conn net.Conn,
	ctx context.Context,
	logger *zap.Logger) (*StackData, error) {
	tempPipeRead, tempPipeWriteClose := common2.Pipe(ctx)
	return &StackData{
		ctx:                ctx,
		connectionType:     connectionType,
		Conn:               Conn,
		ConnWrapper:        nil,
		PipeWriteClose:     tempPipeWriteClose,
		UpgradedConnection: nil,
		PipeRead:           tempPipeRead,
		Logger:             logger,
	}, nil
}

func (self *StackData) Close() error {
	var err error = nil
	err = multierr.Append(err, self.ConnWrapper.Close())
	err = multierr.Append(err, self.PipeWriteClose.Close())
	if self.UpgradedConnection != nil {
		err = multierr.Append(err, self.UpgradedConnection.Close())
	}
	if self.RxNextHandler != nil {
		err = multierr.Append(err, self.RxNextHandler.Close())
	}
	return err
}

func (self *StackData) Start(
	connectionType model.ConnectionType,
	ConnectionCancelFunc model.ConnectionCancelFunc,
	Ctx context.Context,
	ctxCancel context.CancelFunc) (net.Conn, error) {

	var tlsConn net.Conn
	if connectionType == model.ServerConnection {
		cer, err := tls.X509KeyPair(_serverPem, _serverKey)
		if err != nil {
			ctxCancel()
			return nil, err
		}
		config := &tls.Config{
			ServerName:   "127.0.0.1",
			Certificates: []tls.Certificate{cer},
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				return nil
			},
			VerifyConnection: func(state tls.ConnectionState) error {
				return nil
			},
		}
		tlsConn = tls.Server(self.ConnWrapper, config)
	} else {
		config := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "localhost"}
		tlsConn = tls.Client(self.ConnWrapper, config)
	}
	self.UpgradedConnection = common2.NewUpgradedConnectionWrapper(tlsConn)
	s := "tls.conn.read"

	// do not assign self.onInBoundComplete to onComplete as this will do a double close
	// the double close may be handled
	rxHandler, err := RxHandlers.NewRxNextHandler(
		s,
		ConnectionCancelFunc,
		nil,
		self.onInBoundSendData,
		self.onInBoundSendError,
		nil, /*see comment*/
		self.Logger)
	if err != nil {
		return nil, err
	}

	self.RxNextHandler = rxHandler
	go common2.ReadFromIoReader(s, self.UpgradedConnection, Ctx, ConnectionCancelFunc, rxHandler)
	return nil, Ctx.Err()
}

func (self *StackData) setOnInBoundSendData(onData rxgo.NextFunc) error {
	if onData == nil {
		return goerrors.InvalidParam
	}
	self.onInBoundSendData = func(i interface{}) {
		onData(i)
	}
	return nil
}

func (self *StackData) setOnInBoundSendError(onSendError rxgo.ErrFunc) error {
	if onSendError == nil {
		return goerrors.InvalidParam
	}
	self.onInBoundSendError = func(err error) {
		onSendError(err)
	}
	return nil
}

func (self *StackData) setOnInBoundComplete(onComplete rxgo.CompletedFunc) error {
	if onComplete == nil {
		return goerrors.InvalidParam
	}
	self.onInBoundComplete = func() {
		onComplete()
	}
	return nil
}

func (self *StackData) setOnOutBoundSendData(onData rxgo.NextFunc) error {
	if onData == nil {
		return goerrors.InvalidParam
	}
	self.onOutBoundSendData = func(i interface{}) {
		onData(i)
	}
	return nil
}

func (self *StackData) setOnOutBoundSendError(onSendError rxgo.ErrFunc) error {
	if onSendError == nil {
		return goerrors.InvalidParam
	}
	self.onOutBoundSendError = func(err error) {
		onSendError(err)
	}
	return nil
}

func (self *StackData) setOnOutBoundComplete(onComplete rxgo.CompletedFunc) error {
	if onComplete == nil {
		return goerrors.InvalidParam
	}
	self.onOutBoundComplete = func() {
		onComplete()
	}
	return nil
}

func (self *StackData) setConnWrapper(wrapper *common2.ConnWrapper) error {
	self.ConnWrapper = wrapper
	return nil
}

func WrongStackDataError(connectionType model.ConnectionType, stackData interface{}) error {
	return common2.NewWrongStackDataType(
		common.StackName,
		connectionType,
		reflect.TypeOf((*StackData)(nil)),
		reflect.TypeOf(stackData))
}
