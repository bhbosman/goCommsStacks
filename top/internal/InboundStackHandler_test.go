package internal

import (
	"context"
	"github.com/bhbosman/gomessageblock"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestInboundStackHandler_MapReadWriterSizeWithDataAndNoError(t *testing.T) {
	stackHandler := NewInboundStackHandler()
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)
	result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
	assert.NoError(t, err)
	sb := strings.Builder{}
	_, _ = io.Copy(&sb, result)
	assert.Equal(t, s, sb.String())
}

func TestInboundStackHandler_MapReadWriterSizeWithDataAndError(t *testing.T) {
	stackHandler := NewInboundStackHandler()
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	result, err := stackHandler.MapReadWriterSize(cancel, rws)
	assert.Nil(t, result)
	assert.Error(t, err)
}
