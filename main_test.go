package main

import "testing"

func TestGetPointer_ok(t *testing.T) {
	node := initOneTestInternalNode(9)
	pointer := node.getPointer(11)
	assertCondition := pointer.previousKey.value == 10 &&
		pointer.nextKey.value == 20
	if !assertCondition {
		t.FailNow()
	}
}

func TestGetPointer_KeyLessThan10(t *testing.T) {
	node := initOneTestInternalNode(9)
	pointer := node.getPointer(3)
	assertCondition := pointer.previousKey == nil &&
		pointer.nextKey.value == 10
	if !assertCondition {
		t.FailNow()
	}
}

func TestGetPointer_KeyMoreThan90(t *testing.T) {
	node := initOneTestInternalNode(9)
	pointer := node.getPointer(100)
	assertCondition := pointer.previousKey.value == 90 &&
		pointer.nextKey == nil
	if !assertCondition {
		t.FailNow()
	}
}

func TestInsertToLeafNode_ok(t *testing.T) {
	countOfKeys := 9
	node := initOneTestLeafNode(countOfKeys)
	node.insertToLeafNode(25, "")
	checkNode([]int64{10, 20, 25, 30, 40, 50, 60, 70, 80, 90}, node, t)
}

func TestInsertToLeafNode_atStartPosition(t *testing.T) {
	countOfKeys := 9
	node := initOneTestLeafNode(countOfKeys)
	node.insertToLeafNode(5, "")
	checkNode([]int64{5, 10, 20, 30, 40, 50, 60, 70, 80, 90}, node, t)
}

func TestInsertToLeafNode_atEndPosition(t *testing.T) {
	countOfKeys := 9
	node := initOneTestLeafNode(countOfKeys)
	node.insertToLeafNode(100, "")
	checkNode([]int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}, node, t)
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
	checkNode([]int64{10, 20, 30, 50, 60, 70, 80, 90}, node, t)
}

func TestDeleteFromLeafNode_atStartPosition(t *testing.T) {
	node := initOneTestLeafNode(9)
	node.deleteFromLeafNode(10)
	checkNode([]int64{20, 30, 40, 50, 60, 70, 80, 90}, node, t)
}

func TestDeleteFromLeafNode_atEndPosition(t *testing.T) {
	node := initOneTestLeafNode(9)
	node.deleteFromLeafNode(90)
	checkNode([]int64{10, 20, 30, 40, 50, 60, 70, 80}, node, t)
}

func TestDeleteFromLeafNode_emptyList(t *testing.T) {
	node := initOneTestLeafNode(0)
	success := node.deleteFromLeafNode(90)
	if !(node.leafHead == nil && !success) {
		t.FailNow()
	}
}

func TestInsertToTree_oneLeafNode_ok(t *testing.T) {
	tree := BPlusTree{
		order: 3,
		root:  initOneTestLeafNode(4),
	}
	tree.Insert(25, "")
	checkNode([]int64{10, 20, 25, 30, 40}, tree.root, t)
}

func TestInsertToTree_oneLeafNode_emptyLeafNode(t *testing.T) {
	tree := BPlusTree{
		order: 3,
		root:  initOneTestLeafNode(0),
	}
	tree.Insert(25, "")
	assertCondition := tree.root.countOfKeys == 1 &&
		tree.root.leafHead.value == 25 &&
		tree.root.leafHead.nextKey == nil &&
		tree.root.leafHead.previousKey == nil
	if !assertCondition {
		t.FailNow()
	}
}

func TestInsertToTree_rootNodeOverlow_ok(t *testing.T) {
	tree := BPlusTree{
		order: 3,
		root:  initOneTestLeafNode(5),
	}
	tree.Insert(25, "")
	rootLeftPointer := tree.root.internalNodeHead
	rootMiddleKey := tree.root.internalNodeHead.nextKey
	rootRightPointer := tree.root.internalNodeHead.nextKey.nextPointer
	assertCondition := tree.root.countOfKeys == 1 && //some checks
		!tree.root.isLeaf &&
		rootLeftPointer.previousKey == nil &&
		rootMiddleKey.previousPointer == rootLeftPointer &&
		rootRightPointer.previousKey == rootMiddleKey &&
		rootMiddleKey.value == 30
	checkNode([]int64{10, 20, 25}, rootLeftPointer.childNode, t)
	checkNode([]int64{30, 40, 50}, rootRightPointer.childNode, t)
	if !assertCondition {
		t.FailNow()
	}
}

func checkNode(keys []int64, node *BPlusTreeNode, t *testing.T) {
	currentKey := node.leafHead
	for index, key := range keys {
		assertCondition := currentKey.value == key
		if index > 0 && index < len(keys)-1 {
			assertCondition = assertCondition && currentKey.nextKey != nil && currentKey.previousKey != nil
		} else if index > 0 {
			assertCondition = assertCondition && currentKey.nextKey == nil
		} else {
			assertCondition = assertCondition && currentKey.previousKey == nil
		}
		if !assertCondition {
			t.FailNow()
		}
		currentKey = currentKey.nextKey
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
