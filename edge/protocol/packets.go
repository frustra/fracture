package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"math"
)

const (
	Version = 5

	KeepAliveID             = 0x00
	JoinGameID              = 0x01
	ChatMessageID           = 0x02
	SpawnPositionID         = 0x05
	PlayerPositionAndLookID = 0x08
	SpawnPlayerID           = 0x0C
	DestroyEntitiesID       = 0x13
	EntityRelativeMoveID    = 0x15
	EntityLookID            = 0x16
	EntityLookAndMoveID     = 0x17
	EntityTeleportID        = 0x18
	EntityHeadLookID        = 0x19
	BlockChangeID           = 0x23
	EntityMetadataID        = 0x1C
	EntityPropertiesID      = 0x20
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
	Val uint64
}

type Serializable interface {
	Serialize() []byte
}

func (v Varint) Bytes() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, v.Val)
	return buf[0:n]
}

func CreatePacket(id uint64, v ...interface{}) *Packet {
	return &Packet{id, v}
}

func (p RawPacket) Serialize() []byte {
	id := Varint{p.Id}.Bytes()
	size := Varint{uint64(len(p.Buf) + len(id))}.Bytes()
	buf := make([]byte, len(p.Buf)+len(id)+len(size))

	copy(buf, size)
	copy(buf[len(size):], id)
	copy(buf[len(id)+len(size):], p.Buf)

	return buf
}

func (p Packet) Serialize() []byte {
	buf := new(bytes.Buffer) // TODO: buffer pool?
	buf.Write(Varint{p.Id}.Bytes())

	for j := 0; j < len(p.Fields); j++ {
		_, err := serializeValueTo(buf, p.Fields[j])
		if err != nil {
			log.Printf("Error serializing field %d: %s", j, err)
		}
	}

	length := Varint{uint64(buf.Len())}
	return append(length.Bytes(), buf.Bytes()...)
}

func (p Packet) Write(conn io.Writer) error {
	return WritePacket(conn, p)
}

func WriteNewPacket(conn io.Writer, id uint64, v ...interface{}) error {
	return Packet{id, v}.Write(conn)
}

func WritePacket(conn io.Writer, p Serializable) error {
	buf := p.Serialize()
	// _, x := ReadVarint(buf, 0)
	// id, _ := ReadVarint(buf, x)
	// if id != 0x00 && id != 0x26 {
	// 	fmt.Printf("%x:\n%s\n", id, hex.Dump(buf[x+1:]))
	// }
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
	val := int16(buf[start]) << 8
	val |= int16(buf[start+1])
	return int(val), start + 2
}

func ReadInt(buf []byte, start int) (int, int) {
	if start < 0 || start+4 > len(buf) {
		return 0, -1
	}
	val := int32(buf[start]) << 24
	val |= int32(buf[start+1]) << 16
	val |= int32(buf[start+2]) << 8
	val |= int32(buf[start+3])
	return int(val), start + 4
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

func ReadFloat(buf []byte, start int) (float32, int) {
	val, _ := ReadInt(buf, start)
	return math.Float32frombits(uint32(val)), start + 4
}

func ReadDouble(buf []byte, start int) (float64, int) {
	val, _ := ReadLong(buf, start)
	return math.Float64frombits(uint64(val)), start + 8
}

func ReadBool(buf []byte, start int) (bool, int) {
	val, _ := ReadByte(buf, start)
	return val != 0, start + 1
}

func ReadVarint(buf []byte, start int) (uint64, int) {
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
