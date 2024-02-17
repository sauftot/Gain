package BoehmP2P

import (
	"Gain/internal"
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

/*
Insert tests that check which type of data structure is the most efficient for storing NodeMetas in buckets
*/

func TestNewBucketLength(t *testing.T) {
	register := make(chan internal.ContextWithCancel, 10)
	defer close(register)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct := internal.ContextWithCancel{Ctx: rootCtx, Cancel: cancel}

	testMeta := GenerateOwnMeta()
	bucketList := NewBucketList(testMeta, register, ct)
	t.Log(len(bucketList.buckets))
	for key, value := range bucketList.buckets {
		t.Log("Key:", key, "Value:", value)
	}
}

func TestAddToBucketWithVacantSlot(t *testing.T) {
	register := make(chan internal.ContextWithCancel, 10)
	defer close(register)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct := internal.ContextWithCancel{Ctx: rootCtx, Cancel: cancel}

	testMeta := GenerateOwnMeta()
	bucketList := NewBucketList(testMeta, register, ct)
	testNode := GenerateOwnMeta()
	testNode.Latency = 10
	bucketList.Add(*testNode)

	if len(bucketList.buckets[159]) != 2 {
		t.Error("Bucket 159 Length: Expected 2, got", len(bucketList.buckets[159]))
	}

	if bucketList.buckets[159][0].Latency != -1 {
		t.Error("Own Latency: Expected -1, got", bucketList.buckets[159][0].Latency)
	}

	if bucketList.buckets[159][1].Latency != 10 {
		t.Error("testNode Latency: Expected 10, got", bucketList.buckets[159][1].Latency)
	}
}

func TestAddToBucket20_NoSplit_NoRegisterCalled(t *testing.T) {
	register := make(chan internal.ContextWithCancel, 10)
	defer close(register)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct := internal.ContextWithCancel{Ctx: rootCtx, Cancel: cancel}

	testMeta := GenerateOwnMeta()
	bucketList := NewBucketList(testMeta, register, ct)
	for i := 0; i < k; i++ {
		testNode := GenerateOwnMeta()
		testNode.Latency = 10 + i
		bucketList.Add(*testNode)
	}

	if len(bucketList.buckets[159]) != k {
		t.Error("Bucket 159 Length: Expected", k, "got", len(bucketList.buckets[159]))
	}

	select {
	case <-register:
		t.Error("handler was registered! Expected no handler to be registered")
	default:
		break
	}
}

func TestAddToBucket21_Split_NoRegisterCalled(t *testing.T) {
	register := make(chan internal.ContextWithCancel, 10)
	defer close(register)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct := internal.ContextWithCancel{Ctx: rootCtx, Cancel: cancel}

	testMeta := GenerateOwnMeta()
	bucketList := NewBucketList(testMeta, register, ct)
	for i := 0; i < k; i++ {
		testNode := GenerateOwnMeta()
		testNode.Latency = 10 + i
		bucketList.Add(*testNode)
		t.Log("Adding node with index ", testNode.MsbXorDist(testMeta), " to bucket 159")
	}
	t.Log("Length of bucketList: ", len(bucketList.buckets[idBits-1]))
	t.Log("BucketList: ", bucketList.buckets[idBits-1])

	t.Log("Adding node with index ", testMeta.MsbXorDist(&bucketList.buckets[idBits-1][0]), " to bucket 159")
	nodeForSplit := GenerateOwnMeta()
	nodeForSplit.Latency = 100
	bucketList.Add(*nodeForSplit)

	t.Log("Length of old bucket: ", len(bucketList.buckets[idBits-1]))
	t.Log("Length of new bucket: ", len(bucketList.buckets[idBits-2]))
}

/*

TODO: ANSWER: DO NODES CONTAIN THEIR OWN ID IN THEIR BUCKETS?
A: NO

*/

func TestAddToBucket21_NoSplit_RegisterCalled_PingTimedOut(t *testing.T) {
	register := make(chan internal.ContextWithCancel, 10)
	defer close(register)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct := internal.ContextWithCancel{Ctx: rootCtx, Cancel: cancel}

	/*
		Add 21 NodeMetas with index idBits-1 to bucket idBits-1 to test the eviction functionality,
		which should register a handler context for the network handler to respond to.
	*/

	testMeta := GenerateOwnMeta()
	formatString := "%0.159b1"
	id, suc := big.NewInt(0).SetInt64(0).SetString(fmt.Sprintf(formatString, 0), 2)
	if !suc {
		t.Error("Error creating big.Int from string")
	}
	testMeta.ID = *id
	t.Log("OwnMeta ID: ", testMeta)
	bucketList := NewBucketList(testMeta, register, ct)

	for i := 0; i < k; i++ {
		testNode := GenerateOwnMeta()
		zeroes := strings.Repeat("0", idBits-2-i)
		trailingZeroes := strings.Repeat("0", i+1)
		f := "1" + zeroes + "1" + trailingZeroes
		d, s := big.NewInt(0).SetString(f, 2)
		if !s {
			t.Error("Error creating big.Int from string")
		}
		testNode.ID = *d
		testNode.Latency = 10 + i
		bucketList.Add(*testNode)
		t.Log("Adding node with index ", idBits-1-testNode.MsbXorDist(testMeta), " to bucket 159")
	}
	t.Log("Length of bucketList: ", len(bucketList.buckets[idBits-1]))
	t.Log("BucketList: ", bucketList.buckets[idBits-1])

	nodeForSplit := GenerateOwnMeta()
	nodeForSplit.Latency = 100
	h, z := big.NewInt(0).SetInt64(0).SetString("1"+strings.Repeat("0", idBits-1-100)+"1"+strings.Repeat("0", 100), 2)
	if !z {
		t.Error("Error creating big.Int from string")
	}
	nodeForSplit.ID = *h
	t.Log("Overflow node with index: ", idBits-1-testMeta.MsbXorDist(nodeForSplit), " to bucket 159")
	bucketList.Add(*nodeForSplit)

	t.Log("Length of old bucket: ", len(bucketList.buckets[idBits-1]))
	t.Log("Length of new bucket: ", len(bucketList.buckets[idBits-2]))

	select {
	case <-register:
		t.Log("handler was registered!")
	default:
		t.Error("handler was not registered! Expected a handler to be registered")
	}

	t.Log("This Test won't answer the ping. Expect a timeout of 5 seconds...")

	wg.Wait()
}

func TestAddToBucket21_NoSplit_RegisterCalled_PingReceived(t *testing.T) {
	register := make(chan internal.ContextWithCancel, 10)
	defer close(register)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct := internal.ContextWithCancel{Ctx: rootCtx, Cancel: cancel}

	testMeta := GenerateOwnMeta()
	formatString := "%0.159b1"
	id, suc := big.NewInt(0).SetInt64(0).SetString(fmt.Sprintf(formatString, 0), 2)
	if !suc {
		t.Error("Error creating big.Int from string")
	}
	testMeta.ID = *id
	t.Log("OwnMeta ID: ", testMeta)
	bucketList := NewBucketList(testMeta, register, ct)

	firstNode := GenerateOwnMeta()
	zeroes := strings.Repeat("0", idBits-2-50)
	trailingZeroes := strings.Repeat("0", 51)
	f := "1" + zeroes + "1" + trailingZeroes
	d, s := big.NewInt(0).SetString(f, 2)
	if !s {
		t.Error("Error creating big.Int from string")
	}
	firstNode.ID = *d
	firstNode.Latency = 1
	bucketList.Add(*firstNode)

	for i := 0; i < k-1; i++ {
		testNode := GenerateOwnMeta()
		zeroes := strings.Repeat("0", idBits-2-i)
		trailingZeroes := strings.Repeat("0", i+1)
		f := "1" + zeroes + "1" + trailingZeroes
		d, s := big.NewInt(0).SetString(f, 2)
		if !s {
			t.Error("Error creating big.Int from string")
		}
		testNode.ID = *d
		testNode.Latency = 10 + i
		bucketList.Add(*testNode)
		t.Log("Adding node with index ", idBits-1-testNode.MsbXorDist(testMeta), " to bucket 159")
	}
	t.Log("Length of bucketList: ", len(bucketList.buckets[idBits-1]))
	t.Log("BucketList: ", bucketList.buckets[idBits-1])

	nodeForSplit := GenerateOwnMeta()
	nodeForSplit.Latency = 100
	h, z := big.NewInt(0).SetInt64(0).SetString("1"+strings.Repeat("0", idBits-1-100)+"1"+strings.Repeat("0", 100), 2)
	if !z {
		t.Error("Error creating big.Int from string")
	}
	nodeForSplit.ID = *h
	t.Log("Overflow node with index: ", idBits-1-testMeta.MsbXorDist(nodeForSplit), " to bucket 159")
	bucketList.Add(*nodeForSplit)

	t.Log("Length of old bucket: ", len(bucketList.buckets[idBits-1]))
	t.Log("Length of new bucket: ", len(bucketList.buckets[idBits-2]))

	select {
	case re := <-register:
		t.Log("Handler was registered!")
		re.Ctx.Value("callback").(chan Message) <- *NewMessage(re.Ctx.Value("RPCID").(*big.Int), firstNode, PINGRES, nil)
	default:
		t.Error("handler was not registered! Expected a handler to be registered")
	}

	wg.Wait()

	if len(bucketList.buckets[idBits-1]) != k {
		t.Error("Bucket 159 Length: Expected", k, "got", len(bucketList.buckets[idBits-1]))
	}

	if bucketList.buckets[idBits-1][0].ID.Cmp(&firstNode.ID) == 0 {
		t.Error("First node was not moved to the back of the bucket")
	}

	if bucketList.buckets[idBits-1][k-1].ID.Cmp(&firstNode.ID) != 0 {
		t.Error("Last Node is not old first node! Move was incorrect.")
	}
	t.Log("Node was moved to the back of the bucket!")
}
