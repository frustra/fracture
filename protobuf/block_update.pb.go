// Code generated by protoc-gen-gogo.
// source: block_update.proto
// DO NOT EDIT!

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	block_update.proto
	bulk_chunk_request.proto
	bulk_chunk_response.proto
	chat_message.proto
	chunk_request.proto
	chunk_response.proto
	internal_message.proto
	node_meta.proto
	player_action.proto
	player.proto

It has these top-level messages:
	BlockUpdate
*/
package protobuf

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// discarding unused import gogoproto "code.google.com/p/gogoprotobuf/gogoproto/gogo.pb"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type BlockUpdate struct {
	X                int64  `protobuf:"zigzag64,1,req,name=x" json:"x"`
	Y                uint32 `protobuf:"varint,2,req,name=y" json:"y"`
	Z                int64  `protobuf:"zigzag64,3,req,name=z" json:"z"`
	Destroy          bool   `protobuf:"varint,4,opt,name=destroy" json:"destroy"`
	Uuid             string `protobuf:"bytes,5,opt,name=uuid" json:"uuid"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *BlockUpdate) Reset()         { *m = BlockUpdate{} }
func (m *BlockUpdate) String() string { return proto.CompactTextString(m) }
func (*BlockUpdate) ProtoMessage()    {}

func (m *BlockUpdate) GetX() int64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *BlockUpdate) GetY() uint32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *BlockUpdate) GetZ() int64 {
	if m != nil {
		return m.Z
	}
	return 0
}

func (m *BlockUpdate) GetDestroy() bool {
	if m != nil {
		return m.Destroy
	}
	return false
}

func (m *BlockUpdate) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func init() {
}