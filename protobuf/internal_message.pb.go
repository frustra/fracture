// Code generated by protoc-gen-gogo.
// source: internal_message.proto
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

type InternalMessage struct {
	ChunkRequest      *ChunkRequest      `protobuf:"bytes,1,opt" json:"ChunkRequest,omitempty"`
	ChunkResponse     *ChunkResponse     `protobuf:"bytes,2,opt" json:"ChunkResponse,omitempty"`
	BulkChunkRequest  *BulkChunkRequest  `protobuf:"bytes,3,opt" json:"BulkChunkRequest,omitempty"`
	BulkChunkResponse *BulkChunkResponse `protobuf:"bytes,4,opt" json:"BulkChunkResponse,omitempty"`
	PlayerAction      *PlayerAction      `protobuf:"bytes,5,opt" json:"PlayerAction,omitempty"`
	BlockUpdate       *BlockUpdate       `protobuf:"bytes,6,opt" json:"BlockUpdate,omitempty"`
	ChatMessage       *ChatMessage       `protobuf:"bytes,7,opt" json:"ChatMessage,omitempty"`
	Subscription      *Subscription      `protobuf:"bytes,8,opt" json:"Subscription,omitempty"`
	XXX_unrecognized  []byte             `json:"-"`
}

func (m *InternalMessage) Reset()         { *m = InternalMessage{} }
func (m *InternalMessage) String() string { return proto.CompactTextString(m) }
func (*InternalMessage) ProtoMessage()    {}

func (m *InternalMessage) GetChunkRequest() *ChunkRequest {
	if m != nil {
		return m.ChunkRequest
	}
	return nil
}

func (m *InternalMessage) GetChunkResponse() *ChunkResponse {
	if m != nil {
		return m.ChunkResponse
	}
	return nil
}

func (m *InternalMessage) GetBulkChunkRequest() *BulkChunkRequest {
	if m != nil {
		return m.BulkChunkRequest
	}
	return nil
}

func (m *InternalMessage) GetBulkChunkResponse() *BulkChunkResponse {
	if m != nil {
		return m.BulkChunkResponse
	}
	return nil
}

func (m *InternalMessage) GetPlayerAction() *PlayerAction {
	if m != nil {
		return m.PlayerAction
	}
	return nil
}

func (m *InternalMessage) GetBlockUpdate() *BlockUpdate {
	if m != nil {
		return m.BlockUpdate
	}
	return nil
}

func (m *InternalMessage) GetChatMessage() *ChatMessage {
	if m != nil {
		return m.ChatMessage
	}
	return nil
}

func (m *InternalMessage) GetSubscription() *Subscription {
	if m != nil {
		return m.Subscription
	}
	return nil
}

func init() {
}
func (this *InternalMessage) GetValue() interface{} {
	if this.ChunkRequest != nil {
		return this.ChunkRequest
	}
	if this.ChunkResponse != nil {
		return this.ChunkResponse
	}
	if this.BulkChunkRequest != nil {
		return this.BulkChunkRequest
	}
	if this.BulkChunkResponse != nil {
		return this.BulkChunkResponse
	}
	if this.PlayerAction != nil {
		return this.PlayerAction
	}
	if this.BlockUpdate != nil {
		return this.BlockUpdate
	}
	if this.ChatMessage != nil {
		return this.ChatMessage
	}
	if this.Subscription != nil {
		return this.Subscription
	}
	return nil
}

func (this *InternalMessage) SetValue(value interface{}) bool {
	switch vt := value.(type) {
	case *ChunkRequest:
		this.ChunkRequest = vt
	case *ChunkResponse:
		this.ChunkResponse = vt
	case *BulkChunkRequest:
		this.BulkChunkRequest = vt
	case *BulkChunkResponse:
		this.BulkChunkResponse = vt
	case *PlayerAction:
		this.PlayerAction = vt
	case *BlockUpdate:
		this.BlockUpdate = vt
	case *ChatMessage:
		this.ChatMessage = vt
	case *Subscription:
		this.Subscription = vt
	default:
		return false
	}
	return true
}
