package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"io"
	"net"
	"reflect"
)

type Data struct {
	connectionType       model.ConnectionType
	Conn                 net.Conn
	ConnWrapper          *common.ConnWrapper
	PipeWriteClose       io.WriteCloser
	PipeRead             io.ReadCloser
	Ctx                  context.Context
	cancelFunc           context.CancelFunc
	connectionCancelFunc model.ConnectionCancelFunc
	UpgradedConnection   *common.UpgradedConnectionWrapper
	RxNextHandler        io.Closer
	Logger               *zap.Logger
	goFunctionCounter    GoFunctionCounter.IService
	OutboundHandler      goCommsDefinitions.IRxNextHandler
	InboundHandler       goCommsDefinitions.IRxNextHandler
}

func NewStackData(
	connectionType model.ConnectionType,
	Conn net.Conn,
	ctx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	goFunctionCounter GoFunctionCounter.IService,
) (*Data, error) {
	tempPipeRead, tempPipeWriteClose := common.Pipe(ctx)
	return &Data{
		connectionType:       connectionType,
		Conn:                 Conn,
		ConnWrapper:          nil,
		PipeWriteClose:       tempPipeWriteClose,
		PipeRead:             tempPipeRead,
		Ctx:                  ctx,
		cancelFunc:           cancelFunc,
		connectionCancelFunc: connectionCancelFunc,
		UpgradedConnection:   nil,
		Logger:               logger,
		goFunctionCounter:    goFunctionCounter,
	}, nil
}

func (self *Data) Close() error {
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

func (self *Data) Start() (net.Conn, error) {

	var tlsConn interface {
		net.Conn
		ConnectionState() tls.ConnectionState
	}
	if self.connectionType == model.ServerConnection {
		cer, err := tls.X509KeyPair(_serverPem, _serverKey)
		if err != nil {
			self.cancelFunc()
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
	self.UpgradedConnection = common.NewUpgradedConnectionWrapper(tlsConn)
	s := "tls.conn.read"

	// do not assign self.onInBoundComplete to onComplete as this will do a double close
	// the double close may be handled
	rxHandler, err := RxHandlers.NewRxNextHandler2(
		s,
		self.connectionCancelFunc,
		nil,
		self.InboundHandler, /*see comment*/
		self.Logger)
	if err != nil {
		return nil, err
	}

	self.RxNextHandler = rxHandler

	err = self.goFunctionCounter.GoRun("TlsConnection.data.Start",
		func() {
			//
			common.ReadFromIoReader(
				self.UpgradedConnection,
				self.Ctx,
				self.cancelFunc,
				rxHandler,
			)
		},
	)
	if err != nil {
		return nil, err
	}

	// this function is part of the GoFunctionCounter count
	return nil, self.Ctx.Err()
}

func (self *Data) SetConnWrapper(wrapper *common.ConnWrapper) error {
	self.ConnWrapper = wrapper
	return nil
}

func WrongStackDataError(connectionType model.ConnectionType, stackData interface{}) error {
	return common.NewWrongStackDataType(
		goCommsDefinitions.TlsStackName,
		connectionType,
		reflect.TypeOf((*Data)(nil)),
		reflect.TypeOf(stackData))
}
