package edge

import (
	"math"
)

var packId = map[string]byte{
	"keepalive": 0x00,
	"login":     0x01,
	"handshake": 0x02,

	"chatmsg":    0x03,
	"timeupdate": 0x04,
	"entequip":   0x05,
	"spawnpos":   0x06,
	"useent":     0x07,
	"updatehp":   0x08,
	"respawn":    0x09,

	"player":        0x0A,
	"playerpos":     0x0B,
	"playerlook":    0x0C,
	"playerposlook": 0x0D,
	"playerdig":     0x0E,
	"playerplace":   0x0F,

	"heldchange": 0x10,
	"usebed":     0x11,
	"animation":  0x12,
	"entaction":  0x13,

	"spawnent":      0x14,
	"spawnitem":     0x15,
	"collectitem":   0x16,
	"spawnobj":      0x17,
	"spawnmob":      0x18,
	"spawnpainting": 0x19,
	"spawnexp":      0x1A,

	"entvelocity":     0x1C,
	"destroyent":      0x1D,
	"entity":          0x1E,
	"entmove":         0x1F,
	"entlook":         0x20,
	"entlookmove":     0x21,
	"entteleport":     0x22,
	"entheadlook":     0x23,
	"entstatus":       0x26,
	"attachent":       0x27,
	"entmeta":         0x28,
	"enteffect":       0x29,
	"removeenteffect": 0x2A,
	"setexp":          0x2B,

	"chunkdata":        0x33,
	"multiblockchange": 0x34,
	"blockchange":      0x35,
	"blockaction":      0x36,
	"blockanimation":   0x37,
	"mapchunkbulk":     0x38,
	"explosion":        0x3C,
	"soundorparticle":  0x3D,
	"soundeffect":      0x3E,
	"gamestate":        0x46,
	"globalent":        0x47,

	"openwin":           0x64,
	"closewin":          0x65,
	"clickwin":          0x66,
	"setslot":           0x67,
	"setwinitems":       0x68,
	"updatewinproperty": 0x69,
	"confirmtrans":      0x6A,
	"creativeinvaction": 0x6B,
	"enchantitem":       0x6C,
	"updatesign":        0x82,
	"itemdata":          0x83,

	"updatetileent":      0x84,
	"incrementstatistic": 0xCB,
	"playerlistitem":     0xC9,
	"playerabilities":    0xCA,
	"tabcomplete":        0xCB,
	"clientsettings":     0xCC,
	"clientstatuses":     0xCD,
	"pluginmsg":          0xFA,

	"encryptresponse": 0xFC,
	"encryptrequest":  0xFD,
	"listping":        0xFE,
	"kick":            0xFF,
}
var packName = map[byte]string{}

func InitPackets() {
	for k, v := range packId {
		packName[v] = k
	}
}

type Packet struct {
	Id     byte
	Fields []interface{}
}

func CreatePacket(name string, v ...interface{}) Packet {
	id := packId[name]
	return Packet{id, v}
}

func (p Packet) Serialize() []byte {
	buf := make([]byte, 65536)
	n := 1
	buf[0] = p.Id
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
		case int32:
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
			n = toUtf16ByteArray(i, buf, n)
		}
	}
	return buf[0:n]
}

func toUtf16ByteArray(str string, target []byte, start int) int {
	target[start] = (byte)((len(str) >> 8) & 0xFF)
	target[start+1] = (byte)((len(str)) & 0xFF)
	for i := 0; i < len(str); i++ {
		target[start+i*2+2] = 0
		target[start+i*2+3] = str[i]
	}
	return start + len(str)*2 + 2
}

func toUtf8ByteArray(str string, target []byte, start int) int {
	for i := 0; i < len(str); i++ {
		target[start+i] = str[i]
	}
	target[start+len(str)] = 0x00
	return start + len(str) + 1
}

func JoinStrings(v ...interface{}) string {
	buf := make([]byte, 65536)
	n := 0
	for i := 0; i < len(v); i++ {
		switch i := v[i].(type) {
		case []byte:
			for j := 0; j < len(i); j++ {
				buf[n] = i[j]
				n++
			}
			buf[n] = 0x00
			n += 1
		case string:
			n = toUtf8ByteArray(i, buf, n)
		}
	}
	return string(buf[0:n])
}

func ReadString(buf []byte, start int) (string, int) {
	size := int(buf[start]) << 8
	size |= int(buf[start+1])
	return string(buf[start+2 : start+size*2+2]), start + size*2 + 2
}

func ReadShort(buf []byte, start int) (int, int) {
	val := int(buf[start]) << 8
	val |= int(buf[start+1])
	return val, start + 2
}

func ReadInt(buf []byte, start int) (int, int) {
	val := int(buf[start]) << 24
	val |= int(buf[start+1]) << 16
	val |= int(buf[start+2]) << 8
	val |= int(buf[start+3])
	return val, start + 4
}

func ReadBytes(buf []byte, start int, length int) ([]byte, int) {
	return buf[start : start+length], start + length
}
