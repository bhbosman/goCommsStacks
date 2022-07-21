package internal

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	common3 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Inbound(
	ConnectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	opts ...rxgo.Option,
) func() (common3.IStackBoundDefinition, error) {
	return func() (common3.IStackBoundDefinition, error) {
		return common3.NewBoundDefinition(
				func(
					stackData common3.IStackCreateData,
					pipeData common3.IPipeCreateData,
					obs rxgo.Observable,
				) (rxgo.Observable, error) {
					if pd, ok := pipeData.(*InboundStackHandler); ok {
						mapHandler, err := RxHandlers.NewRxMapHandler(
							goCommsDefinitions.CompressionStackName,
							ConnectionCancelFunc,
							logger,
							pd,
						)
						if err != nil {
							return nil, err
						}

						return obs.Map(mapHandler.Handler, opts...), nil
					}
					return nil, goerrors.InvalidType
				},
				InboundPipeState(goCommsDefinitions.CompressionStackName)),
			nil

	}
}
