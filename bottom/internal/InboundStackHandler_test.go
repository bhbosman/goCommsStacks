package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/gomessageblock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"strings"
	"testing"
)

func TestInboundStackHandler_MapReadWriterSizeWithDataAndNoError(t *testing.T) {
	stackHandler, err := NewInboundStackHandler(zap.NewNop())
	assert.NoError(t, err)
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)
	result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
	assert.NoError(t, err)
	sb := strings.Builder{}
	_, _ = io.Copy(&sb, result)
	assert.Equal(t, s, sb.String())
}

func TestInboundStackHandler_MapReadWriterSizeWithDataAndError(t *testing.T) {
	stackHandler, err := NewInboundStackHandler(zap.NewNop())
	assert.NoError(t, err)

	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	result, err := stackHandler.MapReadWriterSize(cancel, rws)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestInboundStackHandler_ErrorState(t *testing.T) {
	stackHandler, err := NewInboundStackHandler(zap.NewNop())
	assert.NoError(t, err)
	stackHandler.errorState = fmt.Errorf("some error")
	assert.Error(t, stackHandler.ErrorState())
	assert.Error(t, stackHandler.ReadMessage(&struct{}{}))
	_, err = stackHandler.MapReadWriterSize(context.Background(), nil)
	assert.Error(t, err)
}

func TestInboundStackHandler_ReadMessage(t *testing.T) {
	stackHandler, err := NewInboundStackHandler(zap.NewNop())
	assert.NoError(t, err)
	assert.NoError(t, stackHandler.ReadMessage(struct{}{}))
}

func TestInboundNoLogger(t *testing.T) {
	_, err := NewInboundStackHandler(nil)
	assert.Error(t, err)
}
