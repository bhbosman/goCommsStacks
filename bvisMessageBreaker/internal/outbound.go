package internal

import (
	"github.com/bhbosman/goCommsStacks/bvisMessageBreaker/common"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Outbound(ConnectionCancelFunc model.ConnectionCancelFunc, marker [4]byte, logger *zap.Logger, opts ...rxgo.Option) common2.BoundResult {
	return func() (common2.IStackBoundDefinition, error) {
		return common2.NewBoundDefinition(
				func(stackData common2.IStackCreateData, pipeData common2.IPipeCreateData, obs rxgo.Observable) (string, rxgo.Observable, error) {
					outboundStackHandler, err := NewOutboundStackHandler(marker)
					if err != nil {
						return common.StackName, nil, err
					}

					mapHandler, err := RxHandlers.NewRxMapHandler(common.StackName, ConnectionCancelFunc, logger, outboundStackHandler)
					if err != nil {
						return common.StackName, nil, err
					}

					return common.StackName, obs.Map(mapHandler.Handler, opts...), nil
				}, nil),
			nil
	}
}
