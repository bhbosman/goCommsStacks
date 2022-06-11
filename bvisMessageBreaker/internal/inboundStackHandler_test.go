package internal

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestInboundStackHandler_MapReadWriterSizeWithDataAndNoError(t *testing.T) {
	stackHandler := NewInboundStackHandler([4]byte{1, 2, 3, 4})
	t.Run("OneMessage", func(t *testing.T) {
		s := "HelloWorld"
		inputData := gomessageblock.NewReaderWriter()
		addToInputData := CreateMultiByteArrayMessages([4]byte{1, 2, 3, 4}, inputData)
		err := addToInputData([]byte(s))
		if !assert.NoError(t, err) {
			return
		}
		var result goprotoextra.ReadWriterSize
		err = stackHandler.NextReadWriterSize(
			inputData,
			func(rws goprotoextra.ReadWriterSize) error {
				result = rws
				return nil
			},
			func(size int) error {
				return nil
			})
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, result) {
			return
		}
		sb := strings.Builder{}
		_, _ = io.Copy(&sb, result)
		assert.Equal(t, s, sb.String())

	})

	t.Run("FourMessage", func(t *testing.T) {
		ss := [4]string{
			"HelloWorld1",
			"HelloWorld2",
			"HelloWorld3",
			"HelloWorld4",
		}

		inputData := gomessageblock.NewReaderWriter()
		addToInputData := CreateMultiByteArrayMessages([4]byte{1, 2, 3, 4}, inputData)
		for _, s := range ss {
			err := addToInputData([]byte(s))
			if !assert.NoError(t, err) {
				return
			}
		}
		callback := func() func(rws goprotoextra.ReadWriterSize) error {
			index := 0
			return func(rws goprotoextra.ReadWriterSize) error {
				sb := strings.Builder{}
				_, _ = io.Copy(&sb, rws)
				if ss[index] != sb.String() {
					return fmt.Errorf("wrong value")
				}
				index++
				return nil
			}
		}
		cb := callback()
		err := stackHandler.NextReadWriterSize(inputData, cb, func(size int) error { return nil })
		if !assert.NoError(t, err) {
			return
		}
	})
}

func TestInboundStackHandler_MapReadWriterSizeWithDataAndError(t *testing.T) {
	stackHandler := NewInboundStackHandler(
		[4]byte{1, 2, 3, 4})
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	stackHandler.OnError(cancel.Err())
	var result goprotoextra.ReadWriterSize
	_ = stackHandler.NextReadWriterSize(
		rws,
		func(rws goprotoextra.ReadWriterSize) error {
			result = rws
			return nil
		},
		func(size int) error {
			return nil
		})
	assert.Nil(t, result)
	assert.Error(t, cancel.Err())
}

func TestInboundStackHandler_NextReadWriterSizeBrokenMessage(t *testing.T) {
	marker := [4]byte{1, 2, 3, 4}

	buildMessage := func(marker [4]byte, s string, hash io.Writer) []byte {
		stringAsBytes := []byte(s)
		if hash != nil {
			_, _ = hash.Write(stringAsBytes)
		}
		result := make([]byte, len(stringAsBytes)+8)

		copy(result[0:4], marker[:])
		binary.LittleEndian.PutUint32(result[4:8], uint32(len(stringAsBytes)))
		copy(result[8:len(stringAsBytes)+8], stringAsBytes)
		return result
	}
	t.Run("One Complete Message", func(t *testing.T) {
		stackHandler := NewInboundStackHandler(
			marker)
		s := "HelloWorld"
		message := buildMessage(marker, s, nil)
		rws := gomessageblock.NewReaderWriterBlock(message)
		var result goprotoextra.ReadWriterSize
		_ = stackHandler.NextReadWriterSize(
			rws,
			func(rws goprotoextra.ReadWriterSize) error {
				result = rws
				return nil
			},
			func(size int) error {
				return nil
			})
		if !assert.NotNil(t, result) {
			return
		}
		sb := strings.Builder{}
		_, _ = io.Copy(&sb, result)
		assert.Equal(t, s, sb.String())
	})
	t.Run("One Complete Message in bytes", func(t *testing.T) {
		stackHandler := NewInboundStackHandler(
			marker)
		s := "HelloWorld"
		message := buildMessage(marker, s, nil)
		var result goprotoextra.ReadWriterSize

		for _, b := range message {
			rws := gomessageblock.NewReaderWriterBlock([]byte{b})
			_ = stackHandler.NextReadWriterSize(
				rws,
				func(rws goprotoextra.ReadWriterSize) error {
					if !assert.Nil(t, result) {
						return fmt.Errorf("SAdsadasdas")
					}
					result = rws
					return nil
				},
				func(size int) error {
					return nil
				})
		}
		if !assert.NotNil(t, result) {
			return
		}
		sb := strings.Builder{}
		_, _ = io.Copy(&sb, result)
		assert.Equal(t, s, sb.String())
	})
	t.Run("Wrong Signature", func(t *testing.T) {
		stackHandler := NewInboundStackHandler(marker)
		wrongMarker := [4]byte{255, 255, 255, 255}
		message := buildMessage(wrongMarker, "With Wrong Market", nil)
		rws := gomessageblock.NewReaderWriterBlock(message)
		err := stackHandler.NextReadWriterSize(
			rws,
			func(rws goprotoextra.ReadWriterSize) error {
				return fmt.Errorf("should not be here")
			},
			func(size int) error {
				return nil
			})
		assert.Error(t, err)
	})

}
