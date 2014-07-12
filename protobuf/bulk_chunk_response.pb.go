// Code generated by protoc-gen-gogo.
// source: bulk_chunk_response.proto
// DO NOT EDIT!

package protobuf

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type BulkChunkResponse struct {
	Chunks           []*ChunkResponse `protobuf:"bytes,1,rep,name=chunks" json:"chunks,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *BulkChunkResponse) Reset()         { *m = BulkChunkResponse{} }
func (m *BulkChunkResponse) String() string { return proto.CompactTextString(m) }
func (*BulkChunkResponse) ProtoMessage()    {}

func (m *BulkChunkResponse) GetChunks() []*ChunkResponse {
	if m != nil {
		return m.Chunks
	}
	return nil
}

func init() {
}
