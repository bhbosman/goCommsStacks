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

type OutboundStackHandler struct {
	isClosed          bool
	compressionStream *gomessageblock.ReaderWriter
	compression       *flate.Writer
	compressionMutex  *sync.Mutex
	errorState        error
}

func (self *OutboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *OutboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func NewOutboundStackHandler() (*OutboundStackHandler, error) {
	compressionStream := gomessageblock.NewReaderWriter()
	compression, err := flate.NewWriter(compressionStream, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	compressionMutex := &sync.Mutex{}
	return &OutboundStackHandler{
		compressionStream: compressionStream,
		compression:       compression,
		compressionMutex:  compressionMutex,
	}, nil
}

func (self *OutboundStackHandler) MapReadWriterSize(ctx context.Context, size goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	self.compressionMutex.Lock()
	defer self.compressionMutex.Unlock()
	if self.isClosed {
		return nil, goerrors.InvalidState
	}
	uncompressedSize, err := io.Copy(self.compression, size)
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, err
	}
	err = self.compression.Flush()
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, err
	}
	b := [8]byte{}
	binary.LittleEndian.PutUint64(b[:], uint64(uncompressedSize))

	if ctx.Err() != nil {
		return nil, err
	}
	_, err = size.Write(b[:])
	if err != nil {
		return nil, err
	}

	if ctx.Err() != nil {
		return nil, err
	}
	_, err = io.Copy(size, self.compressionStream)
	if err != nil {
		return nil, err
	}

	return size, nil
}

func (self *OutboundStackHandler) Destroy() error {
	self.compressionMutex.Lock()
	defer self.compressionMutex.Unlock()
	return self.close()
}

func (self *OutboundStackHandler) Start(ctx context.Context) error {
	self.compressionMutex.Lock()
	defer self.compressionMutex.Unlock()
	return ctx.Err()
}

func (self *OutboundStackHandler) End() error {
	self.compressionMutex.Lock()
	defer self.compressionMutex.Unlock()
	return self.close()
}

func (self *OutboundStackHandler) close() error {
	if !self.isClosed {
		self.isClosed = true
		var err error
		err = multierr.Append(err, self.compression.Close())
		return err
	}
	return nil
}
