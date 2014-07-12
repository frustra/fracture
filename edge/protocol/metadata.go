package protocol

import (
	"bytes"
)

type MetadataID struct {
	index byte
	mask  int
}

var (
	OnFireID    = MetadataID{0, 0x01}
	CrouchedID  = MetadataID{0, 0x02}
	SprintingID = MetadataID{0, 0x08}
	BlockingID  = MetadataID{0, 0x10}
	InvisibleID = MetadataID{0, 0x20}
	AirID       = MetadataID{1, 0}

	HealthID        = MetadataID{6, 0}
	PotionColorID   = MetadataID{7, 0}
	AmbientPotionID = MetadataID{8, 0}
	ArrowCountID    = MetadataID{9, 0}
	NameTagID       = MetadataID{10, 0}
	ShowNameTagID   = MetadataID{11, 0}
)

type Metadata struct {
	values map[byte]interface{}
}

func NewMetadata(v ...interface{}) *Metadata {
	m := &Metadata{
		values: make(map[byte]interface{}),
	}

	var key MetadataID
	var hasKey bool

	for _, val := range v {
		if !hasKey {
			key = val.(MetadataID)
			hasKey = true
		} else {
			m.Set(key, val)
			hasKey = false
		}
	}

	if hasKey {
		panic("wrong number of arguments")
	}
	return m
}

func (m *Metadata) Serialize() []byte {
	buf := new(bytes.Buffer)

	for key, val := range m.values {
		var typeId byte

		switch val.(type) {
		case byte:
			typeId = 0
		case int16:
			typeId = 1
		case uint16:
			typeId = 1
		case int32:
			typeId = 2
		case uint32:
			typeId = 2
		case float32:
			typeId = 3
		case string:
			typeId = 4
		}

		buf.WriteByte(key | typeId<<5)
		serializeValueTo(buf, val)
	}

	buf.WriteByte(127)
	return buf.Bytes()
}

func (m *Metadata) Set(key MetadataID, val interface{}) {
	if key.mask != 0 {
		old, existed := m.values[key.index]

		switch v := val.(type) {
		case bool:
			if !existed {
				old = byte(0)
			}
			if v {
				val = old.(byte) | byte(key.mask)
			} else {
				val = old.(byte) & ^byte(key.mask)
			}
		default:
			panic("unimplemented")
		}
	}
	m.values[key.index] = val
}
