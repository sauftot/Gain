package BoehmP2P

import (
	"crypto/rand"
	"math/big"
)

type NodeMeta struct {
	IP      string
	Port    string
	ID      big.Int
	Latency int
}

func (id *NodeMeta) Xor(other *NodeMeta) *big.Int {
	return new(big.Int).Xor(&id.ID, &other.ID)
}

/*
MsbXorDist returns the floor of the log2 of the xor result
*/
func (id *NodeMeta) MsbXorDist(other *NodeMeta) int {
	xor := id.Xor(other)
	var msb = 0
	for i := idBits; i >= 0; i-- {
		if xor.Bit(i) != 0 {
			msb = i
		}
	}
	return msb
}

func GenerateOwnMeta() *NodeMeta {
	id := GenerateID()
	//nodeLogger.Log("Generated Node ID: " + id.String())
	meta := new(NodeMeta)
	meta.IP = ""
	meta.Port = ""
	meta.ID = id
	meta.Latency = -1
	return meta
}

func GenerateID() big.Int {
	intBytes := make([]byte, 20)
	_, err := rand.Read(intBytes)
	if err != nil {
		nodeLogger.Log("Error generating NodeMeta! panic!")
		panic(err)
	}
	id := *big.NewInt(0).SetBytes(intBytes)
	return id
}
