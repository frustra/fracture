// Code generated by protoc-gen-gogo.
// source: player_action.proto
// DO NOT EDIT!

package protobuf

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// discarding unused import gogoproto "code.google.com/p/gogoprotobuf/gogoproto/gogo.pb"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type PlayerAction_Action int32

const (
	PlayerAction_JOIN          PlayerAction_Action = 0
	PlayerAction_LEAVE         PlayerAction_Action = 1
	PlayerAction_MOVE_RELATIVE PlayerAction_Action = 2
	PlayerAction_MOVE_ABSOLUTE PlayerAction_Action = 3
)

var PlayerAction_Action_name = map[int32]string{
	0: "JOIN",
	1: "LEAVE",
	2: "MOVE_RELATIVE",
	3: "MOVE_ABSOLUTE",
}
var PlayerAction_Action_value = map[string]int32{
	"JOIN":          0,
	"LEAVE":         1,
	"MOVE_RELATIVE": 2,
	"MOVE_ABSOLUTE": 3,
}

func (x PlayerAction_Action) Enum() *PlayerAction_Action {
	p := new(PlayerAction_Action)
	*p = x
	return p
}
func (x PlayerAction_Action) String() string {
	return proto.EnumName(PlayerAction_Action_name, int32(x))
}
func (x *PlayerAction_Action) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(PlayerAction_Action_value, data, "PlayerAction_Action")
	if err != nil {
		return err
	}
	*x = PlayerAction_Action(value)
	return nil
}

type PlayerAction struct {
	Player           *Player             `protobuf:"bytes,1,req,name=player" json:"player,omitempty"`
	Action           PlayerAction_Action `protobuf:"varint,2,req,name=action,enum=protobuf.PlayerAction_Action" json:"action"`
	Uuid             string              `protobuf:"bytes,3,opt,name=uuid" json:"uuid"`
	Flags            int32               `protobuf:"varint,4,opt,name=flags" json:"flags"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *PlayerAction) Reset()         { *m = PlayerAction{} }
func (m *PlayerAction) String() string { return proto.CompactTextString(m) }
func (*PlayerAction) ProtoMessage()    {}

func (m *PlayerAction) GetPlayer() *Player {
	if m != nil {
		return m.Player
	}
	return nil
}

func (m *PlayerAction) GetAction() PlayerAction_Action {
	if m != nil {
		return m.Action
	}
	return PlayerAction_JOIN
}

func (m *PlayerAction) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *PlayerAction) GetFlags() int32 {
	if m != nil {
		return m.Flags
	}
	return 0
}

func init() {
	proto.RegisterEnum("protobuf.PlayerAction_Action", PlayerAction_Action_name, PlayerAction_Action_value)
}
