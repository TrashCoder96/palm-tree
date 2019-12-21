package main

import "testing"

func TestGetPointer_ok(t *testing.T) {
	node := initOneTestInternalNode(9)
	pointer := node.getPointer(11)
	if !(pointer.previousKey.value == 10 && pointer.nextKey.value == 20) {
		t.FailNow()
	}
}

func TestGetPointer_KeyLessThan10(t *testing.T) {
	node := initOneTestInternalNode(9)
	pointer := node.getPointer(3)
	if !(pointer.previousKey == nil && pointer.nextKey.value == 10) {
		t.FailNow()
	}
}

func TestGetPointer_KeyMoreThan90(t *testing.T) {
	node := initOneTestInternalNode(9)
	pointer := node.getPointer(100)
	if !(pointer.previousKey.value == 90 && pointer.nextKey == nil) {
		t.FailNow()
	}
}

func initOneTestInternalNode(countOfKeys int) *BPlusTreeNode {
	node := BPlusTreeNode{
		isLeaf:      false,
		countOfKeys: countOfKeys,
	}
	value := int64(0)
	node.internalNodeHead = &bPlusTreePointer{}
	previousPointer := node.internalNodeHead
	for i := 0; i < countOfKeys; i++ {
		value = value + 10
		newKey := bPlusTreeKey{
			previousPointer: previousPointer,
			value:           value,
		}
		previousPointer.nextKey = &newKey
		newPointer := bPlusTreePointer{
			previousKey: &newKey,
		}
		newKey.nextPointer = &newPointer
		previousPointer = &newPointer
	}
	return &node
}
