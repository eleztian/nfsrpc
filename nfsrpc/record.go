package nfsrpc

import (
	"io"
	"encoding/binary"
	"bytes"
)

const (
	// This is maximum size in bytes for an individual record fragment.
	// The entire RPC message (record) has no size restriction imposed
	// by RFC 5531. Refer: include/linux/sunrpc/msg_prot.h
	maxRecordFragmentSize = (1 << 31) - 1

	// Max size of RPC message that a client is allowed to send.
	maxRecordSize = 1 * 1024 * 1024
)

func Write(conn io.Writer, data []byte) (int64, error) {
	dataSize := int64(len(data))

	var totalBytesWritten int64
	var lastFragment bool

	fragmentHeaderBytes := make([]byte, 4)
	for {
		remainingBytes := dataSize - totalBytesWritten
		if remainingBytes <= maxRecordFragmentSize {
			lastFragment = true
		}
		fragmentSize := uint32(minOf(maxRecordFragmentSize, remainingBytes))

		// Create fragment header
		binary.BigEndian.PutUint32(fragmentHeaderBytes, createFragmentHeader(fragmentSize, lastFragment))

		// Write fragment header and fragment body to network
		bytesWritten, err := conn.Write(append(fragmentHeaderBytes, data[totalBytesWritten:fragmentSize]...))
		if err != nil {
			return totalBytesWritten, err
		}
		totalBytesWritten += int64(bytesWritten)

		if lastFragment {
			break
		}

	}
	return totalBytesWritten, nil
}

func minOf(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func createFragmentHeader(size uint32, lastFragment bool) uint32 {
	fragmentHeader := size &^ (1 << 31)

	if lastFragment {
		fragmentHeader |= (1 << 31)
	}

	return fragmentHeader
}

func Read(conn io.Reader) ([]byte, error) {
	readBuf := bytes.NewBuffer(make([]byte, 0, maxRecordSize))
	var fragmentHeader uint32
	for {
		if err := binary.Read(conn, binary.BigEndian, &fragmentHeader); err != nil {
			return nil, err
		}

		fragmentSize := getFragmentSize(fragmentHeader)
		if fragmentSize > maxRecordFragmentSize {
			return nil, ErrInvalidFragmentSize
		}
		// 剩余容量不容存储，溢出
		if int(fragmentSize) > (readBuf.Cap() - readBuf.Len()) {
			return nil, ErrRPCMessageSizeExceeded
		}
		bytesCopied, err := io.CopyN(readBuf, conn, int64(fragmentSize))
		if err != nil || (bytesCopied != int64(fragmentSize)) {
			return nil, err
		}
		if isLastFragment(fragmentHeader) {
			break
		}
	}
	return readBuf.Bytes(), nil
}

func getFragmentSize(fragmentHeader uint32) uint32 {
	return fragmentHeader &^ (1 << 31)
}
func isLastFragment(fragmentHeader uint32) bool {
	return (fragmentHeader >> 31) == 1
}