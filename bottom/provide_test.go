package bottom

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gomessageblock"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"testing"
)

func TestProvideBottomStack(t *testing.T) {
	createFxAppFactory := func(logger *zap.Logger) (*common.StackDefinition, error) {
		out := struct {
			fx.In
			Data []*common.StackDefinition `group:"StackDefinition"`
		}{}

		fxApp := fxtest.New(t,
			fx.Provide(common.ProvideRxOptions),
			fx.Provide(func() (context.Context, context.CancelFunc) { return context.WithCancel(context.Background()) }),
			fx.Provide(func(cancelFunc context.CancelFunc) model.ConnectionCancelFunc {
				return func(context string, inbound bool, err error) { cancelFunc() }
			}),
			fx.Provide(func() *zap.Logger { return logger }),
			ProvideBottomStack(),
			fx.Populate(&out))
		assert.NoError(t, fxApp.Err())
		return out.Data[0], fxApp.Err()
	}
	standardBottomTest := func(
		inboundFactory func(logger *zap.Logger) (common.IBoundResult, error)) func(t *testing.T) {
		return func(t *testing.T) {
			inboundInstance, err := inboundFactory(zap.NewNop())
			assert.NoError(t, err)
			t.Run("Pipe with 1 message", func(t *testing.T) {
				result, _ := inboundInstance.GetBoundResult()
				definition, _ := result()
				ch := make(chan rxgo.Item)
				channel := rxgo.FromChannel(ch)
				_, outCh, err := definition.GetPipeDefinition()(nil, nil, channel)
				assert.NoError(t, err)
				assert.NotNil(t, outCh)

				cancelCtx, cancelFunc := context.WithCancel(context.Background())
				dataOut := 0
				dataIn := 0
				go func(cancelFunc context.CancelFunc) {
					for item := range channel.Observe() {
						switch v := item.V.(type) {
						case *gomessageblock.ReaderWriter:
							dataOut += v.Size()
						}
					}
					cancelFunc()
				}(cancelFunc)

				rws, _ := gomessageblock.NewReaderWriterString("Hello World")
				dataIn += rws.Size()
				ch <- rxgo.Of(rws)
				close(ch)
				<-cancelCtx.Done()
				assert.Equal(t, dataIn, dataOut)
			})
		}
	}

	t.Run(
		"Inbound",
		standardBottomTest(
			func(logger *zap.Logger) (common.IBoundResult, error) {
				instance, err := createFxAppFactory(logger)
				if err != nil {
					return nil, err
				}
				return instance.Inbound, nil
			}))
	t.Run(
		"Outbound",
		standardBottomTest(
			func(logger *zap.Logger) (common.IBoundResult, error) {
				instance, err := createFxAppFactory(logger)
				if err != nil {
					return nil, err
				}
				return instance.Outbound, nil
			}))

	t.Run("Inbound-Logger = nil", func(t *testing.T) {
		stackDefinition, err := createFxAppFactory(nil)
		assert.NoError(t, err)
		result, _ := stackDefinition.Inbound.GetBoundResult()
		definition, _ := result()
		ch := make(chan rxgo.Item)
		channel := rxgo.FromChannel(ch)
		_, _, err = definition.GetPipeDefinition()(nil, nil, channel)
		assert.Error(t, err)
		close(ch)
	})

	t.Run("Outbound-Logger = nil", func(t *testing.T) {
		stackDefinition, err := createFxAppFactory(nil)
		assert.NoError(t, err)
		result, _ := stackDefinition.Outbound.GetBoundResult()
		definition, _ := result()
		ch := make(chan rxgo.Item)
		channel := rxgo.FromChannel(ch)
		_, _, err = definition.GetPipeDefinition()(nil, nil, channel)
		assert.Error(t, err)
		close(ch)
	})
}
