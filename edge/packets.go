package edge

var packId = map[string]byte{
	"kick": 0xFF,
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
	for i := 0; i < len(p.Fields); i++ {
		switch i := p.Fields[i].(type) {
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
	target[start] = (byte)((len(str) >> 8) & 0xFF)
	target[start+1] = (byte)((len(str)) & 0xFF)
	for i := 0; i < len(str); i++ {
		target[start+i+2] = str[i]
	}
	return start + len(str) + 2
}
