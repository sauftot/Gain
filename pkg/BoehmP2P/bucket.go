package BoehmP2P

import (
	"Gain/internal"
	"context"
	"math/big"
	"net"
	"time"
)

// the key of the map is the most significant bit of the xor result that the bucket covers
type BucketList struct {
	ownMeta  *NodeMeta
	buckets  map[int][]NodeMeta
	ctx      internal.ContextWithCancel
	register chan internal.ContextWithCancel
}

func NewBucketList(ownMeta *NodeMeta, register chan internal.ContextWithCancel) *BucketList {
	b := make(map[int][]NodeMeta)
	b[159] = make([]NodeMeta, 0, k)
	b[159] = append(b[159], *ownMeta)
	bl := new(BucketList)
	bl.buckets = b
	bl.ownMeta = ownMeta
	bl.register = register
	return bl
}

func (b *BucketList) pingHandler(ctx internal.ContextWithCancel, bucketIndex int, node NodeMeta) {
	defer wg.Done()
	pingResponses := ctx.Ctx.Value("callback").(chan Message)
	defer close(pingResponses)
	// ping the node and wait for a response till the context is cancelled or timed out
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(b.buckets[bucketIndex][0].IP), Port: port})
	if err != nil {
		nodeLogger.Error("Error pinging node: ", err)
		ctx.Cancel()
		// evict node and insert new one
	} else {
		msg := NewMessage(ctx.Ctx.Value("RPCID").(*big.Int), b.ownMeta, PING, nil)
		_, err = conn.Write(msg.ToBytes())
	}
	for running := true; running; {
		select {
		case <-ctx.Ctx.Done():
			running = false
		case msg := <-pingResponses:
			if b.buckets[bucketIndex][0].ID.Cmp(&msg.From.ID) == 0 {
				b.buckets[bucketIndex] = append(b.buckets[bucketIndex][1:], b.buckets[bucketIndex][0])
				return
			} else {
				// bucket was changed before ping response. Maybe the pinged node sent us something before responding to the ping, we can return
				return
			}
		}
	}
	b.buckets[bucketIndex] = append(b.buckets[bucketIndex][1:], node)
}

func (b *BucketList) Add(node NodeMeta) {
	// xor the node with our own ID
	// find the most significant bit of the xor result
	// add the node to the bucket with the most significant bit as the key

	// find the bucket for the xor result MSB
	msb := b.ownMeta.MsbXorDist(&node)
	var index = 159
	for key := range b.buckets {
		if key >= msb && key < index {
			index = key
		}
	}

	// check if the bucket has space
	if len(b.buckets[index]) < k {
		b.buckets[index] = append(b.buckets[index], node)
	} else {
		// check if we can split the bucket

		// if the next bucket does not exist, create it and move the appropriate nodes from the index bucket to the new bucket
		_, exists := b.buckets[index-1]
		if !exists {
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
		// if the bucket does exist, or we could not move any nodes to the new bucket because they were all correctly placed, we need to evict

		// ping the first node in the bucket
		// if it does not respond, remove it and add the new node
		// if it does respond, move it to the back of the bucket

		// create a new context with RPC ID
		pingId := GenerateID()
		ct := context.WithValue(b.ctx.Ctx, "RPCID", &pingId)

		// create a callback channel for the network handler to send the response to
		responseChannel := make(chan Message)
		ct2 := context.WithValue(ct, "callback", responseChannel)

		// generate a deadline for the context
		ct3, cancel := context.WithDeadline(ct2, time.Now().Add(5*time.Second))
		cwc := internal.ContextWithCancel{Ctx: ct3, Cancel: cancel}

		// send the context to the network handler
		b.register <- cwc

		// start the ping handler
		wg.Add(1)
		go b.pingHandler(cwc, index, node)
	}
}
