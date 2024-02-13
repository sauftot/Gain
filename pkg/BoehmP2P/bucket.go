package BoehmP2P

import (
	"Gain/internal"
	"context"
	"net"
	"time"
)

// the key of the map is the most significant bit of the xor result that the bucket covers
type BucketList struct {
	ownMeta       *NodeMeta
	buckets       map[int][]NodeMeta
	ctx           internal.ContextWithCancel
	pingResponses chan *Message
}

func NewBucketList(ownMeta *NodeMeta) *BucketList {
	b := make(map[int][]NodeMeta)
	b[159] = make([]NodeMeta, 0, k)
	b[159] = append(b[159], *ownMeta)
	bl := new(BucketList)
	bl.buckets = b
	bl.ownMeta = ownMeta
	return bl
}

func (b *BucketList) pingHandler(ctx internal.ContextWithCancel, bucketIndex int, node NodeMeta) {
	// ping the node and wait for a response till the context is cancelled or timed out
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(b.buckets[bucketIndex][0].IP), Port: port})
	if err != nil {
		nodeLogger.Error("Error pinging node: ", err)
		ctx.Cancel()
		// evict node and insert new one
	} else {
		msgID := GenerateID()
		msg := NewMessage(&msgID, b.ownMeta, PING, nil)
		_, err = conn.Write(msg.ToBytes())
	}
	for {
		select {
		case <-ctx.Ctx.Done():
			return
		case msg := <-b.pingResponses:
			// find the node in the bucket
			// if it is not found, ignore the message
			// if it is found, move it to the back of the bucket
		}
	}

}

func (b *BucketList) Add(node NodeMeta) {
	// xor the node with the ownMeta
	// find the most significant bit of the xor result
	// add the node to the bucket with the most significant bit as the key
	msb := b.ownMeta.MsbXorDist(&node)
	var index = 159
	for key := range b.buckets {
		if key >= msb && key < index {
			index = key
		}
	}
	if len(b.buckets[index]) < k {
		b.buckets[index] = append(b.buckets[index], node)
	} else {
		// check if we can split the bucket, if not then check if we should evict a node
		_, ok := b.buckets[index-1]
		if !ok {
			b.buckets[index-1] = make([]NodeMeta, 0, k)
			// move the appropriate nodes from the index bucket to the new bucket
			indexesToRemove := make([]int, 0, k)
			for i, n := range b.buckets[index] {
				if n.MsbXorDist(b.ownMeta) != index {
					b.buckets[index-1] = append(b.buckets[index-1], n)
					indexesToRemove = append(indexesToRemove, i)
				}
			}
			for _, i := range indexesToRemove {
				b.buckets[index] = append(b.buckets[index][:i], b.buckets[index][i+1:]...)
			}
			// add the new node to the correct bucket
			if len(b.buckets[index]) < k {
				b.buckets[index] = append(b.buckets[index], node)
				return
			}
		}
		// ping the first node in the bucket
		// if it does not respond, remove it and add the new node
		// if it does respond, move it to the back of the bucket
		// message := NewMessage(big.NewInt(0), b.ownMeta, PING, nil)
		ct, cancel := context.WithDeadline(b.ctx.Ctx, time.Now().Add(5*time.Second))
		cwc := internal.ContextWithCancel{Ctx: ct, Cancel: cancel}
		wg.Add(1)
		go b.pingHandler(cwc, index-1, node)
	}
}
