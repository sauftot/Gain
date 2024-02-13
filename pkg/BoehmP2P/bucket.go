package BoehmP2P

import "Gain/internal"

// the key of the map is the most significant bit of the xor result that the bucket covers
type BucketList struct {
	ownMeta *NodeMeta
	buckets map[int][]NodeMeta
	ctx     internal.ContextWithCancel
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

func (b *BucketList) Add(node NodeMeta, toNode chan Message) {
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
		buck, ok := b.buckets[index-1]
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
		// TODO: how to implement pinging a node?
	}
}
