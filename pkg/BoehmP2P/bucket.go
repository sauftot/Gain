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

func NewBucketList(ownMeta *NodeMeta, register chan internal.ContextWithCancel, ctx internal.ContextWithCancel) *BucketList {
	b := make(map[int][]NodeMeta)
	b[159] = make([]NodeMeta, 0, k)
	bl := new(BucketList)
	bl.buckets = b
	bl.ownMeta = ownMeta
	bl.register = register
	bl.ctx = ctx
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
		default:

		}
	}
	b.buckets[bucketIndex] = append(b.buckets[bucketIndex][1:], node)
}

/*
Add adds a node to the bucket list. It handles splitting buckets, pinging nodes and evicting nodes.
Please check appropriate tests for more information.
*/
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
			adjustedLen := len(b.buckets[index])
			for i := 0; i < adjustedLen; i++ {
				if idBits-1-b.buckets[index][i].MsbXorDist(b.ownMeta) != index {
					b.buckets[index-1] = append(b.buckets[index-1], b.buckets[index][i])
					b.buckets[index] = append(b.buckets[index][:i], b.buckets[index][i+1:]...)
					adjustedLen--
				}
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
		responseChannel := make(chan Message, 10)
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
