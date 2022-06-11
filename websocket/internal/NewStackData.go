package internal

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/zap"
	"net"
	"time"
)

func NewStackData(
	logger *zap.Logger,
	connectionType model.ConnectionType,
	conn net.Conn,
	stackCancelFunc model.ConnectionCancelFunc,
	ctx context.Context) (*StackData, error) {

	pipeRead, pipeWriteCloser := common.Pipe(ctx)

	return &StackData{
		ctx:                ctx,
		conn:               conn,
		connectionType:     connectionType,
		UpgradedConnection: nil,
		stackCancelFunc:    stackCancelFunc,
		LastPongReceived:   time.Now(),
		connWrapper:        nil,
		pipeWriteClose:     pipeWriteCloser,
		pipeRead:           pipeRead,
		Logger:             logger,
	}, nil
}
