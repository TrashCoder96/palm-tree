package main

import (
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

//PrintTree struct
func (bpt *BPlusTree) PrintTree() {

}

//BPlusTreeNode struct
type BPlusTreeNode struct {
	countOfKeys      int
	isLeaf           bool
	internalNodeHead *bPlusTreePointer //only for internal node
	leafHead         *bPlusTreeKey     //only for leaf node
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

func (bptn *BPlusTreeNode) deleteFromLeafNode(key int64) bool {
	currentLeaf := bptn.leafHead
	var previousLeaf *bPlusTreeKey
	for currentLeaf != nil {
		if currentLeaf.value == key {
			if previousLeaf != nil {
				previousLeaf.nextKey = currentLeaf.nextKey.nextKey
			} else {
				bptn.leafHead = currentLeaf.nextKey.nextKey
			}
			return true
		}
		previousLeaf = currentLeaf
		currentLeaf = currentLeaf.nextKey
	}
	return false
}
