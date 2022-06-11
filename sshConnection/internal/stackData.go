package internal

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
	"io"
	"net"
)

type StackData struct {
	connectionType      model.ConnectionType
	sshConn             ssh.Conn
	conn                net.Conn
	ctx                 context.Context
	logger              *zap.Logger
	onInBoundSendData   rxgo.NextFunc
	onInBoundSendError  rxgo.ErrFunc
	onInBoundComplete   rxgo.CompletedFunc
	onOutBoundSendData  rxgo.NextFunc
	onOutBoundSendError rxgo.ErrFunc
	onOutBoundComplete  rxgo.CompletedFunc
	PipeWriteClose      io.WriteCloser
	PipeRead            io.ReadCloser
	ConnWrapper         *common.ConnWrapper
}

func (self *StackData) Close() error {
	var errList error = nil

	if self.conn != nil {
		errList = multierr.Append(errList, self.conn.Close())
	} else {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}

	if self.sshConn != nil {
		errList = multierr.Append(errList, self.sshConn.Close())
	} else {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}

	if self.PipeWriteClose != nil {
		errList = multierr.Append(errList, self.PipeWriteClose.Close())
	} else {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}

	if self.ConnWrapper != nil {
		errList = multierr.Append(errList, self.ConnWrapper.Close())
	} else {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}

	return errList
}

func (self *StackData) Start(Ctx context.Context) (common.IInputStreamForStack, error) {
	if self.connectionType == model.ClientConnection {
		return nil, goerrors.InvalidParam
	}
	config := &ssh.ServerConfig{
		Config: ssh.Config{
			Rand:           nil,
			RekeyThreshold: 0,
			KeyExchanges:   nil,
			Ciphers:        nil,
			MACs:           nil,
		},
		NoClientAuth:     false,
		MaxAuthTries:     0,
		PasswordCallback: nil,
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return &ssh.Permissions{
				Extensions: map[string]string{
					"key-id": "sssssss",
				},
			}, nil
		},
		KeyboardInteractiveCallback: nil,
		AuthLogCallback:             nil,
		ServerVersion:               "",
		BannerCallback:              nil,
		GSSAPIWithMICConfig:         nil,
	}

	key, err := ssh.ParsePrivateKey(_defaultcert)
	if err != nil {
		return nil, err
	}
	config.AddHostKey(key)

	serverConn, channels, requests, err := ssh.NewServerConn(self.ConnWrapper, config)
	if err != nil {
		return nil, err
	}

	go self.handleChannelsAndRequest(channels, requests)
	self.sshConn = serverConn

	if err != nil {
		return nil, err
	}

	return nil, Ctx.Err()
}

func (self *StackData) Stop() error {
	return nil
}

func (self *StackData) handleChannelsAndRequest(channels <-chan ssh.NewChannel, requests <-chan *ssh.Request) {
loop:
	for self.ctx.Err() == nil {
		select {
		case <-self.ctx.Done():
			break loop
		case channel, ok := <-channels:
			if !ok {
				break loop
			}
			err := channel.Reject(ssh.ResourceShortage, "ssh.ResourceShortage")
			if err != nil {
				return
			}
			//accept, i, err := channel.Accept()
			//go func() {
			//	accept.CloseWrite()
			//	for request := range i {
			//		accept.Write()
			//	}
			//}()

			break
		case _, ok := <-requests:
			if !ok {
				break loop
			}
			break
		}
	}

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

func (self *StackData) setConnWrapper(wrapper *common.ConnWrapper) error {
	self.ConnWrapper = wrapper
	return nil
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

func NewStackData(
	connectionType model.ConnectionType,
	conn net.Conn,
	ctx context.Context,
	logger *zap.Logger) *StackData {
	tempPipeRead, tempPipeWriteClose := common.Pipe(ctx)

	return &StackData{
		connectionType: connectionType,
		conn:           conn,
		ctx:            ctx,
		logger:         logger,
		PipeWriteClose: tempPipeWriteClose,
		PipeRead:       tempPipeRead,
	}
}
