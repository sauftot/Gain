package BoehmP2P

import (
	json2 "encoding/json"
	"math/big"
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
	RPCID big.Int
	From  NodeMeta
	RPC   uint8
	Data  []byte
}

func NewMessage(rpcID *big.Int, from *NodeMeta, rpc uint8, data []byte) *Message {
	return &Message{
		RPCID: *rpcID,
		From:  *from,
		RPC:   rpc,
		Data:  data,
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
