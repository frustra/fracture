// Code generated by protoc-gen-gogo.
// source: protobuf/node_meta.proto
// DO NOT EDIT!

package protobuf

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type NodeMeta struct {
	Addr             *string `protobuf:"bytes,1,req,name=addr" json:"addr,omitempty"`
	Type             *string `protobuf:"bytes,2,req,name=type" json:"type,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *NodeMeta) Reset()         { *m = NodeMeta{} }
func (m *NodeMeta) String() string { return proto.CompactTextString(m) }
func (*NodeMeta) ProtoMessage()    {}

func (m *NodeMeta) GetAddr() string {
	if m != nil && m.Addr != nil {
		return *m.Addr
	}
	return ""
}

func (m *NodeMeta) GetType() string {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return ""
}

func init() {
}
