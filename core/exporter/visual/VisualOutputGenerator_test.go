package visual

import (
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/core/parser"
	"IG-Parser/core/tree"
	"fmt"
	"os"
	"testing"
)

/*
This file contains tests related to visual output generation.
*/

/*
Tests the generation of basic tree output for visual output. Suppresses shared
elements in output.
*/
func TestVisualOutputBasicWithoutSharedElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualTreeOutputBasicWithoutSharedElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the generation of basic tree output for visual output. Includes shared
elements in output.
*/
func TestVisualOutputBasicWithSharedElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Activate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(true)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualTreeOutputBasicWithSharedElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the generation of tree output with nested properties for visual output.
*/
func TestVisualOutputNestedProperties(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Bdir,p{E(operation) F(has been vetted) Cex(before certification)} " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + output.String())

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualTreeOutputNestedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the generation of tree output for visual output, including annotations.
*/
func TestVisualOutputAnnotations(t *testing.T) {

	text := "A[gov=enforcer,anim=animate](National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=monitor](inspect and), I[act=enforce](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir[gov=monitored,anim=animate](approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex[ref=(Act,part)](for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E[gov=enforcer,anim=animate](Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A[gov=monitor,anim=animate](NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Activate annotations
	tabular.SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualTreeOutputAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests proper output of proper linkages of complex private properties alongside shared properties for visual output.
Tests flat output for properties (labels, as opposed to tree structure).
*/
func TestVisualOutputPropertyNodesFlatPrinting(t *testing.T) {
	text := "A(General Manager) A,p(shared quality) A1(Region Manager) A1,p(left quality) A1,p(right quality) A1,p(third quality)"

	// Activate annotations
	tabular.SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualTreeComplexPrivateNodesFlat.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper flat output of combined shared properties for visual output.
Tests flat output for properties (labels, as opposed to tree structure).
*/
func TestVisualOutputSharedPropertyNodesFlatPrinting(t *testing.T) {
	text := "The A(Program Manager) D(may) I(initiate) Bdir,p((suspension [XOR] revocation)) Bdir(proceedings)"

	// Activate annotations
	tabular.SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualSharedPropertyNodesFlatPrinting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper output of proper linkages of complex private properties alongside shared properties for visual output.
Tests tree structure output for properties.
*/
func TestVisualOutputPropertyNodesTreePrinting(t *testing.T) {
	text := "A(General Manager) A,p(shared quality) A1(Region Manager) A1,p(left quality) A1,p(right quality) A1,p(third quality)"

	// Activate annotations
	tabular.SetIncludeAnnotations(true)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Activate binary tree printing
	tree.SetBinaryPrinting(true)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualTreeComplexPrivateNodesTree.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests the generation of basic tree output for visual output, but as non-binary tree (i.e., collapsing entries associated with same logical operator for given component).
*/
func TestVisualOutputBasicNonBinaryTree(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualTreeOutputBasicNonBinary.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the generation of complex tree output for visual output as non-binary tree (i.e., collapsing entries associated with same logical operator for given component).
Does not decompose property trees.
*/
func TestVisualOutputComplexNonBinaryTree(t *testing.T) {

	// Complex entry
	text := "The Congress finds and declares that it is the E(national policy) F([is] to (encourage [AND] assist)) the P(states) Cex{ A(states) I(to exercise) Cex(effectively) their Bdir(responsibilities) Bdir,p(in the coastal zone) Cex(through the (development [AND] implementation) of management programs to achieve wise use of the (land [AND] water) resources of the coastal zone, giving full consideration to (ecological [AND] cultural [AND] historic [AND] esthetic) values as well as the needs for compatible economic development), Cex{which E(programs) M(should) Cex(at least) F(provide for)— (A) the P1(protection) P1,p1(of natural resources, including (wetlands [AND] floodplains [AND] estuaries [AND] beaches [AND] dunes [AND] barrier islands [AND] coral reefs [AND] fish and wildlife and their habitat) within the coastal zone), the P2(management) P2,p2((of coastal development to minimize the loss of (life [AND] property) caused by improper development in (flood-prone [AND] storm surge [AND] geological hazard [AND] erosion-prone) areas [AND] in areas likely to be (affected by [OR] vulnerable to) (sea level rise [AND] land subsidence [AND] saltwater intrusion) [AND] by the destruction of natural protective features such as (beaches [AND] dunes [AND] wetlands [AND] barrier islands))), (C) the P3(management) P3,p(of coastal development to (improve [AND] safeguard [AND] restore) the quality of coastal waters, [AND] to protect (natural resources [AND] existing uses of those waters)), (D) P4,p1(priority) P4(consideration) P4,p2(being given to (coastal-dependent (uses [AND] orderly processes) for siting major facilities related to (national defense [AND] energy [AND] fisheries development [AND] recreation [AND] ports [AND] transportation), [AND] the location to the maximum extent practicable of new (commercial [AND] industrial) developments (in [XOR] adjacent) to areas where such development already exists)), (E) P5,p1(public) P5(access) P5,p2(to the coasts for recreation purposes), (F) P6(assistance) P6,p(in the redevelopment of (deteriorating urban (waterfronts [AND] ports) [AND] sensitive (preservation [AND] restoration) of (historic [AND] cultural [AND] esthetic) coastal features)), (G) P7(the (coordination [AND] simplification) of procedures) P7,p1(in order to ensure expedited governmental decision making for the management of coastal resources), (H) P8((continued (consultation [AND] coordination) with, [AND] the giving of adequate consideration to the views of affected Federal agencies)), (I) P9(the giving of ((timely [AND] effective) notification of , [AND] opportunities for (public [AND] local) government participation in coastal management decision making)), (J) P10(assistance) P10,p1(to support comprehensive (planning [AND] conservation [AND] management) for living marine resources) P10,p1,p1(including planning for (the siting of (pollution control [AND] aquaculture facilities) within the coastal zone [AND] improved coordination between ((State [AND] Federal) coastal zone management agencies [AND] (State [AND] wildlife) agencies))) }}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualTreeOutputComplexNonBinary.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests multiple combinations embedded in second level phrase/side of first-order statement.
Includes shared elements in output.
*/
func TestVisualOutputMultiLevelEmbeddedCombinationsWithSharedElements(t *testing.T) {

	// Entry containing multiple combinations on right side of first-order combination
	text := "Cex(( left1 [XOR] shared (left [AND] right) via (left2 [XOR] right2)))"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Activate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(true)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualMultiCombinationsPhraseWithSharedElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests multiple combinations embedded in second level phrase/side of first-order statement.
Suppresses shared elements in output.
*/
func TestVisualOutputMultiLevelEmbeddedCombinationsWithoutSharedElements(t *testing.T) {

	// Entry containing multiple combinations on right side of first-order combination
	text := "Cex(( left1 [XOR] shared (left [AND] right) via (left2 [XOR] right2)))"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Activate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualMultiCombinationsPhraseWithoutSharedElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for the default example statement that showcases most IG features.
*/
func TestVisualOutputDefaultExample(t *testing.T) {

	// Default example
	text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that Cac{Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

	// Activate annotations
	tabular.SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualDefaultExample.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for the default example statement that showcases most IG features, but moving activation conditions to the top.
*/
func TestVisualOutputDefaultExampleActivationConditionsFirst(t *testing.T) {

	// Default example
	text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that Cac{Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

	// Activate annotations
	tabular.SetIncludeAnnotations(true)
	// Deactivate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Move activation conditions to the top
	tree.SetMoveActivationConditionsToFront(true)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualDefaultExampleActivationConditionsFirst.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests escaping of symbols (e.g., quotation marks) and internal parentheses for visual output.
*/
func TestVisualOutputEscapingSymbols(t *testing.T) {

	// Statement with quotation marks and internal parentheses
	text := "The E(corporation) M(shall) F(be) P(a \"Type B\" corporation) Cex(pursuant to Section 201(b) of the New York State Not-for-Profit Corporation Law.)"

	// Activate annotations
	tabular.SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualEscapingSymbols.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests special symbols parsing for visual output.
*/
func TestVisualOutputSpecialSymbols(t *testing.T) {

	// Statement with quotation marks and internal parentheses
	text := "The E(cor#po$ration) M(sh<all) F(b>e) P[1%25](a \"Type B\" cor=poration) Cex[#<=>27.14](pur.suant to Se:ct!ion 201(b) of the N;ew York St,ate £Not-for-Profit €Corporatio$n Law.)"

	// Activate annotations
	tabular.SetIncludeAnnotations(true)
	// Activate flat printing
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualSpecialSymbols.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests linear multi-level nesting in visual output (i.e., Cac{Cac{Cac{}}}).
*/
func TestVisualOutputLinearMultiLevelNesting(t *testing.T) {

	// Statement with multiple levels of linear nesting (i.e., no combinations)
	text := "A,p(First) A(Actor) I(action1) I(action2) Bdir{A(actor2) I(actionLevel2) Cac{A(actor3) I(actionLevel3) Bdir(some object)}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualLinearMultilevelNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests complex multi-level nesting in visual output (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestVisualOutputComplexMultiLevelNesting(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualComplexMultilevelNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests embedding of properties (,p) in nested statements embedded in combinations.
*/
func TestVisualOutputComponentNestedStatementCombinationsWithProperties(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations that contain properties
	text := "Cac{Cac{A,p(Resident) A(Program Manager) I((suspects [OR] establishes)) Bdir(violations)} [AND] Cac{E(Program Manager) F(is authorized) for the P,p(relevant) P(region)}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualComponentNestedStatementCombinationsWithProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests moderately complex complete statement featuring nested activation condition combinations, properties in statements, as well as Or else combinations.
*/
func TestVisualOutputModeratelyComplexStatementWithNestedCombinationsPropertiesAndOrElse(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations on activation condition, or else, and includes properties
	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), D(must) I(inspect), I((review [AND] (revise [AND] resubmit))) Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) Cex(for compliance with the (Act or [XOR] regulations in this part)) if Cac{Cac{A(Program Manager) I((suspects [OR] establishes)) Bdir(violations)} [AND] Cac{E(Program Manager) F(is authorized) for the P,p(relevant) P(region)}}, or else O{O{A,p(Manager's) A(supervisor) D(may) I((suspend [XOR] revoke)) Bdir,p(Program Manager's) Bdir(authority)} [XOR] O{A(regional board) D(may) I((warn [OR] fine)) Bdir,p(violating) Bdir(Program Manager)}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualModeratelyComplexStatementWithNestedCombinationsPropertiesAndOrElse.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests 2nd-order nested activation condition combinations that include properties in all nested statements.
*/
func TestVisualOutput2ndOrderNestedStatementCombinationsWithProperties(t *testing.T) {

	// Statement with 2nd-order nesting with statement combinations on activation conditions, and that includes properties in all nested statements
	text := "Cac{Cac{A(actor1) I(act1) Bdir,p(prop1) Bdir(bdir1) Cex(cex1)}  [OR] Cac{A(actor2) I(act2) Bdir,p(prop2) Bdir(bdir2) Cex(cex2)}} [AND] Cac{A(actor2) I(act2) Bdir,p(prop2) Bdir(bdir2) Cex(cex2)}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisual2ndOrderNestedStatementCombinationsWithProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests higher-order statement combinations ({Cac{ ... } [AND] Cac{ ... } [XOR] {Cac{ ... } [AND] Cac{ ... }}}) or component-level nesting combinations in visual output
using an extended version of the default example statement.
NOTE: Nesting works until seven levels at this stage
*/
func TestVisualOutputHigherOrderStatementNestedComponentCombinationsDefaultExample(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that Cac{{Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] {Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)} [OR] Cac{A(actor5) I(act5)}}} [XOR] Cac{A(actor5) I(act5)}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualHigherOrderStatementNestedComponentCombinationsDefaultStatement.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests higher-order statement combinations ({Cac{ ... } [AND] Cac{ ... } [XOR] {Cac{ ... } [AND] Cac{ ... }}}) or component-level nesting combinations in visual output
using an extended version of the default example statement.
Tests the statement with complexity output.
NOTE: Nesting works until seven levels at this stage
*/
func TestVisualOutputHigherOrderStatementNestedComponentCombinationsDefaultExampleWithComplexity(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that Cac{{Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] {Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)} [OR] Cac{A(actor5) I(act5)}}} [XOR] Cac{A(actor6) I(act6)}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(true)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualHigherOrderStatementNestedComponentCombinationsDefaultStatementWithComplexity.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests higher-order statement combinations ({Cac{ ... } [AND] Cac{ ... } [XOR] {Cac{ ... } [AND] Cac{ ... }}}) or component-level nesting combinations in visual output.
NOTE: Nesting works until seven levels at this stage
*/
func TestVisualOutputHigherOrderStatementNestedComponentCombinations(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "Cac{Cac{Cac{A(actor1) I(aim1) Bdir(object1)} [AND] Cac{A(actor2) I(aim2) Bdir(object2)} [AND] Cac{A(actor4) I(aim4)}} [OR] Cac{Cac{A(actor3) I(aim3) Bdir(object3)} [XOR] Cac{Cac{A(actor6) I(aim6) Bdir(object6)} [AND] Cac{Cac{A(actor7) I(aim7) Bdir(object7)} [XOR] Cac{A(actor8) I(aim8)}}}}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualHigherOrderStatementNestedComponentCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", outputString)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the presence of excess symbols or missing separating whitespaces in input and parser tolerance (based on Regex).
*/
func TestVisualOutputExcessiveSymbolsOrMissingWhitespaceInNestedComponentCombinations(t *testing.T) {

	// Standard book statement, but including excess comma preceding logical operator (, [OR]),
	// excessive text in combination parentheses (unnecessary Text,.;  , unnecessary text),
	// and missing whitespace between logical operator and component specification ([OR]Cac)
	text := "The A(Program Manager) D(may) I(initiate) Bdir,p((suspension [XOR] revocation)) Bdir(proceedings) against a Bind,p(certified) Bind(operation): " +
		"Cac{unnecessary Text,.Cac{when the A(Program Manager) I(believes) that Bdir{a A,p(certified) A(operation) I((has violated [OR] is not in compliance)) " +
		"Bdir(with (the Act [OR] regulations in this part))}}, [OR]Cac{when a A((certifying agent [OR] State organic program’s governing State official)) " +
		"I(fails to enforce) Bdir((the Act [OR] regulations in this part)).} , unnecessary text}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualExcessiveSymbolsInNestedComponentCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests component-level statement combinations that contain embedded component-level nesting.

Variation of test #TestVisualOutputComponentLevelNestingInNestedComponentCombinationsInComponentPairs

Key difference: combinations as opposed to component pairs
*/
func TestVisualOutputComponentLevelNestingInNestedComponentCombinations(t *testing.T) {

	// Statement with component-level nesting embedded in statement combinations (i.e., Cac{Cac{ Bdir{} } [AND] Cac{ ... }})
	text := "A(Program Manager) D(may) I(administer) Bdir(sanctions) Cac{Cac{A(Program Manager) I(suspects) Bdir{A(farmer) " +
		"I((violates [OR] does not comply)) with Bdir(regulations)}} [OR] Cac{A(Program Manager) I(has witnessed) " +
		"Bdir,p(farmer's) Bdir(non-compliance) Cex(in the past)}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// Generate tree output
	output, err1 := s.PrintTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(), tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err1.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating visual tree output. Error: ", err1.Error())
	}

	outputString := output.String()

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err2 := os.ReadFile("TestOutputVisualComponentLevelNestingInNestedComponentCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests component-level statement combinations that contain embedded component-level nesting and are wrapped in component pairs (extrapolation)

Variation of test #TestVisualOutputComponentLevelNestingInNestedComponentCombinations

Key difference: component pairs
*/
func TestVisualOutputComponentLevelNestingInNestedComponentCombinationsInComponentPairs(t *testing.T) {

	// Statement with component-level nesting embedded in statement combinations (i.e., {Cac{ Bdir{} } [AND] Cac{ ... }})
	text := "A(Program Manager) D(may) I(administer) Bdir(sanctions) {Cac{A(Program Manager) I(suspects) Bdir{A(farmer) " +
		"I((violates [OR] does not comply)) with Bdir(regulations)}} [OR] Cac{A(Program Manager) I(has witnessed) " +
		"Bdir,p(farmer's) Bdir(non-compliance) Cex(in the past)}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualComponentLevelNestingInNestedComponentCombinationsInComponentPairs.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the combined use of component combinations, component pairs, as well as nested component statements.
*/
func TestVisualOutputNestedCombinationsComponentLevelNestingAndComponentPairs(t *testing.T) {

	// Statement with multi-level nesting, combinations as well as component pairs
	text := "A(actor1) I(aim1) Cac{Cac{A(actor2) I(aim2)} [XOR] Cac{A(actor3) I(aim3)}} {Bdir(directobject1) Bind(indirectobject1) [OR] Bdir{ A(actor4) I(aim4) Bdir(directobject2) Cac{A(actor5) I(aim5)}} Bind(indirectobject2)}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualNestedCombinationsComponentLevelNestingAndComponentPairs.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for multi-level pair combinations.
*/
func TestVisualOutputMultiLevelComponentPair(t *testing.T) {

	// Statement with multi-level nesting, combinations as well as component pairs
	text := "{ A(actor1) I(aim1) [XOR] {A(actor2) I(aim2) [AND] A(actor3) I(aim3)}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualMultiLevelComponentPair.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for multi-level pair combinations
*/
func TestVisualOutputComplexNestedCombinationsWithComponentPairs(t *testing.T) {

	// Statement with multi-level nesting, combinations as well as component pairs
	text := "{Cac{Cac{A(actor1) I(aim1) Bdir(object1)} [AND] Cac{A(actor2) I(aim2) Bdir(object2)} [AND] Cac{A(actor4) I(aim4)}} [OR] Cac{Cac{A(actor3) I(aim3) Bdir(object3)} [XOR] Cac{Cac{A(actor6) I(aim6) Bdir(object6)} [AND] Cac{Cac{A(actor7) I(aim7) Bdir(object7)} [XOR] Cac{A(actor8) I(aim8)}}}}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualComplexNestedCombinationsWithComponentPairs.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for component pairs in nested components.
*/
func TestVisualOutputComponentPairsInNestedComponents(t *testing.T) {

	// Statement with multi-level nesting, combinations as well as component pairs
	text := "A(Individuals) D(must) { I(monitor) Bdir(compliance) [AND] I(report) Bdir(violation) } Cac(in the case of (repeated offense [OR] other reasons)) O{ A(actor2) D(must) {I(enforce) Bdir(compliance) [OR] I(delegate) Bdir(enforcement)}}"

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualComponentPairsInNestedComponents.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for degree of variability with component and statement combinations.
Suppresses shared elements in output.
*/
func TestVisualOutputDegreeOfVariabilityComponentAndStatementCombinationsWithoutSharedElements(t *testing.T) {

	// Statement with various combinations
	text := "A(certifying agent [AND] (borrower [OR] wife)) M(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act [OR] regulations of this part) " +
		"concerning " +
		"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		") " +
		"Cac{Cac{A(actor2) I(aim2)} [XOR] Cac{A(actor3) I(aim3)}} " +
		"Cex(for compliance with the (Act [XOR] regulations in this part))."

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(true)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualDegreeOfVariabilityComponentAndStatementCombinationsWithoutSharedElements.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for degree of variability with component and statement combinations.
Includes shared elements in output.
*/
func TestVisualOutputDegreeOfVariabilityComponentAndStatementCombinationsWithSharedElements(t *testing.T) {

	// Statement with various combinations
	text := "A(certifying agent [AND] (borrower [OR] wife)) M(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act [OR] regulations of this part) " +
		"concerning " +
		"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		") " +
		"Cac{Cac{A(actor2) I(aim2)} [XOR] Cac{A(actor3) I(aim3)}} " +
		"Cex(for compliance with the (Act [XOR] regulations in this part))."

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(true)
	// Activate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(true)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualDegreeOfVariabilityComponentAndStatementCombinationsWithSharedElements.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for degree of variability with component pair combinations.
*/
func TestVisualOutputDegreeOfVariabilityComponentPairCombinations(t *testing.T) {

	// Statement with various combinations
	text := "A(certifying agent [AND] (borrower [OR] wife)) M(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act [OR] regulations of this part) " +
		"concerning " +
		"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		") " +
		"{Cac{A(actor2) I(aim2)} [XOR] Cac{A(actor3) I(aim3)}} " +
		"Cex(for compliance with the (Act [XOR] regulations in this part))."

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(true)
	// Deactivate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(false)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualDegreeOfVariabilityComponentPairCombinations.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for nested properties, including component pairs, simple nested statements and primitive properties.

Prints properties as nested property tree.
*/
func TestVisualOutputNestedPropertiesIncludingComponentPairsAndNestedPropertyAndPrimitivePropertyAsPropertyTree(t *testing.T) {

	// Statement with multi-level nesting, combinations as well as component pairs
	text := "Such E(notification) M(shall) F(provide): (1) A P(description of each noncompliance); " +
		"(2) The P(facts upon which the notification of noncompliance is based); and " +
		"(3) The P1(date) " +
		// Component-pair property
		"P1,p{by which the A(certified operation) D(must) " +
		"{I(rebut [XOR] correct) Bdir,p(each) Bdir(noncompliance) [AND] I(submit) Bdir,p(supporting) " +
		"Bdir(documentation) of Bdir,p(each such correction) Cac(when correction is possible)}} " +
		// Primitive property
		"P1,p(private component) " +
		// nested property
		"P1,p{where E(date) F(is defined) in the P(Gregorian calendar)}."

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing
	tree.SetFlatPrinting(false)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Activate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(true)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualPrivateComplexNestedPropertiesPropertyTree.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for nested properties, including component pairs, simple nested statements and primitive properties.

Prints properties as string associated with component.
*/
func TestVisualOutputNestedPropertiesIncludingComponentPairsAndNestedPropertyAndPrimitivePropertyAsFlatString(t *testing.T) {

	// Statement with multi-level nesting, combinations as well as component pairs
	text := "Such E(notification) M(shall) F(provide): (1) A P(description of each noncompliance); " +
		"(2) The P(facts upon which the notification of noncompliance is based); and " +
		"(3) The P1(date) " +
		// Component-pair property
		"P1,p{by which the A(certified operation) D(must) " +
		"{I(rebut [XOR] correct) Bdir,p(each) Bdir(noncompliance) [AND] I(submit) Bdir,p(supporting) " +
		"Bdir(documentation) of Bdir,p(each such correction) Cac(when correction is possible)}} " +
		// Primitive property
		"P1,p(private component) " +
		// nested property
		"P1,p{where E(date) F(is defined) in the P(Gregorian calendar)}."

	// Deactivate annotations
	tabular.SetIncludeAnnotations(false)
	// Deactivate flat printing -- the key difference to the preceding test
	tree.SetFlatPrinting(true)
	// Deactivate binary tree printing
	tree.SetBinaryPrinting(false)
	// Deactivate moving of activation conditions
	tree.SetMoveActivationConditionsToFront(false)
	// Deactivate DoV
	tabular.SetIncludeDegreeOfVariability(false)
	// Activate shared elements
	tree.SetIncludeSharedElementsInVisualOutput(true)

	// Parse statement
	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	output, err2 := stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), tabular.IncludeAnnotations(),
		tabular.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Error when generating node tree:", err2)
	}

	outputString := output

	fmt.Println("Generated output: " + outputString)

	// Read reference file
	content, err3 := os.ReadFile("TestOutputVisualPrivateComplexNestedPropertiesFlatPrinting.test")
	if err3 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	fmt.Println("Output:", output)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err4 := tabular.WriteToFile("errorOutput.error", outputString, true)
		if err4 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err4.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}
