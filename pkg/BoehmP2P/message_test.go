package BoehmP2P

import "testing"

// TestParsing tests the marshalling and parsing of messages
func TestParsing(t *testing.T) {
	id := GenerateOwnMeta()
	id.latency = 1
	rpc := PING
	rpcID := GenerateID()
	msg := NewMessage(&rpcID, id, rpc, []byte("test"))
	bytes := msg.ToBytes()
	t.Log("Message: ", msg)
	t.Log("Bytes: ", bytes)
	parsedMsg := ParseMessage(bytes)
	t.Log("Parsed Message: ", parsedMsg)
	if parsedMsg.RPCID.String() != rpcID.String() {
		t.Log("Original RPCID: ", rpcID.String())
		t.Log("Parsed RPCID: ", parsedMsg.RPCID.String())
		t.Error("RPCID does not match")
	}
	if parsedMsg.From.ID.String() != id.ID.String() {
		t.Log("Original From ID: ", id.ID.String())
		t.Log("Parsed From ID: ", parsedMsg.From.ID.String())
		t.Error("From ID does not match")
	}
	if parsedMsg.RPC != rpc {
		t.Log("Original RPC: ", rpc)
		t.Log("Parsed RPC: ", parsedMsg.RPC)
		t.Error("RPC does not match")
	}
	if string(parsedMsg.Data) != "test" {
		t.Log("Original Data: ", "test")
		t.Log("Parsed Data: ", string(parsedMsg.Data))
		t.Error("Data does not match")
	}
}

// TestParsingNilData tests the marshalling and parsing of messages with nil data
func TestParsingNilData(t *testing.T) {
	id := GenerateOwnMeta()
	id.latency = 1
	rpc := PING
	rpcID := GenerateID()
	msg := NewMessage(&rpcID, id, rpc, nil)
	bytes := msg.ToBytes()
	t.Log("Message: ", msg)
	t.Log("Bytes: ", bytes)
	parsedMsg := ParseMessage(bytes)
	t.Log("Parsed Message: ", parsedMsg)
	if parsedMsg.RPCID.String() != rpcID.String() {
		t.Log("Original RPCID: ", rpcID.String())
		t.Log("Parsed RPCID: ", parsedMsg.RPCID.String())
		t.Error("RPCID does not match")
	}
	if parsedMsg.From.ID.String() != id.ID.String() {
		t.Log("Original From ID: ", id.ID.String())
		t.Log("Parsed From ID: ", parsedMsg.From.ID.String())
		t.Error("From ID does not match")
	}
	if parsedMsg.RPC != rpc {
		t.Log("Original RPC: ", rpc)
		t.Log("Parsed RPC: ", parsedMsg.RPC)
		t.Error("RPC does not match")
	}
	if parsedMsg.Data != nil {
		t.Log("Original Data: ", nil)
		t.Log("Parsed Data: ", parsedMsg.Data)
		t.Error("Data does not match")
	}
}
