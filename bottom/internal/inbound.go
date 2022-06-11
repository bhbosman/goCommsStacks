package internal

import (
	"github.com/bhbosman/goCommsStacks/bottom/common"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func Inbound(ConnectionCancelFunc model.ConnectionCancelFunc, logger *zap.Logger, opts ...rxgo.Option) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(stackData common2.IStackCreateData, pipeData common2.IPipeCreateData, obs rxgo.Observable) (string, rxgo.Observable, error) {
					if pipeData != nil {
						return "", nil, goerrors.InvalidParam
					}
					if stackData != nil {
						return "", nil, goerrors.InvalidParam
					}
					var errorList error
					stackHandler, err := NewInboundStackHandler(logger)
					if err != nil {
						errorList = multierr.Append(errorList, err)
					}
					mapHandler, err := RxHandlers.NewRxMapHandler(common.StackName, ConnectionCancelFunc, logger, stackHandler)
					if err != nil {
						errorList = multierr.Append(errorList, err)
					}
					if errorList != nil {
						return common.StackName, nil, errorList
					}
					return common.StackName, obs.Map(mapHandler.Handler, opts...), nil
				},
				nil),
			nil
	}
}