package BoehmP2P

import (
	"testing"
)

/*
Insert tests that check which type of data structure is the most efficient for storing NodeMetas in buckets
*/

func TestNewBucketLength(t *testing.T) {
	testMeta := GenerateOwnMeta()
	bucketList := NewBucketList(testMeta)
	t.Log(len(bucketList.buckets))
	for key, value := range bucketList.buckets {
		t.Log("Key:", key, "Value:", value)
	}
}
