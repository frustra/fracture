package network

import (
	"io"
)

// Type proxyByteReader wraps an io.Reader with the ByteReader interface
type proxyByteReader struct {
	io.Reader
}

func (r *proxyByteReader) ReadByte() (c byte, err error) {
	buf := make([]byte, 1)
	if _, err := r.Read(buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}
