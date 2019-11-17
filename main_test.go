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
		leafNode.countOfKeys = leafNode.countOfKeys + 1
		key := &bPlusTreeKey{
			value: int64(i*10 + 1),
		}
		pointer := &bPlusTreePointer{}
		key.nextPointer = pointer
		currentPointer.nextKey = key
		currentPointer = currentPointer.nextKey.nextPointer
	}
	subtree := leafNode.cutByTwoNodes()
	currentPointer = leafNode.internalNodeHead
	i := 1
	for currentPointer.nextKey != nil {
		if currentPointer.nextKey.value != int64(i) {
			t.FailNow()
		}
		currentPointer = currentPointer.nextKey.nextPointer
		i = i + 10
	}
	if subtree.nextKey.value != int64(i) {
		t.FailNow()
	}
	i = i + 10
	currentPointer = subtree.nextKey.nextPointer.childNode.internalNodeHead
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
	for i := 1; i < 6; i++ {
		leafNode.countOfKeys = leafNode.countOfKeys + 1
		key := &bPlusTreeKey{
			value: int64(i*10 + 1),
		}
		currentTailKey.nextKey = key
		currentTailKey = currentTailKey.nextKey
	}
	subtree := leafNode.cutByTwoNodes()
	currentKey := leafNode.leafHead
	i := 1
	for currentKey != nil {
		if currentKey.value != int64(i) {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
		i = i + 10
	}
	if subtree.nextKey.value != int64(i) {
		t.FailNow()
	}
	currentKey = subtree.nextKey.nextKey
	for currentKey != nil {
		if currentKey.value != int64(i) {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
		i = i + 10
	}
	t.Log()
}

func Test_Insert10Values(t *testing.T) {
	tree := initTree(3)
	for i := 1; i < 7; i++ {
		tree.Insert(int64(i), "")
	}
	tree.PrintTree()
}
