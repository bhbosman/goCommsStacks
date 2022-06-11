package internal

import (
	"context"
	"github.com/bhbosman/gomessageblock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"strings"
	"testing"
)

func TestOutboundStackHandler_MapReadWriterSizeWithDataAndNoError(t *testing.T) {
	stackHandler, err := NewOutboundStackHandler(zap.NewNop())
	assert.NoError(t, err)
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)
	result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
	assert.NoError(t, err)
	sb := strings.Builder{}
	_, _ = io.Copy(&sb, result)
	assert.Equal(t, s, sb.String())
}

func TestOutboundStackHandler_MapReadWriterSizeWithDataAndError(t *testing.T) {
	stackHandler, err := NewOutboundStackHandler(zap.NewNop())
	assert.NoError(t, err)
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	result, err := stackHandler.MapReadWriterSize(cancel, rws)
	assert.Nil(t, result)
	assert.Error(t, err)
}
