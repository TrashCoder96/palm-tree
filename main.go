package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	process(os.Args)
}

func process(params []string) {
}

func initTree(order int) *BPlusTree {
	tree := BPlusTree{order: order}
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

func (bptn *BPlusTreeNode) getPointer(key int64) *bPlusTreePointer {
	if !bptn.isLeaf {
		currentPointer := bptn.internalNodeHead
		nextKeyValueMoreOrEqualsKey := false
		nextKeyIsNil := false
		for {
			nextKeyIsNil = currentPointer.nextKey == nil
			if !nextKeyIsNil {
				nextKeyValueMoreOrEqualsKey = currentPointer.nextKey.value >= key
			}
			if nextKeyValueMoreOrEqualsKey || nextKeyIsNil {
				break
			} else {
				currentPointer = currentPointer.nextKey.nextPointer
			}
		}
		if nextKeyIsNil && !nextKeyValueMoreOrEqualsKey {
			return currentPointer
		} else if nextKeyValueMoreOrEqualsKey {
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

func (bptn *BPlusTreeNode) cutByTwoNodes() *bPlusTreePointer {
	tail := &BPlusTreeNode{}
	middleKey := &bPlusTreeKey{}
	var currentKey *bPlusTreeKey
	if bptn.isLeaf {
		//devide by two leaf nodes
		currentKey = bptn.leafHead
		for i := 0; i < bptn.countOfKeys/2; i++ {
			currentKey = currentKey.nextKey
		}
		tail.isLeaf = true
		tail.leafHead = currentKey
		currentKey.previousKey.nextKey = nil
		currentKey.previousKey = nil
		middleKey.value = currentKey.value
		countInNode := bptn.countOfKeys / 2
		bptn.countOfKeys = countInNode
		tail.countOfKeys = countInNode
	} else {
		//devide by two internal nodes
		currentKey = bptn.internalNodeHead.nextKey
		for i := 0; i < bptn.countOfKeys/2; i++ {
			currentKey = currentKey.nextPointer.nextKey
		}
		tail.internalNodeHead = &bPlusTreePointer{
			nextKey:   currentKey.nextPointer.nextKey,
			childNode: currentKey.nextPointer.childNode,
		}
		middleKey.value = currentKey.value
		currentKey.previousPointer.nextKey = nil
		currentKey.nextPointer.previousKey = nil
		countInNode := bptn.countOfKeys / 2
		bptn.countOfKeys = countInNode
		tail.countOfKeys = countInNode - 1
	}
	leftPointer := &bPlusTreePointer{
		childNode: bptn,
		nextKey:   middleKey,
	}
	rightPointer := &bPlusTreePointer{
		childNode:   tail,
		previousKey: middleKey,
	}
	middleKey.nextPointer = rightPointer
	middleKey.previousPointer = leftPointer
	return leftPointer
}

func (bpt *BPlusTree) insert(key int64, value string, node *BPlusTreeNode) *bPlusTreePointer {
	if node != nil {
		if node.isLeaf {
			node.insertToLeafNode(key, value)
			if node.countOfKeys > 2*bpt.order-1 {
				subtree := node.cutByTwoNodes()
				return subtree
			}
			return nil
		}
		currentPointer := node.getPointer(key)
		//if internal node
		subtree := bpt.insert(key, value, currentPointer.childNode)
		if subtree != nil {
			subtreeRightPointer := subtree.nextKey.nextPointer
			subtreeLeftPointer := subtree
			if currentPointer == node.internalNodeHead {
				node.internalNodeHead = subtreeLeftPointer
				subtreeRightPointer.nextKey = currentPointer.nextKey
				subtreeRightPointer.nextKey.previousPointer = subtreeRightPointer
			} else if currentPointer.nextKey == nil {
				currentPointer.previousKey.nextPointer = subtreeLeftPointer
				subtreeLeftPointer.previousKey = currentPointer.previousKey
			} else {
				currentPointer.previousKey.nextPointer = subtreeLeftPointer
				subtreeLeftPointer.previousKey = currentPointer.previousKey
				subtreeRightPointer.nextKey = currentPointer.nextKey
				currentPointer.nextKey.previousPointer = subtreeRightPointer
			}
			node.countOfKeys = node.countOfKeys + 1
			if node.countOfKeys > 2*bpt.order-1 {
				subtree := node.cutByTwoNodes()
				return subtree
			}
		}
		return nil
	}
	panic("Operation is not allowed")
}

//Insert function
func (bpt *BPlusTree) Insert(key int64, value string) {
	if bpt.root == nil {
		bpt.root = &BPlusTreeNode{
			isLeaf:      true,
			countOfKeys: 1,
			leafHead: &bPlusTreeKey{
				value: key,
			},
		}
		return
	}
	if subtree := bpt.insert(key, value, bpt.root); subtree != nil {
		newNode := BPlusTreeNode{
			internalNodeHead: subtree,
			countOfKeys:      1,
		}
		bpt.root = &newNode
	}
}

//Delete function
func (bpt *BPlusTree) Delete(key int64) {
	bpt.delete(key, bpt.root)
}

//Find function
func (bpt *BPlusTree) Find(key int64) bool {
	return bpt.find(key, bpt.root)
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

func (bpt *BPlusTree) delete(key int64, node *BPlusTreeNode) error {
	if key == 8449 {
		fmt.Println(key)
	}
	if node != nil {
		if node.isLeaf {
			if ok := node.deleteFromLeafNode(key); ok {
				return nil
			}
			return errors.New("Key not found")
		}
		currentPointer := node.getPointer(key)
		if itemNotFoundErr := bpt.delete(key, currentPointer.childNode); itemNotFoundErr != nil {
			return itemNotFoundErr
		}
		if currentPointer.nextKey == nil {
			bpt.redistributeNodesIfPossible(currentPointer.previousKey.previousPointer, node)
		} else {
			bpt.redistributeNodesIfPossible(currentPointer, node)
		}
		return nil
	}
	panic("Operation is not allowed")
}

func (bpt *BPlusTree) redistributeNodesIfPossible(subtree *bPlusTreePointer, node *BPlusTreeNode) {
	leftPointer := subtree
	middleKey := leftPointer.nextKey
	rightPointer := middleKey.nextPointer
	leftPointerChildNodeLessThanOrderMinusOne := leftPointer.childNode.countOfKeys <= bpt.order-1   //t - 1
	rightPointerChildNodeLessThanOrderMinusOne := rightPointer.childNode.countOfKeys <= bpt.order-1 //t - 1
	if leftPointerChildNodeLessThanOrderMinusOne && rightPointerChildNodeLessThanOrderMinusOne {
		merge(subtree)
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
	}
	leftPointer.childNode.countOfKeys = leftPointer.childNode.countOfKeys - 1
	rightPointer.childNode.countOfKeys = rightPointer.childNode.countOfKeys + 1
}

func merge(subtree *bPlusTreePointer) {
	leftPointer := subtree
	middleKey := leftPointer.nextKey
	rightPointer := middleKey.nextPointer
	lowLevelIsLeaves := leftPointer.childNode.isLeaf && rightPointer.childNode.isLeaf
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
		middleKey.nextPointer = leftPointer.childNode.internalNodeHead
		leftPointer.childNode.internalNodeHead.previousKey = middleKey
		leftPointer.nextKey = rightPointer.nextKey
		if rightPointer.nextKey != nil {
			rightPointer.nextKey.previousPointer = rightPointer
		}
	}
	leftPointer.childNode.countOfKeys += rightPointer.childNode.countOfKeys
}

func (bptn *BPlusTreeNode) deleteFromLeafNode(key int64) bool {
	currentLeaf := bptn.leafHead
	for currentLeaf != nil {
		if currentLeaf.value == key {
			if currentLeaf.previousKey == nil {
				bptn.leafHead = currentLeaf.nextKey
				currentLeaf.nextKey.previousKey = nil
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
