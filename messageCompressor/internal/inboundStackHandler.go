package internal

import (
	"compress/flate"
	"context"
	"encoding/binary"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"io"
	"sync"
)

type InboundStackHandler struct {
	isClosed           bool
	decompressorStream *gomessageblock.ReaderWriter
	decompressor       io.ReadCloser
	decompressorMutex  *sync.Mutex
	errorState         error
}

func (self *InboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *InboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func (self *InboundStackHandler) close() error {
	if !self.isClosed {
		self.isClosed = true
		var err error
		err = multierr.Append(err, self.decompressor.Close())
		return err
	}
	return nil
}

func NewInboundStackHandler() (*InboundStackHandler, error) {
	decompressorStream := gomessageblock.NewReaderWriter()
	decompressor := flate.NewReader(decompressorStream)
	decompressorMutex := &sync.Mutex{}

	return &InboundStackHandler{
		decompressorStream: decompressorStream,
		decompressor:       decompressor,
		decompressorMutex:  decompressorMutex,
	}, nil
}

func (self *InboundStackHandler) MapReadWriterSize(ctx context.Context, size goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	self.decompressorMutex.Lock()
	defer self.decompressorMutex.Unlock()
	if self.isClosed {
		return nil, goerrors.InvalidState
	}
	b := [8]byte{}
	_, err := size.Read(b[:])
	if err != nil {
		self.errorState = err
		return nil, self.errorState
	}
	uncompressedLength := int64(binary.LittleEndian.Uint64(b[:]))
	_, err = io.Copy(self.decompressorStream, size)
	if err != nil {
		self.errorState = err
		return nil, self.errorState
	}

	_, err = io.CopyN(size, self.decompressor, uncompressedLength)
	if err != nil {
		self.errorState = err
		return nil, self.errorState
	}

	return size, nil
}

func (self *InboundStackHandler) Destroy() error {
	self.decompressorMutex.Lock()
	defer self.decompressorMutex.Unlock()
	return self.close()
}

func (self *InboundStackHandler) Start(ctx context.Context) error {
	self.decompressorMutex.Lock()
	defer self.decompressorMutex.Unlock()
	return ctx.Err()
}

func (self *InboundStackHandler) End() error {
	self.decompressorMutex.Lock()
	defer self.decompressorMutex.Unlock()
	return self.close()
}
