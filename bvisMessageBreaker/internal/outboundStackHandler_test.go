package internal

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestOutboundStackHandler_MapReadWriterSizeWithDataAndNoError(t *testing.T) {
	marker := [4]byte{1, 2, 3, 4}
	stackHandler, _ := NewOutboundStackHandler(marker)
	s := "HelloWorld"
	var err error
	var result goprotoextra.ReadWriterSize
	rws, _ := gomessageblock.NewReaderWriterString(s)
	result, err = stackHandler.MapReadWriterSize(context.Background(), rws)
	assert.NoError(t, err)
	var p4 [4]byte
	var n int
	n, err = result.Read(p4[:])
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	c := bytes.Compare(p4[:], marker[:])
	assert.Equal(t, 0, c)

	n, err = result.Read(p4[:])
	assert.NoError(t, err)
	assert.Equal(t, 4, n)

	l := binary.LittleEndian.Uint32(p4[:])
	assert.Equal(t, uint32(10), l)

	sb := strings.Builder{}
	_, _ = io.Copy(&sb, result)
	assert.Equal(t, s, sb.String())
}

func TestOutboundStackHandler_MapReadWriterSizeWithDataAndError(t *testing.T) {
	stackHandler, _ := NewOutboundStackHandler(
		[4]byte{1, 2, 3, 4})
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	result, err := stackHandler.MapReadWriterSize(cancel, rws)
	assert.Nil(t, result)
	assert.Error(t, err)
}
