package parser

import (
	"IG-Parser/core/tree"
	"testing"
)

func TestMapCombinationOverwrite(t *testing.T) {

	node1 := tree.Node{Entry: "Node1"}
	node2 := tree.Node{Entry: "Node2"}
	node3 := tree.Node{Entry: "Node3"}
	node4 := tree.Node{Entry: "Node4"}

	leftMap := make(map[*tree.Node]*tree.Node)
	leftMap[&node1] = &node3
	leftMap[&node2] = &node1
	leftMap[&node3] = &node2

	rightMap := make(map[*tree.Node]*tree.Node)
	rightMap[&node1] = &node2
	rightMap[&node4] = &node1

	mapRes := CombineMaps(leftMap, rightMap, true)
	if len(mapRes) != 4 {
		t.Fatal("Map length should be 4, but is", len(mapRes))
	}

	if mapRes[&node1] != &node2 {
		t.Fatal("Entry should have been &node2, but is", mapRes[&node1])
	}

}

func TestMapCombinationNoOverwrite(t *testing.T) {

	node1 := tree.Node{Entry: "Node1"}
	node2 := tree.Node{Entry: "Node2"}
	node3 := tree.Node{Entry: "Node3"}
	node4 := tree.Node{Entry: "Node4"}

	leftMap := make(map[*tree.Node]*tree.Node)
	leftMap[&node1] = &node3
	leftMap[&node2] = &node1
	leftMap[&node3] = &node2

	rightMap := make(map[*tree.Node]*tree.Node)
	rightMap[&node1] = &node2
	rightMap[&node4] = &node1

	mapRes := CombineMaps(leftMap, rightMap, false)
	if len(mapRes) != 4 {
		t.Fatal("Map length should be 4, but is", len(mapRes))
	}

	if mapRes[&node1] != &node3 {
		t.Fatal("Entry should have been &node3, but is", mapRes[&node1])
	}

}
