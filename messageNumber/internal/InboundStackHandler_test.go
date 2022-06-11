package internal

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/gomessageblock"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestInboundStackHandler_MapReadWriterSizeWithDataAndNoError(t *testing.T) {
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
	t.Run("One Message", func(t *testing.T) {
		stackHandler, _ := NewInboundStackHandler()
		s := "HelloWorld"
		byteArray := buildMessage(1, s, nil)
		rws := gomessageblock.NewReaderWriterBlock(byteArray)
		result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
		assert.NoError(t, err)
		sb := strings.Builder{}
		_, _ = io.Copy(&sb, result)
		assert.Equal(t, s, sb.String())
	})

	t.Run("Two Message", func(t *testing.T) {
		stackHandler, _ := NewInboundStackHandler()
		s := "HelloWorld"
		byteArray := buildMessage(1, s, nil)
		rws := gomessageblock.NewReaderWriterBlock(byteArray)
		result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
		assert.NoError(t, err)
		sb := strings.Builder{}
		_, _ = io.Copy(&sb, result)
		assert.Equal(t, s, sb.String())

		s = "HelloWorld2"
		byteArray = buildMessage(2, s, nil)
		rws = gomessageblock.NewReaderWriterBlock(byteArray)
		result, err = stackHandler.MapReadWriterSize(context.Background(), rws)
		assert.NoError(t, err)
		sb = strings.Builder{}
		_, _ = io.Copy(&sb, result)
		assert.Equal(t, s, sb.String())
	})

	t.Run("out of seq", func(t *testing.T) {
		stackHandler, _ := NewInboundStackHandler()
		s := "HelloWorld"
		byteArray := buildMessage(4, s, nil)
		rws := gomessageblock.NewReaderWriterBlock(byteArray)
		result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

}

func TestInboundStackHandler_MapReadWriterSizeWithDataAndError(t *testing.T) {
	stackHandler, _ := NewInboundStackHandler()
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	result, err := stackHandler.MapReadWriterSize(cancel, rws)
	assert.Nil(t, result)
	assert.Error(t, err)
}
