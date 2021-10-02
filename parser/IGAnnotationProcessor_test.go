package parser

import (
	"IG-Parser/tree"
	"fmt"
	"testing"
)

func TestExtractComponentLinkageNonAggregateLinkages(t *testing.T) {

	text := "A,p(Certified) A,p1(non-suspended) A,p1(previously reviewed) A1(Operator) or A,p2(recognized) A2(Handler) D(must not) I((produce [AND] trade))"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	// Parse statement, which should consider private suffices
	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing error for statement", text)
	}

	fmt.Println(s.String())

	if s.Attributes.Left.PrivateNodeLinks[0].Entry.(string) != "non-suspended" {
		t.Fatal("Did not correct identify private node value. Should have been 'non-suspended', but is", s.Attributes.Left.PrivateNodeLinks[0].Entry.(string))
	}

	if s.Attributes.Left.PrivateNodeLinks[1].Entry.(string) != "previously reviewed" {
		t.Fatal("Did not correct identify private node value. Should have been 'previously reviewed', but is", s.Attributes.Left.PrivateNodeLinks[1].Entry.(string))
	}

	if s.Attributes.Right.PrivateNodeLinks[0].Entry.(string) != "recognized" {
		t.Fatal("Did not correct identify private node value. Should have been 'recognized', but is", s.Attributes.Right.PrivateNodeLinks[0].Entry.(string))
	}

}

func TestExtractComponentLinkageAggregateLinkages(t *testing.T) {

	text := "A,p(Certified) A,p1(non-suspended) A,p1(previously reviewed) A1(Operator) or A,p2(recognized) A2(Handler) D(must not) I((produce [AND] trade))"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = true

	// Parse statement, which should consider private suffices
	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing error for statement", text)
	}

	fmt.Println(s.String())

	if s.Attributes.Left.PrivateNodeLinks[0].Entry.(string) != "non-suspended" {
		t.Fatal("Did not correct identify private node value. Should have been 'non-suspended', but is", s.Attributes.Left.PrivateNodeLinks[0].Entry.(string))
	}

	if s.Attributes.Left.PrivateNodeLinks[1].Entry.(string) != "previously reviewed" {
		t.Fatal("Did not correct identify private node value. Should have been 'previously reviewed', but is", s.Attributes.Left.PrivateNodeLinks[1].Entry.(string))
	}

	if s.Attributes.Right.PrivateNodeLinks[0].Entry.(string) != "recognized" {
		t.Fatal("Did not correct identify private node value. Should have been 'recognized', but is", s.Attributes.Right.PrivateNodeLinks[0].Entry.(string))
	}

}

/*
Tests matching of nodes in different trees based on suffix
 */
func TestFindNodesLinkedViaSuffix(t *testing.T) {

	one := tree.Node{Entry: "one"}
	two := tree.Node{Entry: "two", Suffix: "p1,55"}
	three := tree.Node{Entry: "three", Suffix: "2"}
	four := tree.Node{Entry: "four", Suffix: "p1"}

	oneTwo := tree.Node{LogicalOperator: "AND"}
	oneTwo.InsertLeftNode(&one)
	oneTwo.InsertRightNode(&two)

	threeFour := tree.Node{LogicalOperator: "OR"}
	threeFour.InsertLeftNode(&three)
	threeFour.InsertRightNode(&four)

	rootOne := tree.Node{}
	rootOne.InsertLeftNode(&oneTwo)
	rootOne.InsertRightNode(&threeFour)

	five := tree.Node{Entry: "five", Suffix: "p1,123"}
	six := tree.Node{Entry: "six"}
	seven := tree.Node{Entry: "seven", Suffix: "2"}
	eight := tree.Node{Entry: "eight", Suffix: "2,89"}

	fiveSix := tree.Node{LogicalOperator: "XOR"}
	fiveSix.InsertLeftNode(&five)
	fiveSix.InsertRightNode(&six)

	sevenEight := tree.Node{LogicalOperator: "AND"}
	sevenEight.InsertLeftNode(&seven)
	sevenEight.InsertRightNode(&eight)

	rootTwo := tree.Node{}
	rootTwo.InsertLeftNode(&fiveSix)
	rootTwo.InsertRightNode(&sevenEight)

	resMap := FindNodesLinkedViaSuffix(&rootOne, &rootTwo)

	// Check general results

	fmt.Println(resMap)

	if len(resMap) != 3 {
		t.Fatal("Results map has wrong number of entries.")
	}

	// Node two

	if len(resMap[&two]) != 1 {
		t.Fatal("Could not find entries for node two.")
	}

	if resMap[&two][0] != &five {
		t.Fatal("Wrong linkage. Node two should be linked to node five, but is", fmt.Sprint(resMap[&two][0]))
	}

	// Node three

	if len(resMap[&three]) != 2 {
		t.Fatal("Could not find entries for node three.")
	}

	if resMap[&three][0] != &seven {
		t.Fatal("Wrong linkage. Node three should be linked to node seven, but is", fmt.Sprint(resMap[&three][0]))
	}

	if resMap[&three][1] != &eight {
		t.Fatal("Wrong linkage. Node three should be linked to node eight, but is", fmt.Sprint(resMap[&three][1]))
	}

	// Node four

	if len(resMap[&four]) != 1 {
		t.Fatal("Could not find entries for node four.")
	}

	if resMap[&four][0] != &five {
		t.Fatal("Wrong linkage. Node four should be linked to node five, but is", fmt.Sprint(resMap[&four][0]))
	}
}

