package main

import (
	"log"
	"os"
)

func main() {
	process(os.Args)
}

func process(params []string) {
}

func initTree(degree int) *BPlusTree {
	tree := BPlusTree{degree: degree}
	return &tree
}

//BPlusTree struct
type BPlusTree struct {
	degree int
	root   *BPlusTreeNode
}

//BPlusTreeNode struct
type BPlusTreeNode struct {
	count            int
	isLeaf           bool
	internalNodeHead *bPlusTreePointer //only for internal node
	leafHead         *bPlusTreeKey     //only for leaf node
}

func (bptn *BPlusTreeNode) insertToLeafNode(key int64, value string) {
	newLeaf := bPlusTreeKey{value: key}
	currentLeaf := bptn.leafHead
	var previousLeaf *bPlusTreeKey
	for currentLeaf != nil || currentLeaf.value <= key {
		previousLeaf = currentLeaf
		currentLeaf = currentLeaf.nextKey
	}
	newLeaf.nextKey = currentLeaf
	previousLeaf.nextKey = &newLeaf
	bptn.count = bptn.count + 1
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

func (bptn *BPlusTreeNode) cutTailWithMiddleKey() (tail *BPlusTreeNode, middleKey *bPlusTreeKey) {
	tail = &BPlusTreeNode{
		count: bptn.count / 2,
	}
	middleKey = &bPlusTreeKey{}
	if bptn.isLeaf {
		//devide by two leaf nodes
		currentKey := bptn.leafHead
		var previousKey *bPlusTreeKey
		for i := 1; i < tail.count; i++ {
			previousKey = currentKey
			currentKey = currentKey.nextKey
		}
		tail.isLeaf = true
		tail.leafHead = currentKey
		previousKey.nextKey = nil
		middleKey.value = currentKey.value
	} else {
		//devide by two internal nodes
		currentKey := bptn.internalNodeHead.nextKey
		var previousKey *bPlusTreeKey
		for i := 1; i < tail.count; i++ {
			previousKey = currentKey
			currentKey = currentKey.nextPointer.nextKey
		}
		middleKey = currentKey
		tail.internalNodeHead = currentKey.nextPointer
		currentKey.nextPointer = nil
		previousKey.nextPointer.nextKey = nil
	}
	return
}

type bPlusTreeKey struct {
	value       int64
	nextPointer *bPlusTreePointer
	nextKey     *bPlusTreeKey
}

type bPlusTreePointer struct {
	nextKey   *bPlusTreeKey
	childNode *BPlusTreeNode
}

//Find function
func (bpt *BPlusTree) Find(key int64) {
}

//Insert function
func (bpt *BPlusTree) Insert(key int64, value string) {
	if bpt.root == nil {
		bpt.root = &BPlusTreeNode{
			isLeaf: true,
			leafHead: &bPlusTreeKey{
				value: key,
			},
		}
		return
	}
	bpt.insert(key, value, bpt.root)
}

func (bpt *BPlusTree) insert(key int64, value string, node *BPlusTreeNode) {
	if node != nil {
		if node.isLeaf {
			node.insertToLeafNode(key, value)
			if node.count > bpt.degree {
				//tail := node.cutTailWithMiddleKey()

			}
		} else {
			bpt.insert(key, value, node)
		}
	} else {
		log.Panicln("Panic! Operation not allowed. Tree node is nil")
	}
}

//Update function
func (bpt *BPlusTree) Update(key int64, value string) {
}

//Delete function
func (bpt *BPlusTree) Delete(key int64) {
}
