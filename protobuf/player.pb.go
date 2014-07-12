// Code generated by protoc-gen-gogo.
// source: player.proto
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

type Player struct {
	Uuid             string  `protobuf:"bytes,1,req,name=uuid" json:"uuid"`
	Username         string  `protobuf:"bytes,2,req,name=username" json:"username"`
	X                float64 `protobuf:"fixed64,3,req,name=x" json:"x"`
	Y                float64 `protobuf:"fixed64,4,req,name=y" json:"y"`
	Z                float64 `protobuf:"fixed64,5,req,name=z" json:"z"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Player) Reset()         { *m = Player{} }
func (m *Player) String() string { return proto.CompactTextString(m) }
func (*Player) ProtoMessage()    {}

func (m *Player) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Player) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Player) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Player) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Player) GetZ() float64 {
	if m != nil {
		return m.Z
	}
	return 0
}

func init() {
}
