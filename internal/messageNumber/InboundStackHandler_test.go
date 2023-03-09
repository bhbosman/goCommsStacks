package messageNumber_test

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsStacks/internal/messageNumber"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gomessageblock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"strings"
	"testing"
)

func TestWithMapHandler(t *testing.T) {
	buildMessage := func(number uint64, s string, hash io.Writer) []byte {
		stringAsBytes := []byte(s)
		if hash != nil {
			_, _ = hash.Write(stringAsBytes)
		}
		result := make([]byte, len(stringAsBytes)+8)

		binary.LittleEndian.PutUint64(result[0:8], number)
		copy(result[8:len(stringAsBytes)+8], stringAsBytes)
		return result
	}
	createRws := func(s string) *gomessageblock.ReaderWriter {
		byteArray := buildMessage(1, s, nil)
		return gomessageblock.NewReaderWriterBlock(byteArray)
	}
	t.Run(
		"Nested",
		func(t *testing.T) {
			mapHandlerFactory := func() *RxHandlers.RxMapHandler {
				stackHandler, err := messageNumber.NewInboundStackHandler()
				assert.NoError(t, err)
				assert.NotNil(t, stackHandler)

				mapHandler, err := RxHandlers.NewRxMapHandler(
					goCommsDefinitions.MessageNumberStackName,
					func(context string, inbound bool, err error) {

					},
					zap.NewNop(),
					stackHandler,
				)
				assert.NoError(t, err)
				assert.NotNil(t, mapHandler)
				return mapHandler
			}
			t.Run(
				"d",
				func(t *testing.T) {
					ctx := context.Background()
					pipe01 := mapHandlerFactory()
					pipe02 := mapHandlerFactory()
					s := "Hello World"
					result, err := pipe01.Handler(ctx, createRws(s))
					assert.NoError(t, err)
					if rws, ok := result.(*gomessageblock.ReaderWriter); ok {
						err = rws.Skip(len(s))
						assert.NoError(t, err)
					}
					result, err = pipe02.Handler(ctx, result)
					assert.NoError(t, err)
				},
			)
		},
	)

	t.Run(
		"Standard Inbound Tests",
		func(t *testing.T) {
			t.Run(
				"One Message",
				func(t *testing.T) {
					stackHandler, _ := messageNumber.NewInboundStackHandler()
					s := "HelloWorld"
					rws := createRws(s)
					result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
					assert.NoError(t, err)
					sb := strings.Builder{}
					_, _ = io.Copy(&sb, result.(io.Reader))
					assert.Equal(t, s, sb.String())
				},
			)

			t.Run(
				"Two Message",
				func(t *testing.T) {
					stackHandler, _ := messageNumber.NewInboundStackHandler()
					s := "HelloWorld"
					byteArray := buildMessage(1, s, nil)
					rws := gomessageblock.NewReaderWriterBlock(byteArray)
					result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
					assert.NoError(t, err)
					sb := strings.Builder{}
					_, _ = io.Copy(&sb, result.(io.Reader))
					assert.Equal(t, s, sb.String())

					s = "HelloWorld2"
					byteArray = buildMessage(2, s, nil)
					rws = gomessageblock.NewReaderWriterBlock(byteArray)
					result, err = stackHandler.MapReadWriterSize(context.Background(), rws)
					assert.NoError(t, err)
					sb = strings.Builder{}
					_, _ = io.Copy(&sb, result.(io.Reader))
					assert.Equal(t, s, sb.String())
				},
			)

			t.Run(
				"out of seq",
				func(t *testing.T) {
					stackHandler, _ := messageNumber.NewInboundStackHandler()
					s := "HelloWorld"
					byteArray := buildMessage(4, s, nil)
					rws := gomessageblock.NewReaderWriterBlock(byteArray)
					result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
					assert.Error(t, err)
					assert.Nil(t, result)
				},
			)
		},
	)

}
