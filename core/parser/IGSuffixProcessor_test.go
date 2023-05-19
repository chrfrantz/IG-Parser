package parser

import (
	"IG-Parser/core/tree"
	"fmt"
	"testing"
)

func TestExtractComponentLinkageNonAggregateLinkages(t *testing.T) {

	text := "A,p(Certified) A1,p(non-suspended) A1,p(previously reviewed) A1(Operator) or A2,p(recognized) A2(Handler) D(must not) I((produce [AND] trade))"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	// Parse statement, which should consider private suffices
	stmt, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing error for statement", text)
	}

	s := stmt[0].Entry.(*tree.Statement)

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

	text := "A,p(Certified) A1,p1(non-suspended) A1,p1(previously reviewed) A1(Operator) or A2,p(recognized) A2(Handler) D(must not) I((produce [AND] trade))"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = true

	// Parse statement, which should consider private suffices
	stmt, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing error for statement", text)
	}

	s := stmt[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	if s.AttributesPropertySimple.Entry.(string) != "Certified" {
		t.Fatal("Did not correctly parse shared property.")
	}

	if s.AttributesPropertySimple.Parent != nil {
		t.Fatal("Did not correctly clean parent reference of shared property.")
	}

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
	nine := tree.Node{Entry: "nine"}

	fiveSix := tree.Node{LogicalOperator: "XOR"}
	fiveSix.InsertLeftNode(&five)
	fiveSix.InsertRightNode(&six)

	sevenEight := tree.Node{LogicalOperator: "AND"}
	sevenEight.InsertLeftNode(&seven)
	sevenEight.InsertRightNode(&eight)
	sevenEight.InsertRightNode(&nine)

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

/*
Tests the correct parsing of a component-pair-nested private property
(Pattern: 'P1,p{ A( ...) { I( ...) Bdir( ...) [XOR] I( ... ) Bdir( ... ) } ... }')
alongside a primitive property (Pattern: 'P1,p(something)')
*/
func TestNestedPrivatePropertyPrimitiveAndNestedComponentPair(t *testing.T) {

	// Problem area is Bdir,p and Bdir (parentheses)
	text := "Such E(notification) M(shall) F(provide): (1) A P(description of each noncompliance); " +
		"(2) The P(facts upon which the notification of noncompliance is based); and " +
		"(3) The P1(date) P1,p{by which the A(certified operation) D(must) " +
		"{I(rebut [XOR] correct) Bdir,p(each) Bdir(noncompliance) [AND] " +
		"I(submit) Bdir,p(supporting) Bdir(documentation) of Bdir,p(each such correction) Cac(when correction is possible)}} P1,p(private component)"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should not have returned error, but returned error ", err)
	}

	stmt := s[0].Entry.(*tree.Statement)

	if stmt.ConstitutedEntity.Entry != "notification" {
		t.Fatal("Inconsistency when parsing Constituted Entity. Should be 'notification'.")
	}

	if stmt.Modal.Entry != "shall" {
		t.Fatal("Inconsistency when parsing Modal. Should be 'shall'.")
	}

	if stmt.ConstitutiveFunction.Entry != "provide" {
		t.Fatal("Inconsistency when parsing Constitutive Function. Should be 'provide'.")
	}

	if stmt.ConstitutingProperties.Left.Left.Entry != "description of each noncompliance" {
		t.Fatal("Inconsistency when parsing left left Constituting Property. Should be 'description of each noncompliance'.")
	}

	if stmt.ConstitutingProperties.Left.Right.Entry != "facts upon which the notification of noncompliance is based" {
		t.Fatal("Inconsistency when parsing left right Constituting Property. Should be 'facts upon which the notification of noncompliance is based'.")
	}

	if stmt.ConstitutingProperties.Right.Entry != "date" {
		t.Fatal("Inconsistency when parsing right Constituting Property. Should be 'date'.")
	}

	if stmt.ConstitutingProperties.Right.Suffix != "1" {
		t.Fatal("Inconsistency when parsing right Constituting Property suffix. Should be '1'.")
	}

	// Private nodes

	// primitive node

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[0].Entry != "private component" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'private component'.")
	}

	// complex node

	// second one (first complex one)

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Attributes.Entry != "certified operation" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'certified operation'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry != "must" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'must'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.LogicalOperator != "XOR" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'XOR'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Left.Entry != "rebut" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'rebut'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Right.Entry != "correct" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'correct'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry != "noncompliance" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'noncompliance'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObjectPropertySimple.Entry != "each" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'each'.")
	}

	// third statement (second complex one)

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Attributes.Entry != "certified operation" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'certified operation'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry != "must" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'must'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry != "submit" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'submit'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry != "documentation" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'documentation'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObjectPropertySimple.LogicalOperator != "bAND" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'bAND'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObjectPropertySimple.Left.Entry != "supporting" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'supporting'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObjectPropertySimple.Right.Entry != "each such correction" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'each such correction'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionSimple.Entry != "when correction is possible" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'when correction is possible'.")
	}

}

/*
Tests the correct parsing of a component-pair-nested private property (Pattern: 'P1,p{ A( ...) { I( ...) Bdir( ...) [XOR] I( ... ) Bdir( ... ) } ... }')
alongside a primitive property (Pattern: 'P1,p(something)') and a single nested property (Pattern: 'P1,p{ E( ...) ... }').
*/
func TestNestedPrivatePropertyPrimitiveAndNestedComponentPairAndAdditionalSingleNestedProperty(t *testing.T) {

	// Problem area is Bdir,p and Bdir (parentheses)
	text := "Such E(notification) M(shall) F(provide): (1) A P(description of each noncompliance); " +
		"(2) The P(facts upon which the notification of noncompliance is based); and " +
		"(3) The P1(date) P1,p{by which the A(certified operation) D(must) " +
		"{I(rebut [XOR] correct) Bdir,p(each) Bdir(noncompliance) [AND] " +
		"I(submit) Bdir,p(supporting) Bdir(documentation) of Bdir,p(each such correction) Cac(when correction is possible)}} P1,p(private component) " +
		// additional nested property
		"P1,p{where E(date) F(is defined) in the P(Gregorian calendar)}."

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should not have returned error, but returned error ", err)
	}

	stmt := s[0].Entry.(*tree.Statement)

	if stmt.ConstitutedEntity.Entry != "notification" {
		t.Fatal("Inconsistency when parsing Constituted Entity. Should be 'notification'.")
	}

	if stmt.Modal.Entry != "shall" {
		t.Fatal("Inconsistency when parsing Modal. Should be 'shall'.")
	}

	if stmt.ConstitutiveFunction.Entry != "provide" {
		t.Fatal("Inconsistency when parsing Constitutive Function. Should be 'provide'.")
	}

	if stmt.ConstitutingProperties.Left.Left.Entry != "description of each noncompliance" {
		t.Fatal("Inconsistency when parsing left left Constituting Property. Should be 'description of each noncompliance'.")
	}

	if stmt.ConstitutingProperties.Left.Right.Entry != "facts upon which the notification of noncompliance is based" {
		t.Fatal("Inconsistency when parsing left right Constituting Property. Should be 'facts upon which the notification of noncompliance is based'.")
	}

	if stmt.ConstitutingProperties.Right.Entry != "date" {
		t.Fatal("Inconsistency when parsing right Constituting Property. Should be 'date'.")
	}

	if stmt.ConstitutingProperties.Right.Suffix != "1" {
		t.Fatal("Inconsistency when parsing right Constituting Property suffix. Should be '1'.")
	}

	// Private nodes

	// primitive node

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[0].Entry != "private component" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'private component'.")
	}

	// complex node

	// second one (first complex one)

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Attributes.Entry != "certified operation" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'certified operation'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry != "must" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'must'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.LogicalOperator != "XOR" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'XOR'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Left.Entry != "rebut" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'rebut'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Right.Entry != "correct" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'correct'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry != "noncompliance" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'noncompliance'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[1].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObjectPropertySimple.Entry != "each" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'each'.")
	}

	// third statement (second complex one)

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Attributes.Entry != "certified operation" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'certified operation'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry != "must" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'must'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry != "submit" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'submit'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry != "documentation" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'documentation'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObjectPropertySimple.LogicalOperator != "bAND" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'bAND'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObjectPropertySimple.Left.Entry != "supporting" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'supporting'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObjectPropertySimple.Right.Entry != "each such correction" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'each such correction'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[2].Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionSimple.Entry != "when correction is possible" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'when correction is possible'.")
	}

	// fourth statement (nested one)

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[3].Entry.(*tree.Statement).ConstitutedEntity.Entry != "date" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'date'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[3].Entry.(*tree.Statement).ConstitutiveFunction.Entry != "is defined" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'is defined'.")
	}

	if stmt.ConstitutingProperties.Right.PrivateNodeLinks[3].Entry.(*tree.Statement).ConstitutingProperties.Entry != "Gregorian calendar" {
		t.Fatal("Inconsistency when parsing private node of right-most Constituting Property. Should be 'Gregorian calendar'.")
	}

}
