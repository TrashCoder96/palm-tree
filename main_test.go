package main

import (
	"testing"
)

func Test_DevideByTwoInternalNodes(t *testing.T) {
	leafNode := BPlusTreeNode{
		isLeaf: false,
	}
	leafNode.internalNodeHead = &bPlusTreePointer{}
	currentPointer := leafNode.internalNodeHead
	for i := 0; i < 10; i++ {
		leafNode.count = leafNode.count + 1
		key := &bPlusTreeKey{
			value: int64(i*10 + 1),
		}
		pointer := &bPlusTreePointer{}
		key.nextPointer = pointer
		currentPointer.nextKey = key
		currentPointer = currentPointer.nextKey.nextPointer
	}
	tailNode, middleKey := leafNode.cutTailWithMiddleKey()
	currentPointer = leafNode.internalNodeHead
	i := 1
	for currentPointer.nextKey != nil {
		if currentPointer.nextKey.value != int64(i) {
			t.FailNow()
		}
		currentPointer = currentPointer.nextKey.nextPointer
		i = i + 10
	}
	if middleKey.value != int64(i) {
		t.FailNow()
	}
	i = i + 10
	currentPointer = tailNode.internalNodeHead
	for currentPointer.nextKey != nil {
		if currentPointer.nextKey.value != int64(i) {
			t.FailNow()
		}
		currentPointer = currentPointer.nextKey.nextPointer
		i = i + 10
	}
}

func Test_DevideByTwoLeafNodes(t *testing.T) {
	leafNode := BPlusTreeNode{
		isLeaf: true,
	}
	leafNode.leafHead = &bPlusTreeKey{value: 1}
	currentTailKey := leafNode.leafHead
	for i := 1; i < 9; i++ {
		leafNode.count = leafNode.count + 1
		key := &bPlusTreeKey{
			value: int64(i*10 + 1),
		}
		currentTailKey.nextKey = key
		currentTailKey = currentTailKey.nextKey
	}
	tailNode, middleKey := leafNode.cutTailWithMiddleKey()
	currentKey := leafNode.leafHead
	i := 1
	for currentKey != nil {
		if currentKey.value != int64(i) {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
		i = i + 10
	}
	if middleKey.value != int64(i) {
		t.FailNow()
	}
	currentKey = tailNode.leafHead
	for currentKey != nil {
		if currentKey.value != int64(i) {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
		i = i + 10
	}

}
