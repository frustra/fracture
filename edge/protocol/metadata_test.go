package protocol_test

import (
	"testing"

	. "github.com/frustra/fracture/edge/protocol"
)

var MetadataSerializationTests = map[*Metadata][]byte{
	NewMetadata(ArrowCountID, byte(0x52)): []byte{0x09, 0x52, 0x7f},
	NewMetadata(HealthID, float32(10.4)):  []byte{0x66, 0x41, 0x26, 0x66, 0x66, 0x7f},
}

func TestMetadataSerialization(t *testing.T) {
	for m, expected := range MetadataSerializationTests {
		got := m.Serialize()

		if string(got) != string(expected) {
			t.Errorf("Incorrectly serialized metadata.\n       got: %#v\n  expected: %#v", got, expected)
		}
	}
}
