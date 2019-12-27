package main

import (
	"os"
)

func main() {
	process(os.Args)
}

func process(params []string) {
	tree := initTree(5)
	for i := 1; i <= 100000; i++ {
		tree.Insert(int64(i), "")
	}
	//tree.Find(991)
	for i := 50000; i > 0; i-- {
		tree.Delete(int64(i))
	}
}

func initTree(order int) *BPlusTree {
	tree := BPlusTree{
		order: order,
		root: &BPlusTreeNode{
			isLeaf:      true,
			countOfKeys: 0,
		},
	}
	return &tree
}

//BPlusTree struct
type BPlusTree struct {
	order int
	root  *BPlusTreeNode
}

//BPlusTreeNode struct
type BPlusTreeNode struct {
	countOfKeys      int
	isLeaf           bool
	internalNodeHead *bPlusTreePointer //only for internal node
	leafHead         *bPlusTreeKey     //only for leaf node
}

type bPlusTreeKey struct {
	value           int64
	nextPointer     *bPlusTreePointer
	nextKey         *bPlusTreeKey
	previousKey     *bPlusTreeKey
	previousPointer *bPlusTreePointer
}

type bPlusTreePointer struct {
	nextKey     *bPlusTreeKey
	previousKey *bPlusTreeKey
	childNode   *BPlusTreeNode
}

//Insert function
func (bpt *BPlusTree) Insert(key int64, value string) {
	if bpt.root.leafHead == nil && bpt.root.internalNodeHead == nil {
		bpt.root.leafHead = &bPlusTreeKey{
			value: key,
		}
		bpt.root.countOfKeys = 1
	} else {
		rootPointer := bPlusTreePointer{childNode: bpt.root}
		bpt.insert(key, value, &rootPointer)
		if rootPointer.nextKey != nil {
			newNode := BPlusTreeNode{
				internalNodeHead: &rootPointer,
				countOfKeys:      1,
			}
			bpt.root = &newNode
		}
	}
}

//Delete function
func (bpt *BPlusTree) Delete(key int64) bool {
	if bpt.delete(key, bpt.root) {
		if bpt.root.internalNodeHead != nil && bpt.root.internalNodeHead.nextKey == nil {
			bpt.root = bpt.root.internalNodeHead.childNode
		}
		return true
	}
	return false
}

//Find function
func (bpt *BPlusTree) Find(key int64) bool {
	return bpt.find(key, bpt.root)
}

//function returns true, if cutting has occured
func (bpt *BPlusTree) insert(key int64, value string, pointerToNode *bPlusTreePointer) bool {
	if pointerToNode.childNode.isLeaf {
		pointerToNode.childNode.insertToLeafNode(key, value)
		if pointerToNode.childNode.countOfKeys > 2*bpt.order-1 {
			cutIfPossible(pointerToNode)
			return true
		}
	} else {
		suitablePointer := pointerToNode.childNode.getPointer(key)
		if bpt.insert(key, value, suitablePointer) {
			pointerToNode.childNode.countOfKeys = pointerToNode.childNode.countOfKeys + 1
		}
		if pointerToNode.childNode.countOfKeys > 2*bpt.order-1 {
			cutIfPossible(pointerToNode)
			return true
		}
	}
	return false
}

func (bptn *BPlusTreeNode) getPointer(key int64) *bPlusTreePointer {
	if !bptn.isLeaf {
		currentPointer := bptn.internalNodeHead
		nextKeyValueMoreThanKey := false
		nextKeyIsNil := false
		for {
			nextKeyIsNil = currentPointer.nextKey == nil
			if !nextKeyIsNil {
				nextKeyValueMoreThanKey = currentPointer.nextKey.value > key
			}
			if nextKeyValueMoreThanKey || nextKeyIsNil {
				break
			} else {
				currentPointer = currentPointer.nextKey.nextPointer
			}
		}
		if nextKeyIsNil && !nextKeyValueMoreThanKey {
			return currentPointer
		} else if nextKeyValueMoreThanKey {
			return currentPointer
		} else {
			panic("Operation is not allowed!!!")
		}
	} else {
		panic("Operation is not allowed!!!")
	}
}

func (bptn *BPlusTreeNode) insertToLeafNode(key int64, value string) {
	newLeaf := bPlusTreeKey{value: key}
	if bptn.leafHead == nil {
		bptn.leafHead = &newLeaf
	} else {
		foundNextLeaf := false
		currentLeaf := bptn.leafHead
		for i := 0; i < bptn.countOfKeys-1; i++ {
			if currentLeaf.value < key {
				currentLeaf = currentLeaf.nextKey
			} else {
				foundNextLeaf = true
				break
			}
		}
		newLeaf = bPlusTreeKey{
			value: key,
		}
		if foundNextLeaf {
			newLeaf.nextKey = currentLeaf
			if currentLeaf.previousKey == nil {
				currentLeaf.previousKey = &newLeaf
				bptn.leafHead = &newLeaf
			} else {
				newLeaf.previousKey = currentLeaf.previousKey
				currentLeaf.previousKey.nextKey = &newLeaf
				currentLeaf.previousKey = &newLeaf
			}
		} else {
			currentLeaf.nextKey = &newLeaf
			newLeaf.previousKey = currentLeaf
		}
	}
	bptn.countOfKeys = bptn.countOfKeys + 1
}

func cutIfPossible(pointer *bPlusTreePointer) {
	leftPointer := pointer
	newKey := bPlusTreeKey{
		previousPointer: leftPointer,
	}
	rightPointer := bPlusTreePointer{
		previousKey: &newKey,
		nextKey:     pointer.nextKey,
	}
	newKey.nextPointer = &rightPointer
	if pointer.nextKey != nil {
		pointer.nextKey.previousPointer = &rightPointer
	}
	leftPointer.nextKey = &newKey
	leftNode := pointer.childNode
	if pointer.childNode.isLeaf {
		rightNode := BPlusTreeNode{
			isLeaf:      true,
			countOfKeys: leftNode.countOfKeys / 2,
		}
		keyBeforeNextNode := leftNode.leafHead
		for i := 1; i < leftNode.countOfKeys/2; i++ {
			keyBeforeNextNode = keyBeforeNextNode.nextKey
		}
		leftNode.countOfKeys = rightNode.countOfKeys
		rightPointer.childNode = &rightNode
		newKey.value = keyBeforeNextNode.nextKey.value
		rightNode.leafHead = keyBeforeNextNode.nextKey
		rightNode.leafHead.previousKey = nil
		keyBeforeNextNode.nextKey = nil
	} else {
		rightNode := BPlusTreeNode{
			isLeaf:      false,
			countOfKeys: leftNode.countOfKeys / 2,
		}
		pointerBeforeMiddleKey := leftNode.internalNodeHead
		for i := 1; i < leftNode.countOfKeys/2; i++ {
			pointerBeforeMiddleKey = pointerBeforeMiddleKey.nextKey.nextPointer
		}
		leftNode.countOfKeys = rightNode.countOfKeys - 1
		rightPointer.childNode = &rightNode
		rightNode.internalNodeHead = pointerBeforeMiddleKey.nextKey.nextPointer
		rightNode.internalNodeHead.previousKey = nil
		newKey.value = pointerBeforeMiddleKey.nextKey.value
		pointerBeforeMiddleKey.nextKey = nil
	}
}

func (bpt *BPlusTree) find(key int64, node *BPlusTreeNode) bool {
	if node != nil {
		if node.isLeaf {
			leaf := node.leafHead
			for leaf != nil && leaf.value != key {
				leaf = leaf.nextKey
			}
			if leaf == nil {
				return false
			}
			return true
		}
		pointer := node.getPointer(key)
		return bpt.find(key, pointer.childNode)
	}
	panic("Operation is not allowed")
}

func (bpt *BPlusTree) delete(key int64, node *BPlusTreeNode) bool {
	if node != nil {
		if node.isLeaf {
			return node.deleteFromLeafNode(key)
		}
		currentPointer := node.getPointer(key)
		success := bpt.delete(key, currentPointer.childNode)
		if success {
			if currentPointer.nextKey == nil {
				bpt.redistributeNodesIfPossible(currentPointer.previousKey.previousPointer, node)
			} else {
				bpt.redistributeNodesIfPossible(currentPointer, node)
			}
		}
		return success
	}
	panic("Operation is not allowed")
}

func (bpt *BPlusTree) redistributeNodesIfPossible(subtree *bPlusTreePointer, node *BPlusTreeNode) {
	leftPointer := subtree
	middleKey := leftPointer.nextKey
	rightPointer := middleKey.nextPointer
	//leftPointerChildNodeLessThanOrder := leftPointer.childNode.countOfKeys <= bpt.order
	leftPointerChildNodeLessThanOrderMinusOne := leftPointer.childNode.countOfKeys <= bpt.order-1
	//rightPointerChildNodeLessThanOrder := rightPointer.childNode.countOfKeys <= bpt.order
	rightPointerChildNodeLessThanOrderMinusOne := rightPointer.childNode.countOfKeys <= bpt.order-1
	if leftPointerChildNodeLessThanOrderMinusOne && rightPointerChildNodeLessThanOrderMinusOne {
		merge(subtree)
		node.countOfKeys = node.countOfKeys - 1
	} else if leftPointerChildNodeLessThanOrderMinusOne && !rightPointerChildNodeLessThanOrderMinusOne {
		moveToLeftNode(subtree)
	} else if !leftPointerChildNodeLessThanOrderMinusOne && rightPointerChildNodeLessThanOrderMinusOne {
		moveToRightNode(subtree)
	}
}

func moveToLeftNode(subtree *bPlusTreePointer) {
	leftPointer := subtree
	middleKey := leftPointer.nextKey
	rightPointer := middleKey.nextPointer
	lowLevelIsLeaves := leftPointer.childNode.isLeaf && rightPointer.childNode.isLeaf
	if lowLevelIsLeaves {
		movedItem := rightPointer.childNode.leafHead
		rightPointer.childNode.leafHead = rightPointer.childNode.leafHead.nextKey
		rightPointer.childNode.leafHead.previousKey = nil
		movedItem.nextKey = nil
		tailKey := leftPointer.childNode.leafHead
		for tailKey.nextKey != nil {
			tailKey = tailKey.nextKey
		}
		tailKey.nextKey = movedItem
		movedItem.previousKey = tailKey
		middleKey.value = rightPointer.childNode.leafHead.value
	} else {
		tailPointer := leftPointer.childNode.internalNodeHead
		for tailPointer.nextKey != nil {
			tailPointer = tailPointer.nextKey.nextPointer
		}
		newKey := bPlusTreeKey{
			value:           middleKey.value,
			nextPointer:     rightPointer.childNode.internalNodeHead,
			previousPointer: tailPointer,
		}
		tailPointer.nextKey = &newKey
		rightPointer.childNode.internalNodeHead.previousKey = &newKey
		middleKey.value = rightPointer.childNode.internalNodeHead.nextKey.value
		rightPointer.childNode.internalNodeHead = rightPointer.childNode.internalNodeHead.nextKey.nextPointer
		newKey.nextPointer.nextKey = nil
	}
	leftPointer.childNode.countOfKeys = leftPointer.childNode.countOfKeys + 1
	rightPointer.childNode.countOfKeys = rightPointer.childNode.countOfKeys - 1
}

func moveToRightNode(subtree *bPlusTreePointer) {
	leftPointer := subtree
	middleKey := leftPointer.nextKey
	rightPointer := middleKey.nextPointer
	lowLevelIsLeaves := leftPointer.childNode.isLeaf && rightPointer.childNode.isLeaf
	if lowLevelIsLeaves {
		tailKey := leftPointer.childNode.leafHead
		for tailKey.nextKey != nil {
			tailKey = tailKey.nextKey
		}
		tailKey.previousKey.nextKey = nil
		tailKey.previousKey = nil
		tailKey.nextKey = rightPointer.childNode.leafHead
		rightPointer.childNode.leafHead.previousKey = tailKey
		rightPointer.childNode.leafHead = tailKey
		middleKey.value = rightPointer.childNode.leafHead.value
	} else {
		tailPointer := leftPointer.childNode.internalNodeHead
		for tailPointer.nextKey != nil {
			tailPointer = tailPointer.nextKey.nextPointer
		}
		newKey := bPlusTreeKey{
			value:           middleKey.value,
			nextPointer:     rightPointer.childNode.internalNodeHead,
			previousPointer: tailPointer,
		}
		middleKey.value = tailPointer.previousKey.value
		tailPointer.nextKey = &newKey
		tailPointer.previousKey.previousPointer.nextKey = nil
		tailPointer.previousKey = nil
		rightPointer.childNode.internalNodeHead.previousKey = &newKey
		rightPointer.childNode.internalNodeHead = newKey.previousPointer
	}
	leftPointer.childNode.countOfKeys = leftPointer.childNode.countOfKeys - 1
	rightPointer.childNode.countOfKeys = rightPointer.childNode.countOfKeys + 1
}

func merge(subtree *bPlusTreePointer) {
	leftPointer := subtree
	middleKey := leftPointer.nextKey
	rightPointer := middleKey.nextPointer
	lowLevelIsLeaves := leftPointer.childNode.isLeaf && rightPointer.childNode.isLeaf
	leftPointer.childNode.countOfKeys += rightPointer.childNode.countOfKeys
	if lowLevelIsLeaves {
		tailLeaf := leftPointer.childNode.leafHead
		for tailLeaf.nextKey != nil {
			tailLeaf = tailLeaf.nextKey
		}
		tailLeaf.nextKey = rightPointer.childNode.leafHead
		rightPointer.childNode.leafHead.previousKey = tailLeaf
		leftPointer.nextKey = rightPointer.nextKey
		if rightPointer.nextKey != nil {
			rightPointer.nextKey.previousPointer = rightPointer
		}
	} else {
		tailPointer := leftPointer.childNode.internalNodeHead
		for tailPointer.nextKey != nil {
			tailPointer = tailPointer.nextKey.nextPointer
		}
		tailPointer.nextKey = middleKey
		middleKey.previousPointer = tailPointer
		middleKey.nextPointer = rightPointer.childNode.internalNodeHead
		rightPointer.childNode.internalNodeHead.previousKey = middleKey
		leftPointer.nextKey = rightPointer.nextKey
		if rightPointer.nextKey != nil {
			rightPointer.nextKey.previousPointer = rightPointer
		}
		leftPointer.childNode.countOfKeys++
	}
}

func (bptn *BPlusTreeNode) deleteFromLeafNode(key int64) bool {
	currentLeaf := bptn.leafHead
	for currentLeaf != nil {
		if currentLeaf.value == key {
			if currentLeaf.previousKey == nil {
				bptn.leafHead = currentLeaf.nextKey
				if currentLeaf.nextKey != nil {
					currentLeaf.nextKey.previousKey = nil
				}
			} else if currentLeaf.nextKey == nil {
				currentLeaf.previousKey.nextKey = nil
			} else {
				currentLeaf.previousKey.nextKey = currentLeaf.nextKey
				currentLeaf.nextKey.previousKey = currentLeaf.previousKey
			}
			bptn.countOfKeys = bptn.countOfKeys - 1
			return true
		}
		currentLeaf = currentLeaf.nextKey
	}
	return false
}
