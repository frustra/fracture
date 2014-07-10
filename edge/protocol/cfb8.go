package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"net"
)

type AESConn struct {
	conn    *net.TCPConn
	block   cipher.Block
	readIV  []byte
	writeIV []byte
	scratch []byte
}

func NewAESConn(conn *net.TCPConn, secret []byte) (*AESConn, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}
	readIV := make([]byte, len(secret))
	writeIV := make([]byte, len(secret))
	copy(readIV, secret)
	copy(writeIV, secret)
	return &AESConn{conn, block, readIV, writeIV, make([]byte, block.BlockSize())}, nil
}

func (ac AESConn) Read(buf []byte) (int, error) {
	n, err := ac.conn.Read(buf)
	if err == nil {
		XORKeyStreamRead(ac, buf[:n], buf[:n])
	}
	return n, err
}

func (ac AESConn) Write(buf []byte) (int, error) {
	buf2 := make([]byte, len(buf))
	XORKeyStreamWrite(ac, buf2, buf)
	return ac.conn.Write(buf2)
}

func XORKeyStreamRead(conn AESConn, dst, src []byte) {
	for i := 0; i < len(src); i++ {
		val := src[i]
		copy(conn.scratch, conn.readIV)
		conn.block.Encrypt(conn.readIV, conn.readIV)
		val ^= conn.readIV[0]

		copy(conn.readIV, conn.scratch[1:])
		conn.readIV[15] = src[i]

		dst[i] = val
	}
}

func XORKeyStreamWrite(conn AESConn, dst, src []byte) {
	for i := 0; i < len(src); i++ {
		val := src[i]
		copy(conn.scratch, conn.writeIV)
		conn.block.Encrypt(conn.writeIV, conn.writeIV)
		val ^= conn.writeIV[0]

		copy(conn.writeIV, conn.scratch[1:])
		conn.writeIV[15] = val

		dst[i] = val
	}
}
