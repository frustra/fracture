package protocol_test

import (
	"testing"

	. "github.com/frustra/fracture/edge/protocol"
)

var MetadataSerializationTests = map[*Metadata][]byte{
	NewMetadata(ArrowCountID, byte(0x52)): []byte{0x09, 0x52, 0x7f},
	NewMetadata(HealthID, float32(10.4)):  []byte{3<<5 | 6, 0x41, 0x26, 0x66, 0x66, 0x7f},

	NewMetadata(
		AbsorptionHeartsID, float32(0),
		OnFireID, false,
		UnknownBitFieldID, byte(0),
		AirID, uint16(0x012c),
		ScoreID, uint32(0),
		HealthID, float32(20),
		PotionColorID, int32(0),
		AmbientPotionID, byte(0),
		ArrowCountID, byte(0),
	): []byte{
		0x71, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x21, 0x01, 0x2c,
		0x52, 0x00, 0x00, 0x00, 0x00, 0x66, 0x41, 0xa0, 0x00, 0x00,
		0x47, 0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x09, 0x00, 0x7f,
	},
}

func TestMetadataSerialization(t *testing.T) {
	for m, expected := range MetadataSerializationTests {
		got := m.Serialize()

		if string(got) != string(expected) {
			t.Errorf("Incorrectly serialized metadata.\n       got: %#v\n  expected: %#v", got, expected)
		}
	}
}
