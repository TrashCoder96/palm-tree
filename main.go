package main

import (
	"errors"
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
		currentLeaf := bptn.leafHead
		currentKeyValueMoreOrEqualsKey := false
		currentKeyIsLastKey := false
		for {
			currentKeyIsLastKey = currentLeaf.nextKey == nil
			currentKeyValueMoreOrEqualsKey = currentLeaf.value >= key
			if currentKeyValueMoreOrEqualsKey || currentKeyIsLastKey {
				//append after currentLeaf
				break
			} else {
				currentLeaf = currentLeaf.nextKey
			}
		}
		if currentKeyValueMoreOrEqualsKey && currentLeaf != bptn.leafHead {
			newLeaf.nextKey = currentLeaf
			currentLeaf.previousKey.nextKey = &newLeaf
			newLeaf.previousKey = currentLeaf.previousKey
			currentLeaf.previousKey = &newLeaf
		} else if currentLeaf == bptn.leafHead {
			newLeaf.nextKey = currentLeaf
			currentLeaf.previousKey = &newLeaf
			bptn.leafHead = &newLeaf
		} else if currentKeyIsLastKey {
			currentLeaf.nextKey = &newLeaf
			newLeaf.previousKey = currentLeaf
		} else {
			panic("Operation is not allowed")
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

func (bpt *BPlusTree) delete(key int64, node *BPlusTreeNode) error {
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
			bpt.redistributeNodesIfPossible(currentPointer.previousKey.previousPointer)
		} else {
			bpt.redistributeNodesIfPossible(currentPointer)
		}
		return nil
	}
	panic("Operation is not allowed")
}

func (bpt *BPlusTree) redistributeNodesIfPossible(subtree *bPlusTreePointer) {
	leftPointer := subtree
	middleKey := leftPointer.nextKey
	rightPointer := middleKey.nextPointer
	leftPointerChildNodeLessThanOrderMinusOne := leftPointer.childNode.countOfKeys <= bpt.order-1
	rightPointerChildNodeLessThanOrderMinusOne := rightPointer.childNode.countOfKeys <= bpt.order-1
	if leftPointer.childNode.isLeaf && rightPointer.childNode.isLeaf {
		if leftPointerChildNodeLessThanOrderMinusOne && rightPointerChildNodeLessThanOrderMinusOne {
			tailKey := leftPointer.childNode.leafHead
			for tailKey.nextKey != nil {
				tailKey = tailKey.nextKey
			}
			tailKey.nextKey = rightPointer.childNode.leafHead
			rightPointer.childNode.leafHead.previousKey = tailKey
			leftPointer.nextKey = leftPointer.nextKey.nextPointer.nextKey
		} else if leftPointerChildNodeLessThanOrderMinusOne {
			leftPointer.childNode.insertToLeafNode(rightPointer.childNode.leafHead.value, "")
			rightPointer.childNode.deleteFromLeafNode(rightPointer.childNode.leafHead.value)
			middleKey.value = rightPointer.childNode.leafHead.value
		} else if rightPointerChildNodeLessThanOrderMinusOne {
			rightPointer.childNode.insertToLeafNode(leftPointer.childNode.leafHead.value, "")
			leftPointer.childNode.deleteFromLeafNode(leftPointer.childNode.leafHead.value)
			middleKey.value = rightPointer.childNode.leafHead.value
		}
	} else {
		tailPointer := leftPointer.childNode.internalNodeHead
		for tailPointer.nextKey != nil {
			tailPointer = tailPointer.nextKey.nextPointer
		}
		if leftPointerChildNodeLessThanOrderMinusOne && rightPointerChildNodeLessThanOrderMinusOne {
			tailPointer.nextKey = middleKey
			middleKey.previousPointer = tailPointer
			middleKey.nextPointer = rightPointer.childNode.internalNodeHead
			rightPointer.childNode.internalNodeHead.previousKey = middleKey
			leftPointer.nextKey = leftPointer.nextKey.nextPointer.nextKey
		} else if leftPointerChildNodeLessThanOrderMinusOne {
			newKey := bPlusTreeKey{
				value:           middleKey.value,
				previousPointer: tailPointer,
				nextPointer:     rightPointer.childNode.internalNodeHead,
			}
			tailPointer.nextKey = &newKey
			rightPointer.childNode.internalNodeHead.previousKey = &newKey
			middleKey.value = rightPointer.childNode.internalNodeHead.nextKey.value
			rightPointer.childNode.internalNodeHead = rightPointer.childNode.internalNodeHead.nextKey.nextPointer
			rightPointer.childNode.internalNodeHead.previousKey = nil
			tailPointer.nextKey.nextPointer.nextKey = nil
		} else if rightPointerChildNodeLessThanOrderMinusOne {
			newKey := tailPointer.previousKey
			newKey.previousPointer.nextKey = nil
			newKey.nextPointer = rightPointer.childNode.internalNodeHead
			rightPointer.childNode.internalNodeHead.previousKey = newKey
			tailPointer.nextKey = newKey
			tailPointer.previousKey = nil
			rightPointer.childNode.internalNodeHead = tailPointer
			middleKey.value, rightPointer.childNode.internalNodeHead.nextKey.value = middleKey.value, rightPointer.childNode.internalNodeHead.nextKey.value
		}
	}
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
