package messageNumber_test

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/goCommsStacks/internal/messageNumber"
	"github.com/bhbosman/gomessageblock"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestOutboundStackHandler_MapReadWriterSizeWithDataAndNoError(t *testing.T) {
	stackHandler, _ := messageNumber.NewOutboundStackHandler()
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)
	result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
	assert.NoError(t, err)
	var p8 [8]byte
	n, err := result.(io.Reader).Read(p8[:])
	assert.NoError(t, err)
	assert.Equal(t, 8, n)
	v := binary.LittleEndian.Uint64(p8[:])
	assert.Equal(t, uint64(1), v)

	sb := strings.Builder{}
	_, _ = io.Copy(&sb, result.(io.Reader))
	assert.Equal(t, s, sb.String())
}

func TestOutboundStackHandler_MapReadWriterSizeWithDataAndError(t *testing.T) {
	stackHandler, _ := messageNumber.NewOutboundStackHandler()
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	result, err := stackHandler.MapReadWriterSize(cancel, rws)
	assert.Nil(t, result)
	assert.Error(t, err)
}
