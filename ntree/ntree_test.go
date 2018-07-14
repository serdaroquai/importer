package ntree

import "testing"

func TestAddChild(t *testing.T) {
	node := []*Node{NewNode("0"), NewNode("1"), NewNode("2")}

	node[0].AddChild(node[1]).AddChild(node[2])

	if node[0].Children[0] != node[1] || node[0].Children[1] != node[2] {
		t.Fail()
	} else if node[1].Parent != node[0] || node[2].Parent != node[0] {
		t.Fail()
	}
}

func TestIsRoot(t *testing.T) {
	node := NewNode("2")
	if !node.IsRoot() {
		t.Errorf("Node expected to be a root node")
	}

	NewNode("1").AddChild(node)
	if node.IsRoot() {
		t.Errorf("Node expected not to be a root node")
	}
}

func TestBreadthFirst(t *testing.T) {
	
	node := []*Node{NewNode("0"), NewNode("1"), NewNode("2"), NewNode("3"), NewNode("4"), NewNode("5")}
	tree := NewTree(node[0])

	var result string
	
	tree.BreadthFirst(func(node *Node) {
		result = result + node.Id
	})
	if (result != "0") {
		t.Errorf("Expected %v, but have %v", "0", result)
	}

	// now reset result and add some real data
	result = ""
	node[0].AddChild(node[1]).AddChild(node[2])
	node[1].AddChild(node[3]).AddChild(node[4])
	node[2].AddChild(node[5])
	// 	   0
	//   1   2
	//  3 4   5

	tree.BreadthFirst(func(node *Node) {
		result = result + node.Id
	})

	if (result != "012345") {
		t.Errorf("Expected %v, but have %v", "012345", result)
	}
}
