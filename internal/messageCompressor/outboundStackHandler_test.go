package messageCompressor

import (
	"compress/flate"
	"context"
	"encoding/binary"
	"github.com/bhbosman/gomessageblock"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestOutboundStackHandler_MapReadWriterSizeWithDataAndNoError(t *testing.T) {
	stackHandler, err := NewOutboundStackHandler()
	assert.NoError(t, err)
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)
	result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
	assert.NoError(t, err)

	b := [8]byte{}
	_, err = result.Read(b[:])
	assert.NoError(t, err)

	decompressorStream := gomessageblock.NewReaderWriter()
	decompressor := flate.NewReader(decompressorStream)

	uncompressedLength := int64(binary.LittleEndian.Uint64(b[:]))
	_, err = io.Copy(decompressorStream, result)

	_, err = io.CopyN(result, decompressor, uncompressedLength)

	sb := strings.Builder{}
	_, _ = io.Copy(&sb, result)
	assert.Equal(t, s, sb.String())
}

func TestOutboundStackHandler_MapReadWriterSizeWithDataAndError(t *testing.T) {
	stackHandler, err := NewOutboundStackHandler()
	assert.NoError(t, err)
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)

	cancel, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	result, err := stackHandler.MapReadWriterSize(cancel, rws)
	assert.Nil(t, result)
	assert.Error(t, err)
}
