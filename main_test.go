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

func TestInsertToLeafNode_ok(t *testing.T) {
	countOfKeys := 9
	node := initOneTestLeafNode(countOfKeys)
	node.insertToLeafNode(25, "")
	keys := []int64{10, 20, 25, 30, 40, 50, 60, 70, 80, 90}
	currentKey := node.leafHead
	for _, key := range keys {
		if currentKey.value != key {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
	}
}

func TestInsertToLeafNode_atStartPosition(t *testing.T) {
	countOfKeys := 9
	node := initOneTestLeafNode(countOfKeys)
	node.insertToLeafNode(5, "")
	keys := []int64{5, 10, 20, 30, 40, 50, 60, 70, 80, 90}
	currentKey := node.leafHead
	for _, key := range keys {
		if currentKey.value != key {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
	}
}

func TestInsertToLeafNode_atEndPosition(t *testing.T) {
	countOfKeys := 9
	node := initOneTestLeafNode(countOfKeys)
	node.insertToLeafNode(100, "")
	keys := []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	currentKey := node.leafHead
	for _, key := range keys {
		if currentKey.value != key {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
	}
}

func TestInsertToLeafNode_emptyList(t *testing.T) {
	node := initOneTestLeafNode(0)
	node.insertToLeafNode(100, "")
	if !(node.leafHead != nil && node.leafHead.value == 100 && node.leafHead.nextKey == nil) {
		t.FailNow()
	}
}

func TestDeleteFromLeafNode_ok(t *testing.T) {
	node := initOneTestLeafNode(9)
	node.deleteFromLeafNode(40)
	keys := []int64{10, 20, 30, 50, 60, 70, 80, 90}
	currentKey := node.leafHead
	for _, key := range keys {
		if currentKey.value != key {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
	}
}

func TestDeleteFromLeafNode_atStartPosition(t *testing.T) {
	node := initOneTestLeafNode(9)
	node.deleteFromLeafNode(10)
	keys := []int64{20, 30, 40, 50, 60, 70, 80, 90}
	currentKey := node.leafHead
	for _, key := range keys {
		if currentKey.value != key {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
	}
}

func TestDeleteFromLeafNode_atEndPosition(t *testing.T) {
	node := initOneTestLeafNode(9)
	node.deleteFromLeafNode(90)
	keys := []int64{10, 20, 30, 40, 50, 60, 70, 80}
	currentKey := node.leafHead
	for _, key := range keys {
		if currentKey.value != key {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
	}
}

func TestDeleteFromLeafNode_emptyList(t *testing.T) {
	node := initOneTestLeafNode(0)
	success := node.deleteFromLeafNode(90)
	if !(node.leafHead == nil && !success) {
		t.FailNow()
	}
}

func initOneTestLeafNode(countOfKeys int) *BPlusTreeNode {
	value := int64(10)
	node := BPlusTreeNode{
		isLeaf:      true,
		countOfKeys: countOfKeys,
	}
	if countOfKeys > 0 {
		node.leafHead = &bPlusTreeKey{
			value: value,
		}
	}
	previousKey := node.leafHead
	for i := 1; i < countOfKeys; i++ {
		value = value + 10
		newKey := bPlusTreeKey{
			value:       value,
			previousKey: previousKey,
		}
		previousKey.nextKey = &newKey
		previousKey = previousKey.nextKey
	}
	return &node
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
