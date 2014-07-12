package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

const (
	KeepAliveID             = 0x00
	JoinGameID              = 0x01
	ChatMessageID           = 0x02
	SpawnPositionID         = 0x05
	PlayerPositionAndLookID = 0x08
	SpawnPlayerID           = 0x0C
	EntityTeleportID        = 0x18
	MapChunkBulkID          = 0x26
	PlayerListItemID        = 0x38
	PlayerAbilitiesID       = 0x39
	DisconnectID            = 0x40

	PreAuthKickID       = 0x00
	StatusResponseID    = 0x00
	PingResponseID      = 0x01
	EncryptionRequestID = 0x01
	LoginSuccessID      = 0x02
)

type Packet struct {
	Id     uint64
	Fields []interface{}
}

type RawPacket struct {
	Id  uint64
	Buf []byte
}

type Varint struct {
	Val int64
}

type Uvarint struct {
	Val uint64
}

type Serializable interface {
	Serialize() []byte
}

func (v Varint) Bytes() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, v.Val)
	return buf[0:n]
}

func (v Uvarint) Bytes() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, v.Val)
	return buf[0:n]
}

func CreatePacket(id uint64, v ...interface{}) *Packet {
	return &Packet{id, v}
}

func (p RawPacket) Serialize() []byte {
	id := Uvarint{p.Id}.Bytes()
	size := Uvarint{uint64(len(p.Buf) + len(id))}.Bytes()
	buf := make([]byte, len(p.Buf)+len(id)+len(size))
	copy(buf, size)
	copy(buf[len(size):], id)
	copy(buf[len(id)+len(size):], p.Buf)
	return buf
}

func (p Packet) Serialize() []byte {
	buf := make([]byte, 65536)
	n := binary.PutUvarint(buf, p.Id)
	for j := 0; j < len(p.Fields); j++ {
		switch i := p.Fields[j].(type) {
		case []byte:
			for k := 0; k < len(i); k++ {
				buf[n] = i[k]
				n++
			}
		case byte:
			buf[n] = i
			n++
		case bool:
			if i {
				buf[n] = 0x01
			} else {
				buf[n] = 0x00
			}
			n++
		case int16:
			buf[n] = byte((i >> 8) & 0xFF)
			buf[n+1] = byte(i & 0xFF)
			n += 2
		case uint16:
			buf[n] = byte((i >> 8) & 0xFF)
			buf[n+1] = byte(i & 0xFF)
			n += 2
		case int32:
			buf[n] = byte((i >> 24) & 0xFF)
			buf[n+1] = byte((i >> 16) & 0xFF)
			buf[n+2] = byte((i >> 8) & 0xFF)
			buf[n+3] = byte(i & 0xFF)
			n += 4
		case uint32:
			buf[n] = byte((i >> 24) & 0xFF)
			buf[n+1] = byte((i >> 16) & 0xFF)
			buf[n+2] = byte((i >> 8) & 0xFF)
			buf[n+3] = byte(i & 0xFF)
			n += 4
		case int64:
			buf[n] = byte((i >> 56) & 0xFF)
			buf[n+1] = byte((i >> 48) & 0xFF)
			buf[n+2] = byte((i >> 40) & 0xFF)
			buf[n+3] = byte((i >> 32) & 0xFF)
			buf[n+4] = byte((i >> 24) & 0xFF)
			buf[n+5] = byte((i >> 16) & 0xFF)
			buf[n+6] = byte((i >> 8) & 0xFF)
			buf[n+7] = byte(i & 0xFF)
			n += 8
		case float32:
			k := math.Float32bits(i)
			buf[n] = byte((k >> 24) & 0xFF)
			buf[n+1] = byte((k >> 16) & 0xFF)
			buf[n+2] = byte((k >> 8) & 0xFF)
			buf[n+3] = byte(k & 0xFF)
			n += 4
		case float64:
			k := math.Float64bits(i)
			buf[n] = byte((k >> 56) & 0xFF)
			buf[n+1] = byte((k >> 48) & 0xFF)
			buf[n+2] = byte((k >> 40) & 0xFF)
			buf[n+3] = byte((k >> 32) & 0xFF)
			buf[n+4] = byte((k >> 24) & 0xFF)
			buf[n+5] = byte((k >> 16) & 0xFF)
			buf[n+6] = byte((k >> 8) & 0xFF)
			buf[n+7] = byte(k & 0xFF)
			n += 8
		case string:
			n += binary.PutUvarint(buf[n:], uint64(len(i)))
			for k := 0; k < len(i); k++ {
				buf[n] = i[k]
				n++
			}
		case Varint:
			buf2 := i.Bytes()
			for k := 0; k < len(buf2); k++ {
				buf[n] = buf2[k]
				n++
			}
		case Uvarint:
			buf2 := i.Bytes()
			for k := 0; k < len(buf2); k++ {
				buf[n] = buf2[k]
				n++
			}
		case Serializable:
			buf2 := i.Serialize()
			for k := 0; k < len(buf2); k++ {
				buf[n] = buf2[k]
				n++
			}
		default:
			fmt.Printf("Unknown serialization: (%d) %s\n", j, reflect.ValueOf(p.Fields[j]).String())
		}
	}
	return append(Uvarint{uint64(n)}.Bytes(), buf[0:n]...)
}

func (p Packet) Write(conn io.Writer) error {
	return WritePacket(conn, p)
}

func WriteNewPacket(conn io.Writer, id uint64, v ...interface{}) error {
	return Packet{id, v}.Write(conn)
}

func WritePacket(conn io.Writer, p Serializable) error {
	buf := p.Serialize()
	// fmt.Println(hex.Dump(buf))
	n := 0
	for n < len(buf) {
		n2, err := conn.Write(buf[n:])
		if err != nil {
			return err
		} else {
			n += n2
		}
	}
	return nil
}

func ReadPacket(conn io.Reader) (uint64, []byte, error) {
	buf := make([]byte, 65536)
	size := make([]byte, 256)
	tmp := make([]byte, 1)
	sizen := 0
	n := 0
	length := -1
	for {
		read, err := conn.Read(tmp)
		if read > 0 && err == nil {
			buf[n] = tmp[0]
			n++
			if length >= 0 {
				if n >= length {
					break
				}
			} else if (tmp[0] & 0x80) == 0 {
				len2, _ := binary.Uvarint(buf[0:n])
				length = int(len2)
				copy(size, buf[:n])
				sizen = n
				n = 0
			}
		} else if err != nil {
			return 0, append(size[:sizen], buf[:n]...), err
		}
	}
	id, n2 := binary.Uvarint(buf)
	return id, buf[n2:n], nil
}

func ReadString(buf []byte, start int) (string, int) {
	if start < 0 || start >= len(buf) {
		return "", -1
	}
	size, n := binary.Uvarint(buf[start:])
	if n <= 0 || start+int(size)+n > len(buf) {
		return "", -1
	}
	return string(buf[start+n : start+int(size)+n]), start + int(size) + n
}

func ReadByte(buf []byte, start int) (int, int) {
	if start < 0 || start+1 > len(buf) {
		return 0, -1
	}
	return int(buf[start]), start + 1
}

func ReadShort(buf []byte, start int) (int, int) {
	if start < 0 || start+2 > len(buf) {
		return 0, -1
	}
	val := int(buf[start]) << 8
	val |= int(buf[start+1])
	return val, start + 2
}

func ReadInt(buf []byte, start int) (int, int) {
	if start < 0 || start+4 > len(buf) {
		return 0, -1
	}
	val := int(buf[start]) << 24
	val |= int(buf[start+1]) << 16
	val |= int(buf[start+2]) << 8
	val |= int(buf[start+3])
	return val, start + 4
}

func ReadLong(buf []byte, start int) (int64, int) {
	if start < 0 || start+8 > len(buf) {
		return 0, -1
	}
	val := int64(buf[start]) << 56
	val |= int64(buf[start+1]) << 48
	val |= int64(buf[start+2]) << 40
	val |= int64(buf[start+3]) << 32
	val |= int64(buf[start+4]) << 24
	val |= int64(buf[start+5]) << 16
	val |= int64(buf[start+6]) << 8
	val |= int64(buf[start+7])
	return val, start + 8
}

func ReadVarint(buf []byte, start int) (int64, int) {
	if start < 0 || start >= len(buf) {
		return 0, -1
	}
	result, n := binary.Varint(buf[start:])
	if n <= 0 || start+n > len(buf) {
		return 0, -1
	}
	return result, start + n
}

func ReadUvarint(buf []byte, start int) (uint64, int) {
	if start < 0 || start >= len(buf) {
		return 0, -1
	}
	result, n := binary.Uvarint(buf[start:])
	if n <= 0 || start+n > len(buf) {
		return 0, -1
	}
	return result, start + n
}

func ReadBytes(buf []byte, start int, length int) ([]byte, int) {
	if start < 0 || start+length > len(buf) {
		return nil, -1
	}
	return buf[start : start+length], start + length
}
