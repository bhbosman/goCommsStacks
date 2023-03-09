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

func TestInboundStackHandler_MapReadWriterSizeWithDataAndNoError(t *testing.T) {
	stackHandler, _ := NewInboundStackHandler()
	s := "HelloWorld"
	rws, _ := gomessageblock.NewReaderWriterString(s)

	compressedData := gomessageblock.NewReaderWriter()
	writer, err := flate.NewWriter(compressedData, flate.DefaultCompression)
	assert.NoError(t, err)
	n, err := io.Copy(writer, rws)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), n)
	err = writer.Flush()
	assert.NoError(t, err)

	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], uint64(n))

	_, err = rws.Write(b[:])
	assert.NoError(t, err)
	_, _ = io.Copy(rws, compressedData)

	result, err := stackHandler.MapReadWriterSize(context.Background(), rws)
	assert.NoError(t, err)
	sb := strings.Builder{}
	_, _ = io.Copy(&sb, result.(io.Reader))
	assert.Equal(t, s, sb.String())
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
