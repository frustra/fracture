package protocol

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
)

func serializeValueTo(w *bytes.Buffer, val interface{}) (int, error) {
	switch i := val.(type) {
	case []byte:
		return w.Write(i)

	case byte:
		return 1, w.WriteByte(i)

	case bool:
		if i {
			return 1, w.WriteByte(1)
		}
		return 1, w.WriteByte(0)

	case int16:
		return w.Write([]byte{
			byte(i >> 8),
			byte(i),
		})

	case uint16:
		return w.Write([]byte{
			byte(i >> 8),
			byte(i),
		})

	case int32:
		return w.Write([]byte{
			byte(i >> 24),
			byte(i >> 16),
			byte(i >> 8),
			byte(i),
		})

	case uint32:
		return w.Write([]byte{
			byte(i >> 24),
			byte(i >> 16),
			byte(i >> 8),
			byte(i),
		})

	case int64:
		return w.Write([]byte{
			byte(i >> 56),
			byte(i >> 48),
			byte(i >> 40),
			byte(i >> 32),
			byte(i >> 24),
			byte(i >> 16),
			byte(i >> 8),
			byte(i),
		})

	case float32:
		k := math.Float32bits(i)
		return w.Write([]byte{
			byte(k >> 24),
			byte(k >> 16),
			byte(k >> 8),
			byte(k),
		})

	case float64:
		k := math.Float64bits(i)
		return w.Write([]byte{
			byte(k >> 56),
			byte(k >> 48),
			byte(k >> 40),
			byte(k >> 32),
			byte(k >> 24),
			byte(k >> 16),
			byte(k >> 8),
			byte(k),
		})

	case string:
		nlen, err := w.Write(Uvarint{uint64(len(i))}.Bytes())
		if err != nil {
			return 0, err
		}

		n, err := w.WriteString(i)
		return nlen + n, err

	case Varint:
		return w.Write(i.Bytes())

	case Uvarint:
		return w.Write(i.Bytes())

	case Serializable:
		return w.Write(i.Serialize())

	default:
		return 0, fmt.Errorf("unknown serialization: %s\n", reflect.ValueOf(val).String())
	}
}
