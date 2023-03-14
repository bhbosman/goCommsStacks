package messageCompressor

import (
	"compress/flate"
	"context"
	"encoding/binary"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"io"
	"sync"
)

type inboundStackHandler struct {
	isClosed           bool
	decompressorStream *gomessageblock.ReaderWriter
	decompressor       io.ReadCloser
	decompressorMutex  *sync.Mutex
	errorState         error
}

func (self *inboundStackHandler) PublishCounters(counters *model.PublishRxHandlerCounters) {
}

func (self *inboundStackHandler) EmptyQueue() {
}

func (self *inboundStackHandler) ClearCounters() {
}

func (self *inboundStackHandler) FlatMapHandler(_ context.Context, _ interface{}) (RxHandlers.FlatMapHandlerResult, error) {
	return RxHandlers.NewFlatMapHandlerResult(true, nil, 0, 0, 0, 0), nil
}

func (self *inboundStackHandler) ErrorState() error {
	return self.errorState
}

func (self *inboundStackHandler) close() error {
	if !self.isClosed {
		self.isClosed = true
		var err error
		err = multierr.Append(err, self.decompressor.Close())
		return err
	}
	return nil
}

func NewInboundStackHandler() (RxHandlers.IRxMapStackHandler, error) {
	decompressorStream := gomessageblock.NewReaderWriter()
	decompressor := flate.NewReader(decompressorStream)
	decompressorMutex := &sync.Mutex{}

	return &inboundStackHandler{
		decompressorStream: decompressorStream,
		decompressor:       decompressor,
		decompressorMutex:  decompressorMutex,
	}, nil
}

func (self *inboundStackHandler) MapReadWriterSize(
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
		self.decompressorMutex.Lock()
		defer self.decompressorMutex.Unlock()
		if self.isClosed {
			return nil, goerrors.InvalidState
		}
		b := [8]byte{}
		_, err := rws.Read(b[:])
		if err != nil {
			self.errorState = err
			return nil, self.errorState
		}
		uncompressedLength := int64(binary.LittleEndian.Uint64(b[:]))
		_, err = io.Copy(self.decompressorStream, rws)
		if err != nil {
			self.errorState = err
			return nil, self.errorState
		}

		_, err = io.CopyN(rws, self.decompressor, uncompressedLength)
		if err != nil {
			self.errorState = err
			return nil, self.errorState
		}

		return rws, nil
	default:
		return rws, nil
	}
}

func (self *inboundStackHandler) Destroy() error {
	self.decompressorMutex.Lock()
	defer self.decompressorMutex.Unlock()
	return self.close()
}

func (self *inboundStackHandler) Start(ctx context.Context) error {
	self.decompressorMutex.Lock()
	defer self.decompressorMutex.Unlock()
	return ctx.Err()
}

func (self *inboundStackHandler) End() error {
	self.decompressorMutex.Lock()
	defer self.decompressorMutex.Unlock()
	return self.close()
}
