package main

import (
	"log"
	"math/rand"
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
		key.previousPointer = currentPointer
		pointer.previousKey = key
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
		key.previousKey = currentTailKey
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

func Test_AddKeyToLeaf_emptyList(t *testing.T) {
	node := &BPlusTreeNode{isLeaf: true}
	node.insertToLeafNode(1, "value")
	if node.countOfKeys != 1 {
		t.FailNow()
	}
	if node.leafHead.value != 1 {
		t.FailNow()
	}
	if node.leafHead.nextKey != nil {
		t.FailNow()
	}
}

func Test_AddKeyToLeaf_appendToTail(t *testing.T) {
	leafNode := BPlusTreeNode{
		isLeaf:      true,
		countOfKeys: 1,
	}
	leafNode.leafHead = &bPlusTreeKey{value: 1}
	currentTailKey := leafNode.leafHead
	for i := 1; i < 3; i++ {
		leafNode.countOfKeys = leafNode.countOfKeys + 1
		key := &bPlusTreeKey{
			value: int64(i*10 + 1),
		}
		currentTailKey.nextKey = key
		key.previousKey = currentTailKey
		currentTailKey = currentTailKey.nextKey
	}
	leafNode.insertToLeafNode(31, "value")
	currentKey := leafNode.leafHead
	for i := 0; i < 4; i++ {
		if currentKey.value != int64(i*10+1) {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
	}
}

func Test_AddKeyWithinLeafArray(t *testing.T) {
	leafNode := BPlusTreeNode{
		isLeaf: true,
	}
	leafNode.leafHead = &bPlusTreeKey{value: 1}
	currentTailKey := leafNode.leafHead
	for i := 1; i < 10; i++ {
		if i == 6 {
			continue
		}
		leafNode.countOfKeys = leafNode.countOfKeys + 1
		key := &bPlusTreeKey{
			value: int64(i*10 + 1),
		}
		currentTailKey.nextKey = key
		key.previousKey = currentTailKey
		currentTailKey = currentTailKey.nextKey
	}
	leafNode.insertToLeafNode(61, "value")
	currentKey := leafNode.leafHead
	for i := 0; i < 10; i++ {
		if currentKey.value != int64(i*10+1) {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
	}
}

func Test_AddKeyAtStartOfList(t *testing.T) {
	leafNode := BPlusTreeNode{
		isLeaf: true,
	}
	leafNode.leafHead = &bPlusTreeKey{value: 1}
	currentTailKey := leafNode.leafHead
	for i := 1; i < 3; i++ {
		leafNode.countOfKeys = leafNode.countOfKeys + 1
		key := &bPlusTreeKey{
			value: int64(i*10 + 1),
		}
		currentTailKey.nextKey = key
		key.previousKey = currentTailKey
		currentTailKey = currentTailKey.nextKey
	}
	leafNode.insertToLeafNode(-9, "value")
	currentKey := leafNode.leafHead
	for i := -1; i < 3; i++ {
		if currentKey.value != int64(i*10+1) {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
	}
}

func Test_Insert700Values(t *testing.T) {
	tree := initTree(2)
	for i := 1; i < 700; i++ {
		tree.Insert(int64(i), "")
	}
}

func Test_GetPointer_oneKey(t *testing.T) {
	internalNode := BPlusTreeNode{
		isLeaf: false,
	}
	internalNode.internalNodeHead = &bPlusTreePointer{}
	currentPointer := internalNode.internalNodeHead
	for i := 0; i < 2; i++ {
		internalNode.countOfKeys = internalNode.countOfKeys + 1
		key := &bPlusTreeKey{
			value: int64(i*10 + 1),
		}
		pointer := &bPlusTreePointer{}
		key.nextPointer = pointer
		currentPointer.nextKey = key
		key.previousPointer = currentPointer
		pointer.previousKey = key
		currentPointer = currentPointer.nextKey.nextPointer
	}
	poiner := internalNode.getPointer(2)
	log.Println(poiner)
}

func Test_InsertRandom700Values(t *testing.T) {
	tree := initTree(2)
	for i := 0; i < 300; i++ {
		tree.Insert(rand.Int63n(10000), "")
	}
}
