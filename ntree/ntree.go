package ntree

import (
	"github.com/serdaroquai/importer/fifo"
)

type Tree struct {
	Root *Node
}
type Node struct {
	Id       string
	Parent   *Node
	Children []*Node
}

func (tree *Tree) BreadthFirst(doSomethingWith func(node *Node)) {
	q := fifo.NewFifoQueue()
	for node := tree.Root; node != nil; {
		doSomethingWith(node)
		for _, child := range node.Children {
			q.Add(child)
		}
		if newNode, ok := q.Get(); ok {
			node = newNode.(*Node)
		} else {
			node = (*Node)(nil)
		}
	}
}

func NewNode(id string) *Node {
	return &Node{Id: id}
}

func NewTree(rootNode *Node) *Tree {
	return &Tree{Root: rootNode}
}

func (node Node) IsRoot() bool {
	return node.Parent == nil
}

func (node *Node) AddChild(child *Node) *Node {
	child.Parent = node
	node.Children = append(node.Children, child)
	return node
}
