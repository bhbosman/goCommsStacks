package messageCompressor

import (
	"compress/flate"
	"context"
	"encoding/binary"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"io"
	"sync"
)

type outboundStackHandler struct {
	isClosed          bool
	compressionStream *gomessageblock.ReaderWriter
	compression       *flate.Writer
	compressionMutex  *sync.Mutex
	errorState        error
}

func (self *outboundStackHandler) FlatMapHandler(
	_ context.Context,
	_ interface{},
) (RxHandlers.FlatMapHandlerResult, error) {
	return RxHandlers.NewFlatMapHandlerResult(
		true,
		nil,
		0,
		0,
		0,
		0,
	), nil
}

func (self *outboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *outboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func NewOutboundStackHandler() (RxHandlers.IRxMapStackHandler, error) {
	compressionStream := gomessageblock.NewReaderWriter()
	compression, err := flate.NewWriter(compressionStream, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	compressionMutex := &sync.Mutex{}
	return &outboundStackHandler{
		compressionStream: compressionStream,
		compression:       compression,
		compressionMutex:  compressionMutex,
	}, nil
}

func (self *outboundStackHandler) MapReadWriterSize(
	ctx context.Context,
	unk interface{},
) (interface{}, error) {
	if self.errorState != nil {
		return nil, self.errorState
	}
	if ctx.Err() != nil {
		self.errorState = ctx.Err()
		return nil, ctx.Err()
	}
	switch rws := unk.(type) {
	case goprotoextra.IReadWriterSize:
		self.compressionMutex.Lock()
		defer self.compressionMutex.Unlock()
		if self.isClosed {
			return nil, goerrors.InvalidState
		}
		uncompressedSize, err := io.Copy(self.compression, rws)
		if err != nil {
			return nil, err
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		err = self.compression.Flush()
		if err != nil {
			return nil, err
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		b := [8]byte{}
		binary.LittleEndian.PutUint64(b[:], uint64(uncompressedSize))

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		_, err = rws.Write(b[:])
		if err != nil {
			return nil, err
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		_, err = io.Copy(rws, self.compressionStream)
		if err != nil {
			return nil, err
		}

		return rws, nil
	default:
		return rws, nil
	}
}

func (self *outboundStackHandler) Destroy() error {
	self.compressionMutex.Lock()
	defer self.compressionMutex.Unlock()
	return self.close()
}

func (self *outboundStackHandler) Start(ctx context.Context) error {
	self.compressionMutex.Lock()
	defer self.compressionMutex.Unlock()
	return ctx.Err()
}

func (self *outboundStackHandler) End() error {
	self.compressionMutex.Lock()
	defer self.compressionMutex.Unlock()
	return self.close()
}

func (self *outboundStackHandler) close() error {
	if !self.isClosed {
		self.isClosed = true
		var err error
		err = multierr.Append(err, self.compression.Close())
		return err
	}
	return nil
}
