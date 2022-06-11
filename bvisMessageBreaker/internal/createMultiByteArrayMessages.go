package internal

import (
	"encoding/binary"
	"github.com/bhbosman/gomessageblock"
)

func CreateMultiByteArrayMessages(marker [4]byte, rws *gomessageblock.ReaderWriter) func(s []byte) error {
	return func(s []byte) error {
		block := [4]byte{}
		binary.LittleEndian.PutUint32(block[0:4], uint32(len(s)))
		err := rws.AddReaders(
			gomessageblock.NewReaderWriterBlock(marker[:]),
			gomessageblock.NewReaderWriterBlock(block[:]),
			gomessageblock.NewReaderWriterBlock(s))
		if err != nil {
			return err
		}
		return nil
	}
}
