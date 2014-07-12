// Code generated by protoc-gen-gogo.
// source: chunk_request.proto
// DO NOT EDIT!

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	chunk_request.proto
	chunk_response.proto
	internal_message.proto
	node_meta.proto

It has these top-level messages:
	ChunkRequest
*/
package protobuf

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type ChunkRequest struct {
	X                *int64 `protobuf:"varint,1,req,name=x" json:"x,omitempty"`
	Z                *int64 `protobuf:"varint,2,req,name=z" json:"z,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ChunkRequest) Reset()         { *m = ChunkRequest{} }
func (m *ChunkRequest) String() string { return proto.CompactTextString(m) }
func (*ChunkRequest) ProtoMessage()    {}

func (m *ChunkRequest) GetX() int64 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *ChunkRequest) GetZ() int64 {
	if m != nil && m.Z != nil {
		return *m.Z
	}
	return 0
}

func init() {
}
