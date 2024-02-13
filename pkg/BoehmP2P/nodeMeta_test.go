package BoehmP2P

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	id := GenerateOwnMeta()
	if id == nil {
		t.Error("GenerateOwnMeta returned nil")
	}
	t.Log("Generated NodeMeta: ", id.ID.String())
}

func TestID_Xor(t *testing.T) {
	id1 := GenerateOwnMeta()
	id2 := GenerateOwnMeta()
	t.Log("ID1: ", id1.ID.String())
	t.Log("ID2: ", id2.ID.String())

	xor := id1.Xor(id2)
	if xor == nil {
		t.Error("Xor returned nil")
	}
	t.Log(xor.String())

	xorD := id1.MsbXorDist(id2)
	t.Log("Xor distance (2^", xorD, ")")
}
