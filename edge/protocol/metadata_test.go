package protocol_test

import (
	"testing"

	"github.com/frustra/fracture/edge/protocol"
)

func TestMetadataSerialization(t *testing.T) {
	m := protocol.NewMetadata()
	m.Set(protocol.ArrowCountID, byte(0x52))

	got := m.Serialize()
	expected := "\x09\x52\x7f"

	if string(got) != expected {
		t.Errorf("Incorrectly serialized metadata.\n     got: %#v\n  expected: %#v", got, []byte(expected))
	}
}
