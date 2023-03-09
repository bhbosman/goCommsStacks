package bvisMessageBreaker

import (
	"bytes"
	"encoding/binary"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
)

type OnStateCallback func(onNext func(data *gomessageblock.ReaderWriter) error) (canContinue bool, err error)

type inboundStackHandler struct {
	rw           *gomessageblock.ReaderWriter
	errorState   error
	marker       [4]byte
	currentState OnStateCallback
}

func (self *inboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func (self *inboundStackHandler) onReadHeader() OnStateCallback {
	return func(onNext func(data *gomessageblock.ReaderWriter) error) (bool, error) {
		if self.rw.Size() >= 4 {
			var p [4]byte
			_, err := self.rw.Read(p[:])
			if err != nil {
				self.errorState = err
				return false, self.errorState
			}
			c := bytes.Compare(p[:], self.marker[:])
			if c != 0 {
				self.errorState = goerrors.InvalidSignature
				return false, self.errorState
			}
			self.currentState = self.onReadLength()
			return true, nil
		}
		return false, nil
	}
}

func (self *inboundStackHandler) onReadLength() OnStateCallback {
	return func(onNext func(data *gomessageblock.ReaderWriter) error) (bool, error) {
		size := self.rw.Size()
		if size >= 4 {
			var p [4]byte
			_, err := self.rw.Read(p[:])
			if err != nil {
				self.errorState = err
				return false, self.errorState
			}
			var length = binary.LittleEndian.Uint32(p[:])
			self.currentState = self.onReadData(length)
			return true, nil
		}
		return false, nil
	}
}

func (self *inboundStackHandler) onReadData(length uint32) OnStateCallback {
	return func(onNext func(data *gomessageblock.ReaderWriter) error) (bool, error) {
		if uint32(self.rw.Size()) >= length {
			dataBlock := make([]byte, length)
			_, err := self.rw.Read(dataBlock)
			if err != nil {
				self.errorState = err
				return false, self.errorState
			}
			err = onNext(gomessageblock.NewReaderWriterBlock(dataBlock))
			if err != nil {
				return false, err
			}
			self.currentState = self.onReadHeader()
			return true, nil
		}
		return false, nil
	}
}

func (self *inboundStackHandler) inboundState(onNext func(data *gomessageblock.ReaderWriter) error) error {
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
