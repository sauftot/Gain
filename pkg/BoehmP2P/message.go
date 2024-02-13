package BoehmP2P

import (
	json2 "encoding/json"
)

const (
	// PING FIND_NODE FIND_VALUE STORE BROADCAST RPCs
	PING       uint8 = 66
	FIND_NODE  uint8 = 67
	FIND_VALUE uint8 = 68
	STORE      uint8 = 69
	BROADCAST  uint8 = 70
	DOWNLOAD   uint8 = 71
)

type Message struct {
	From *NodeMeta
	RPC  uint8
	data []byte
}

func NewMessage(from *NodeMeta, rpc uint8, data []byte) *Message {
	return &Message{
		From: from,
		RPC:  rpc,
		data: data,
	}
}

func (m *Message) ToBytes() []byte {
	json, err := json2.Marshal(m)
	if err != nil {
		nodeLogger.Error("Error marshalling message to bytes: ", err)
		return nil
	}
	return json
}

func ParseMessage(data []byte) *Message {
	var receivedMsg Message
	err := json2.Unmarshal(data, &receivedMsg)
	if err != nil {
		// Handle error (e.g., invalid JSON or incorrect data structure)
		nodeLogger.Error("Error unmarshalling message from bytes: ", err)
		return nil
	}
	return &receivedMsg
}
