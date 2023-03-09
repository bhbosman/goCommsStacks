package internal

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goerrors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"net"
	"time"
)

func NewStackData(
	logger *zap.Logger,
	connectionType model.ConnectionType,
	conn net.Conn,
	//stackCancelFunc model.ConnectionCancelFunc,
	ctx context.Context,
	cancelFunc context.CancelFunc,
	cfr intf.IConnectionReactorFactoryExtractValues,
	goFunctionCounter GoFunctionCounter.IService,
) (*Data, error) {
	var errorList error
	if logger == nil {
		errorList = multierr.Append(errorList, goerrors.InvalidParam)
	}
	if conn == nil {
		errorList = multierr.Append(errorList, goerrors.InvalidParam)
	}
	//if stackCancelFunc == nil {
	//	errorList = multierr.Append(errorList, goerrors.InvalidParam)
	//}
	if ctx == nil {
		errorList = multierr.Append(errorList, goerrors.InvalidParam)
	}
	if cfr == nil {
		errorList = multierr.Append(errorList, goerrors.InvalidParam)
	}
	if errorList != nil {
		return nil, errorList
	}

	pipeRead, pipeWriteCloser := common.Pipe(ctx)

	return &Data{
		Ctx:                ctx,
		cancelFunc:         cancelFunc,
		Conn:               conn,
		connectionType:     connectionType,
		UpgradedConnection: nil,
		//stackCancelFunc:    stackCancelFunc,
		LastPongReceived:  time.Now(),
		connWrapper:       nil,
		pipeWriteClose:    pipeWriteCloser,
		PipeRead:          pipeRead,
		Logger:            logger,
		cfr:               cfr,
		goFunctionCounter: goFunctionCounter,
	}, nil
}
