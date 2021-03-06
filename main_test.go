package main

import (
	"testing"
)

func TestGetPointer_ok(t *testing.T) {
	node := initOneTestInternalNodeByKeys([]int64{6, 16})
	pointer := node.getPointer(6)
	assertCondition := pointer.previousKey.value == 6 &&
		pointer.nextKey.value == 16
	if !assertCondition {
		t.FailNow()
	}
}

func TestGetPointer_KeyLessThan10(t *testing.T) {
	node := initOneTestInternalNode(9, 10, 10)
	pointer := node.getPointer(3)
	assertCondition := pointer.previousKey == nil &&
		pointer.nextKey.value == 10
	if !assertCondition {
		t.FailNow()
	}
}

func TestGetPointer_KeyMoreThan90(t *testing.T) {
	node := initOneTestInternalNode(9, 10, 10)
	pointer := node.getPointer(100)
	assertCondition := pointer.previousKey.value == 90 &&
		pointer.nextKey == nil
	if !assertCondition {
		t.FailNow()
	}
}

func TestInsertToLeafNode_ok(t *testing.T) {
	countOfKeys := 9
	node := initOneTestLeafNode(countOfKeys, 10, 10)
	node.insertToLeafNode(25, "")
	checkLeafNode([]int64{10, 20, 25, 30, 40, 50, 60, 70, 80, 90}, node, t)
}

func TestInsertToLeafNode_atStartPosition(t *testing.T) {
	countOfKeys := 9
	node := initOneTestLeafNode(countOfKeys, 10, 10)
	node.insertToLeafNode(5, "")
	checkLeafNode([]int64{5, 10, 20, 30, 40, 50, 60, 70, 80, 90}, node, t)
}

func TestInsertToLeafNode_atEndPosition(t *testing.T) {
	countOfKeys := 9
	node := initOneTestLeafNode(countOfKeys, 10, 10)
	node.insertToLeafNode(100, "")
	checkLeafNode([]int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}, node, t)
}

func TestInsertToLeafNode_emptyList(t *testing.T) {
	node := initOneTestLeafNode(0, 0, 10)
	node.insertToLeafNode(100, "")
	if !(node.leafHead != nil && node.leafHead.value == 100 && node.leafHead.nextKey == nil) {
		t.FailNow()
	}
}

func TestDeleteFromLeafNode_ok(t *testing.T) {
	node := initOneTestLeafNode(9, 10, 10)
	node.deleteFromLeafNode(40)
	checkLeafNode([]int64{10, 20, 30, 50, 60, 70, 80, 90}, node, t)
}

func TestDeleteFromLeafNode_atStartPosition(t *testing.T) {
	node := initOneTestLeafNode(9, 10, 10)
	node.deleteFromLeafNode(10)
	checkLeafNode([]int64{20, 30, 40, 50, 60, 70, 80, 90}, node, t)
}

func TestDeleteFromLeafNode_atEndPosition(t *testing.T) {
	node := initOneTestLeafNode(9, 10, 10)
	node.deleteFromLeafNode(90)
	checkLeafNode([]int64{10, 20, 30, 40, 50, 60, 70, 80}, node, t)
}

func TestDeleteFromLeafNode_emptyList(t *testing.T) {
	node := initOneTestLeafNode(0, 10, 10)
	success := node.deleteFromLeafNode(90)
	if !(node.leafHead == nil && !success) {
		t.FailNow()
	}
}

func TestInsertToTree_oneLeafNode_ok(t *testing.T) {
	tree := BPlusTree{
		order: 5,
		root:  initOneTestLeafNodeByKeys([]int64{7, 29, 32, 33, 43, 45, 67}),
	}
	tree.Insert(52, "")
	checkLeafNode([]int64{7, 29, 32, 33, 43, 45, 52, 67}, tree.root, t)
}

func TestInsertToTree_oneLeafNode_emptyLeafNode(t *testing.T) {
	tree := BPlusTree{
		order: 3,
		root:  initOneTestLeafNode(0, 10, 10),
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
		root:  initOneTestLeafNode(5, 10, 10),
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
	checkLeafNode([]int64{10, 20, 25}, rootLeftPointer.childNode, t)
	checkLeafNode([]int64{30, 40, 50}, rootRightPointer.childNode, t)
	if !assertCondition {
		t.FailNow()
	}
}

func TestInsertToTree_leafNodeOverflow_ok(t *testing.T) {
	tree := BPlusTree{
		order: 2,
		root:  initOneTestInternalNode(2, 10, 10),
	}
	rootMiddlePointer := tree.root.internalNodeHead.nextKey.nextPointer
	rootMiddlePointer.childNode = initOneTestLeafNode(3, 12, 2)
	tree.Insert(13, "")
	checkInternalNode([]int64{10, 14, 20}, tree.root, t)
	checkLeafNode([]int64{12, 13}, tree.root.internalNodeHead.nextKey.nextPointer.childNode, t)
	checkLeafNode([]int64{14, 16}, tree.root.internalNodeHead.nextKey.nextPointer.nextKey.nextPointer.childNode, t)
}

func TestInsertToTree_internalNodeOverflow_ok(t *testing.T) {
	tree := BPlusTree{
		order: 2,
		root:  initOneTestInternalNode(1, 20, 10),
	}
	tree.root.internalNodeHead.childNode = initOneTestInternalNode(2, 16, 3)
	tree.root.internalNodeHead.childNode.internalNodeHead.childNode = initOneTestInternalNode(3, 1, 1)
	tailPointer := tree.root.internalNodeHead.childNode.internalNodeHead.childNode.internalNodeHead
	for tailPointer.nextKey != nil {
		tailPointer = tailPointer.nextKey.nextPointer
	}
	tailPointer.childNode = initOneTestLeafNode(3, 5, 1)
	tree.Insert(8, "")
	checkInternalNode([]int64{20}, tree.root, t)
	checkInternalNode([]int64{2, 16, 19}, tree.root.internalNodeHead.childNode, t)
	checkInternalNode([]int64{1}, tree.root.internalNodeHead.childNode.internalNodeHead.childNode, t)
	checkInternalNode([]int64{3, 7}, tree.root.internalNodeHead.childNode.internalNodeHead.nextKey.nextPointer.childNode, t)
	checkLeafNode([]int64{5, 6}, tree.root.internalNodeHead.
		childNode.internalNodeHead.nextKey.nextPointer.childNode.
		internalNodeHead.nextKey.nextPointer.childNode, t)
	checkLeafNode([]int64{7, 8}, tree.root.internalNodeHead.
		childNode.internalNodeHead.nextKey.nextPointer.childNode.
		internalNodeHead.nextKey.nextPointer.nextKey.nextPointer.childNode, t)
}

func TestDeleteFromTree_onlyLeafNode_ok(t *testing.T) {
	tree := BPlusTree{
		order: 2,
		root:  initOneTestLeafNode(3, 10, 10),
	}
	assertCondition := tree.Delete(20)
	if !assertCondition {
		t.FailNow()
	}
	checkLeafNode([]int64{10, 30}, tree.root, t)
}

func TestDeleteFromTree_onlyLeafNode_fail(t *testing.T) {
	tree := BPlusTree{
		order: 2,
		root:  initOneTestLeafNode(3, 10, 10),
	}
	assertCondition := !tree.Delete(25)
	if !assertCondition {
		t.FailNow()
	}
}

func TestDeleteFromTree_rootAndFirstTwoLeafNode_ok(t *testing.T) {
	tree := BPlusTree{
		order: 2,
		root:  initOneTestInternalNode(3, 20, 10),
	}
	tree.root.internalNodeHead.childNode = initOneTestLeafNode(3, 3, 3)
	tree.root.internalNodeHead.nextKey.nextPointer.childNode = initOneTestLeafNode(2, 21, 1)
	assertCondition := tree.Delete(3)
	if !assertCondition {
		t.FailNow()
	}
	checkInternalNode([]int64{20, 30, 40}, tree.root, t)
	checkLeafNode([]int64{6, 9}, tree.root.internalNodeHead.childNode, t)
	checkLeafNode([]int64{21, 22}, tree.root.internalNodeHead.nextKey.nextPointer.childNode, t)
}

func TestDeleteFromTree_rootAndLastTwoLeafNode_ok(t *testing.T) {
	tree := BPlusTree{
		order: 2,
		root:  initOneTestInternalNode(3, 20, 10),
	}
	tailPointer := tree.root.internalNodeHead
	for tailPointer.nextKey != nil {
		tailPointer = tailPointer.nextKey.nextPointer
	}
	tailPointer.childNode = initOneTestLeafNode(3, 41, 1)
	tailPointer.previousKey.previousPointer.childNode = initOneTestLeafNode(2, 31, 1)
	assertCondition := tree.Delete(41)
	if !assertCondition {
		t.FailNow()
	}
	checkInternalNode([]int64{20, 30, 40}, tree.root, t)
	checkLeafNode([]int64{31, 32}, tailPointer.previousKey.previousPointer.childNode, t)
	checkLeafNode([]int64{42, 43}, tailPointer.childNode, t)
}

func TestDeleteFromTree_rootAndLastLittleTwoNodes_ok(t *testing.T) {
	tree := BPlusTree{
		order: 2,
		root:  initOneTestInternalNode(3, 20, 10),
	}
	tailPointer := tree.root.internalNodeHead
	for tailPointer.nextKey != nil {
		tailPointer = tailPointer.nextKey.nextPointer
	}
	tailPointer.childNode = initOneTestLeafNode(2, 41, 1)
	tailPointer.previousKey.previousPointer.childNode = initOneTestLeafNode(1, 31, 0)
	assertCondition := tree.Delete(42)
	if !assertCondition {
		t.FailNow()
	}
	checkInternalNode([]int64{20, 30}, tree.root, t)
	checkLeafNode([]int64{31, 41}, tree.root.internalNodeHead.nextKey.nextPointer.nextKey.nextPointer.childNode, t)
}

func TestDeleteFromTree_moveToLeftNode_ok(t *testing.T) {
	tree := BPlusTree{
		order: 3,
		root:  initOneTestInternalNode(3, 20, 20),
	}
	tree.root.internalNodeHead.nextKey.nextPointer.childNode = initOneTestLeafNode(3, 21, 1)
	tree.root.internalNodeHead.nextKey.nextPointer.nextKey.nextPointer.childNode = initOneTestLeafNode(5, 41, 1)
	assertCondition := tree.Delete(22)
	if !assertCondition {
		t.FailNow()
	}
	checkInternalNode([]int64{20, 42, 60}, tree.root, t)
	checkLeafNode([]int64{21, 23, 41}, tree.root.internalNodeHead.nextKey.nextPointer.childNode, t)
	n := tree.root.internalNodeHead.nextKey.nextPointer.nextKey.nextPointer.childNode
	checkLeafNode([]int64{42, 43, 44, 45}, n, t)
}

func TestDeleteFromTree_moveToRightNode_ok(t *testing.T) {
	tree := BPlusTree{
		order: 3,
		root:  initOneTestInternalNode(3, 20, 20),
	}
	tailPointer := tree.root.internalNodeHead
	for tailPointer.nextKey != nil {
		tailPointer = tailPointer.nextKey.nextPointer
	}
	tailPointer.previousKey.previousPointer.childNode = initOneTestLeafNode(5, 41, 1)
	tailPointer.childNode = initOneTestLeafNode(3, 61, 1)
	assertCondition := tree.Delete(62)
	if !assertCondition {
		t.FailNow()
	}
	checkInternalNode([]int64{20, 40, 45}, tree.root, t)
	checkLeafNode([]int64{41, 42, 43, 44}, tailPointer.previousKey.previousPointer.childNode, t)
	checkLeafNode([]int64{45, 61, 63}, tailPointer.childNode, t)
}

func TestDeleteFromTree_mergeToOneRootNode_ok(t *testing.T) {
	tree := BPlusTree{
		order: 5,
		root:  initOneTestInternalNode(1, 20, 0),
	}
	tree.root.internalNodeHead.childNode = initOneTestLeafNode(4, 10, 1)
	tree.root.internalNodeHead.nextKey.nextPointer.childNode = initOneTestLeafNode(4, 21, 1)
	assertCondition := tree.Delete(21)
	if !assertCondition {
		t.FailNow()
	}
	checkLeafNode([]int64{10, 11, 12, 13, 22, 23, 24}, tree.root, t)
}

func TestDeleteFromTree_mergeTwoLastNodes_ok(t *testing.T) {
	tree := BPlusTree{
		order: 5,
		root:  initOneTestInternalNode(3, 10, 10),
	}
	tailPointer := tree.root.internalNodeHead
	for tailPointer.nextKey != nil {
		tailPointer = tailPointer.nextKey.nextPointer
	}
	tailPointer.previousKey.previousPointer.childNode = initOneTestLeafNode(4, 21, 1)
	tailPointer.previousKey.previousPointer.previousKey.previousPointer.childNode = initOneTestLeafNode(5, 11, 1)
	assertCondition := tree.Delete(15)
	if !assertCondition {
		t.FailNow()
	}
	checkInternalNode([]int64{10, 30}, tree.root, t)
	checkLeafNode([]int64{11, 12, 13, 14, 21, 22, 23, 24}, tree.root.internalNodeHead.nextKey.nextPointer.childNode, t)
}

func TestDeleteFromTree_mergeTwoInternalNodes_ok(t *testing.T) {
	tree := BPlusTree{
		order: 5,
		root:  initOneTestInternalNodeByKeys([]int64{21}),
	}
	tree.root.internalNodeHead.childNode = initOneTestInternalNodeByKeys([]int64{5, 6, 7, 8, 15})
	tree.root.internalNodeHead.nextKey.nextPointer.childNode = initOneTestInternalNodeByKeys([]int64{25, 26, 27, 28})
	leftChildNodeTailPointer := tree.root.internalNodeHead.childNode.internalNodeHead
	for leftChildNodeTailPointer.nextKey != nil {
		leftChildNodeTailPointer = leftChildNodeTailPointer.nextKey.nextPointer
	}
	leftChildNodeTailPointer.childNode = initOneTestLeafNodeByKeys([]int64{16, 17, 18, 19, 20})
	leftChildNodeTailPointer.previousKey.previousPointer.childNode = initOneTestLeafNodeByKeys([]int64{9, 10, 11, 12})
	assertCondition := tree.Delete(16)
	if !assertCondition {
		t.FailNow()
	}
	checkInternalNode([]int64{5, 6, 7, 8, 21, 25, 26, 27, 28}, tree.root, t)
	necessaryPointer := tree.root.internalNodeHead
	for !(necessaryPointer.nextKey != nil && necessaryPointer.nextKey.value == 21) {
		necessaryPointer = necessaryPointer.nextKey.nextPointer
	}
	checkLeafNode([]int64{9, 10, 11, 12, 17, 18, 19, 20}, necessaryPointer.childNode, t)
}

func TestDeleteFromTree_moveToRightForInternalNodes_ok(t *testing.T) {
	tree := BPlusTree{
		order: 5,
		root:  initOneTestInternalNodeByKeys([]int64{21}),
	}
	tree.root.internalNodeHead.childNode = initOneTestInternalNodeByKeys([]int64{10, 11, 12, 13, 14, 15})
	tree.root.internalNodeHead.nextKey.nextPointer.childNode = initOneTestInternalNodeByKeys([]int64{30, 36, 37, 38})
	tree.root.internalNodeHead.nextKey.nextPointer.childNode.internalNodeHead.childNode = initOneTestLeafNodeByKeys([]int64{22, 23, 24, 25, 26})
	tree.root.internalNodeHead.nextKey.nextPointer.childNode.internalNodeHead.nextKey.nextPointer.childNode = initOneTestLeafNodeByKeys([]int64{31, 32, 33, 34})
	assertCondition := tree.Delete(22)
	if !assertCondition {
		t.FailNow()
	}
	checkInternalNode([]int64{15}, tree.root, t)
	checkInternalNode([]int64{10, 11, 12, 13, 14}, tree.root.internalNodeHead.childNode, t)
	checkInternalNode([]int64{21, 36, 37, 38}, tree.root.internalNodeHead.nextKey.nextPointer.childNode, t)
	checkLeafNode([]int64{23, 24, 25, 26, 31, 32, 33, 34}, tree.root.internalNodeHead.nextKey.nextPointer.childNode.internalNodeHead.nextKey.nextPointer.childNode, t)
}

func TestDeleteFromTree_moveToLeftForInternalNodes_ok(t *testing.T) {
	tree := BPlusTree{
		order: 5,
		root:  initOneTestInternalNodeByKeys([]int64{21}),
	}
	tree.root.internalNodeHead.childNode = initOneTestInternalNodeByKeys([]int64{6, 7, 8, 15})
	tree.root.internalNodeHead.nextKey.nextPointer.childNode = initOneTestInternalNodeByKeys([]int64{25, 26, 27, 28, 29})
	leftChildNodeTailPointer := tree.root.internalNodeHead.childNode.internalNodeHead
	for leftChildNodeTailPointer.nextKey != nil {
		leftChildNodeTailPointer = leftChildNodeTailPointer.nextKey.nextPointer
	}
	leftChildNodeTailPointer.childNode = initOneTestLeafNodeByKeys([]int64{16, 17, 18, 19, 20})
	leftChildNodeTailPointer.previousKey.previousPointer.childNode = initOneTestLeafNodeByKeys([]int64{9, 10, 11, 12})
	assertCondition := tree.Delete(16)
	if !assertCondition {
		t.FailNow()
	}
	checkInternalNode([]int64{25}, tree.root, t)
	checkInternalNode([]int64{6, 7, 8, 21}, tree.root.internalNodeHead.childNode, t)
	checkInternalNode([]int64{26, 27, 28, 29}, tree.root.internalNodeHead.nextKey.nextPointer.childNode, t)
	necessaryPointer := tree.root.internalNodeHead.childNode.internalNodeHead
	for !(necessaryPointer.nextKey != nil && necessaryPointer.nextKey.value == 21) {
		necessaryPointer = necessaryPointer.nextKey.nextPointer
	}
	checkLeafNode([]int64{9, 10, 11, 12, 17, 18, 19, 20}, necessaryPointer.childNode, t)
}

func checkLeafNode(keys []int64, node *BPlusTreeNode, t *testing.T) {
	currentKey := node.leafHead
	actualCountOfKeysInNode := 0
	for index, key := range keys {
		assertCondition := currentKey.value == key
		if index > 0 && index < len(keys)-1 {
			assertCondition = assertCondition &&
				currentKey.nextKey != nil &&
				currentKey.previousKey != nil
		} else if index > 0 {
			assertCondition = assertCondition &&
				currentKey.nextKey == nil
		} else {
			assertCondition = assertCondition &&
				currentKey.previousKey == nil
		}
		if !assertCondition {
			t.FailNow()
		}
		actualCountOfKeysInNode = actualCountOfKeysInNode + 1
		currentKey = currentKey.nextKey
	}
	if node.countOfKeys != actualCountOfKeysInNode {
		t.FailNow()
	}
}

func checkInternalNode(keys []int64, node *BPlusTreeNode, t *testing.T) {
	currentPointer := node.internalNodeHead
	actualCountOfKeysInNode := 0
	for _, key := range keys {
		actualCountOfKeysInNode++
		currentPointer = currentPointer.nextKey.nextPointer
		assertCondition := currentPointer.previousKey.value == key
		if !assertCondition {
			t.FailNow()
		}
	}
	if node.countOfKeys != actualCountOfKeysInNode {
		t.FailNow()
	}
	tailPointer := node.internalNodeHead
	for tailPointer.nextKey != nil {
		tailPointer = tailPointer.nextKey.nextPointer
	}
	currentPointer = tailPointer
	for _, key := range reverse(keys) {
		currentPointer = currentPointer.previousKey.previousPointer
		assertCondition := currentPointer.nextKey.value == key
		if !assertCondition {
			t.FailNow()
		}
	}
}

func reverse(slice []int64) []int64 {
	len := len(slice)
	for i := 0; i < len/2; i++ {
		slice[i], slice[len-i-1] = slice[len-i-1], slice[i]
	}
	return slice
}

func initOneTestLeafNode(countOfKeys int, base int64, step int64) *BPlusTreeNode {
	value := int64(base)
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
		value = value + step
		newKey := bPlusTreeKey{
			value:       value,
			previousKey: previousKey,
		}
		previousKey.nextKey = &newKey
		previousKey = previousKey.nextKey
	}
	return &node
}

func initOneTestLeafNodeByKeys(keys []int64) *BPlusTreeNode {
	node := BPlusTreeNode{
		isLeaf:      true,
		countOfKeys: len(keys),
	}
	node.leafHead = &bPlusTreeKey{
		value: keys[0],
	}
	previousKey := node.leafHead
	for _, key := range keys[1:] {
		newKey := bPlusTreeKey{
			value:       key,
			previousKey: previousKey,
		}
		previousKey.nextKey = &newKey
		previousKey = previousKey.nextKey
	}
	return &node
}

func initOneTestInternalNodeByKeys(keys []int64) *BPlusTreeNode {
	node := BPlusTreeNode{
		isLeaf:      false,
		countOfKeys: len(keys),
	}
	node.internalNodeHead = &bPlusTreePointer{}
	previousPointer := node.internalNodeHead
	for _, key := range keys {
		newKey := bPlusTreeKey{
			previousPointer: previousPointer,
			value:           key,
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

func initOneTestInternalNode(countOfKeys int, base int64, step int64) *BPlusTreeNode {
	node := BPlusTreeNode{
		isLeaf:      false,
		countOfKeys: countOfKeys,
	}
	value := int64(base)
	node.internalNodeHead = &bPlusTreePointer{}
	previousPointer := node.internalNodeHead
	for i := 0; i < countOfKeys; i++ {
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
		value = value + step
	}
	return &node
}
