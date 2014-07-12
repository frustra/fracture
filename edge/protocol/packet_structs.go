package protocol

import (
	"encoding/json"
)

type JsonMessage struct {
	Message    string `json:"text"`
	Color      string `json:"color,omitempty"`
	Bold       bool   `json:"bold,omitempty"`
	Italic     bool   `json:"italic,omitempty"`
	Underlined bool   `json:"underlined,omitempty"`
	Strike     bool   `json:"strikethrough,omitempty"`
	Obfuscated bool   `json:"obfuscated,omitempty"`
}

func (jm JsonMessage) Serialize() []byte {
	buf, err := json.Marshal(jm)
	if err != nil {
		return make([]byte, 0)
	} else {
		return append(Varint{uint64(len(buf))}.Bytes(), buf...)
	}
}

func CreateJsonMessage(message, color string) JsonMessage {
	return JsonMessage{message, color, false, false, false, false, false}
}

type StatusResponse struct {
	data map[string]interface{}
}

func (sr StatusResponse) Serialize() []byte {
	buf, err := json.Marshal(sr.data)
	if err != nil {
		return make([]byte, 0)
	} else {
		return append(Varint{uint64(len(buf))}.Bytes(), buf...)
	}
}

func CreateStatusResponse(version string, protocol, players, maxplayers int, description JsonMessage) StatusResponse {
	response := StatusResponse{make(map[string]interface{})}
	_version := make(map[string]interface{})
	_version["name"] = version
	_version["protocol"] = protocol
	response.data["version"] = _version
	_players := make(map[string]interface{})
	_players["max"] = maxplayers
	_players["online"] = players
	response.data["players"] = _players
	response.data["description"] = description
	return response
}
