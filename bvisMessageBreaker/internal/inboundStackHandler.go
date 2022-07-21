package internal

import (
	"bytes"
	"encoding/binary"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
)

type OnStateCallback func(onNext func(data []byte) error) (canContinue bool, err error)

type InboundStackHandler struct {
	Rw           *gomessageblock.ReaderWriter
	errorState   error
	Marker       [4]byte
	currentState OnStateCallback
}

func (self *InboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func NewInboundStackHandler(
	marker [4]byte,
) RxHandlers.IRxNextStackHandler {
	result := &InboundStackHandler{
		Rw:           gomessageblock.NewReaderWriter(),
		errorState:   nil,
		Marker:       marker,
		currentState: nil,
	}
	result.currentState = result.OnReadHeader()
	return result
}

func (self *InboundStackHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *InboundStackHandler) ReadMessage(_ interface{}) (interface{}, bool, error) {
	return nil, false, nil
}

func (self *InboundStackHandler) Close() error {
	return nil
}

func (self *InboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *InboundStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	onNext func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error,
) error {
	if self.errorState != nil {
		return self.errorState
	}
	err := self.Rw.SetNext(rws)
	if err != nil {
		return err
	}
	return self.inboundState(
		func(dataBlock []byte) error {
			return onNext(gomessageblock.NewReaderWriterBlock(dataBlock))
		},
	)
}

func (self *InboundStackHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *InboundStackHandler) OnReadHeader() OnStateCallback {
	return func(onNext func(data []byte) error) (bool, error) {
		if self.Rw.Size() >= 4 {
			var p [4]byte
			_, err := self.Rw.Read(p[:])
			if err != nil {
				self.errorState = err
				return false, self.errorState
			}
			c := bytes.Compare(p[:], self.Marker[:])
			if c != 0 {
				self.errorState = goerrors.InvalidSignature
				return false, self.errorState
			}
			self.currentState = self.OnReadLength()
			return true, nil
		}
		return false, nil
	}
}

func (self *InboundStackHandler) OnReadLength() OnStateCallback {
	return func(onNext func(data []byte) error) (bool, error) {
		size := self.Rw.Size()
		if size >= 4 {
			var p [4]byte
			_, err := self.Rw.Read(p[:])
			if err != nil {
				self.errorState = err
				return false, self.errorState
			}
			var length = binary.LittleEndian.Uint32(p[:])
			self.currentState = self.OnReadData(length)
			return true, nil
		}
		return false, nil
	}
}

func (self *InboundStackHandler) OnReadData(length uint32) OnStateCallback {
	return func(onNext func(data []byte) error) (bool, error) {
		if uint32(self.Rw.Size()) >= length {
			dataBlock := make([]byte, length)
			_, err := self.Rw.Read(dataBlock)
			if err != nil {
				self.errorState = err
				return false, self.errorState
			}
			err = onNext(dataBlock)
			if err != nil {
				return false, err
			}
			self.currentState = self.OnReadHeader()
			return true, nil
		}
		return false, nil
	}
}

func (self *InboundStackHandler) inboundState(onNext func(data []byte) error) error {
	var err error
	canContinue := true
	for canContinue {
		canContinue, err = self.currentState(onNext)
		if err != nil {
			return err
		}
	}
	return nil
}
