package main

import (
	"IG-Parser/parser"
	"IG-Parser/tree"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"regexp"
	"strings"
)

func main0() {

	text := "National Organic Program's Program Manager, on behalf of the Secretary, may (inspect and [AND] review) (certified production and [AND] handling operations and [AND] accredited certifying agents) for compliance with the (Act or [XOR] regulations in this part)."

	// logical operators
	r, _ := regexp.Compile("\\[[a-zA-Z]+\\]")
	// parenthesized expressions
	r, _ = regexp.Compile("\\(([a-zA-Z]+\\s)+\\[[a-zA-Z]\\]")//+\\s[a-zA-Z]\\)")
	//left hand only
	r, _ = regexp.Compile("\\(([a-zA-Z]+\\s)+\\[AND")
	//left and right hand
	//r, _ = regexp.Compile("\\(([a-zA-Z]+\\s)+\\[AND\\](\\s[a-zA-Z]+)+\\)")
	//left and right hand (left AND/XOR/OR right)
	r, _ = regexp.Compile("\\(([a-zA-Z]+\\s)+\\[(AND|OR|XOR)\\](\\s[a-zA-Z]+)+\\)")
	// left, right and listings separated by logical operator
	r, _ = regexp.Compile("\\(([a-zA-Z]+\\s)+(\\[(AND|OR|XOR)\\](\\s[a-zA-Z]+)+\\s?)+\\)")

	fmt.Println(r.MatchString(text))
	fmt.Println(r.FindAllStringSubmatch(text, -1))

	for k,v := range r.FindAllStringSubmatch(text, -1){
		fmt.Println(k)
		fmt.Println(v[0])
	}

	//fmt.Println(text)
}

var words = "([a-zA-Z',;]+\\s*)+"
var wordsWithParentheses = "([a-zA-Z',;()\\[\\]]+\\s*)+"
var logicalOperators = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"

func main3() {
	diff := diffmatchpatch.New()

	diffs := diff.DiffMain("world", "helloe world etx sdosdgsdklcx ", false)

	fmt.Println(diffs)

	diffs = diff.DiffMain("sox wulli beloved sds", "wulli", false)

	fmt.Println(diffs)
	fmt.Println(diffs[0].Type.String())
}

func main() {
	root := tree.Node{}
	root.InsertLeftNode(root)

	root.InsertRightNode(root)

}

func main5() {



	text := "( shared left ( inner left (Far left side [AND] Left side information [AND] inner right)) [AND] right information)"


	//component := "I"

	//text := "Cex((on behalf of the Secretary) [AND] (for compliance with the (Act or [XOR] regulations in this part)))"

	//r, _ := regexp.Compile(component + "\\(" + wordsWithParentheses + "(\\[" + logicalOperators + "\\]\\s" + wordsWithParentheses + ")*\\)")

	//result := r.FindAllStringSubmatch(text, -1)
	//result := extractComponent(component, text)

	text = "left middle right"

	//text = "(shared left ((left middle [AND] ((left [OR] right) [AND] right)) [AND] (shared right [OR] pisa)) right stuff)"

	//text = "(shared (left [AND] (inner left [OR] inner right)))"

	//text = "(lefty (left (ADND [OR] SGSD) (smalleft [AND] mouse) (inner left [AND] inner right)))"
	//text = "(lefty (left (ADND [OR] SGSD) (smalleft [AND] mouse) (otherLeft [OR] otherRight) inner right ))"

	//text = "(Left side information (source) [AND] middle information and [AND] right-side)"

	text = "(left0 ((left1 (left2 (left3 [AND] mid3 [AND] right3))) right1))"

	//text = "( shared left ( inner left (innermost left ((left-most information [AND] Left side information [AND] middle information) [AND] right information))) shared right)"

	//text = "( shared left (Left side information [XOR] middle information) shared right)"

	//text = "(left (inner left [AND] inner right [AND] nother) other)"

	//text = "(left (inner left [AND] inner right)   (inner2Left [OR] inner2Right) other)"

	/*
	text = "(left (inner left [AND] inner right) middle (inner2Left [OR] inner2Right) other)"
	*/
	//node := tree.Node{}
	node, _, err := parser.ParseDepth(text, false)


	fmt.Println(err.Error())

	fmt.Println("Final: " + node.String())

	fmt.Println(fmt.Sprint(node.ElementOrder))

}

func extractComponent(component string, input string) string {

	// Parentheses count to check for balance
	parCount := 0

	startPos := strings.Index(input, component)
	if startPos == -1 {
		log.Fatal("Component signature not found")
		return ""
	}

	for i, letter := range input[startPos:] {

		switch string(letter) {

		case "(":
			parCount++
		case ")":
			parCount--
			if parCount == 0 {
				return input[startPos:startPos+i+1]
			}
		}
	}

	return ""
}

func main2(){

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), D(may) I(inspect and [AND] review) Bdir(certified production and [AND] handling operations and [AND] accredited certifying agents) Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// left, right and listings separated by logical operator
	r, _ := regexp.Compile("\\(([a-zA-Z]+\\s)+(\\[(AND|OR|XOR)\\](\\s[a-zA-Z]+)+\\s?)+\\)")

	components := "(A|Cex|D|I|Bdir|Bind)"
	// pure
	words := "([a-zA-Z',;]+\\s?)+"
	// with structure
	//words = "([a-zA-Z',;]+\\s?)+"
	logicalOperators := "(AND|OR|XOR)"
	// Attribute
	r, _ = regexp.Compile("A\\(" + words + "\\)")//([a-zA-Z]+\\s)+(\\[(AND|OR|XOR)\\](\\s[a-zA-Z]+)+\\s?)+\\)")
	// Attribute and logical combinations embedded - Detecting components
	r, _ = regexp.Compile(components + "\\(" + words + "(\\[" + logicalOperators + "\\]\\s" + words + ")*\\)")
	// Generic combinations - Detecting generic combinations
	r, _ = regexp.Compile("\\s" + "\\(" + words + "(\\[" + logicalOperators + "\\]\\s" + words + ")*\\)")

	fmt.Println(r.MatchString(text))
	fmt.Println(r.FindAllStringSubmatch(text, -1))

	for k,v := range r.FindAllStringSubmatch(text, -1){
		fmt.Println(k)
		fmt.Println(v[0])
	}

	// if generic combinations, check if they contain components

		// if they contain components, store reference to component types combination for later decomposition
		//Concatenation of "shared parts"? Test

	// if components, check if they contain combinations
		//if they contain combinations, expand into "shared part" + either combination

	// create map containing each component as key

}

func parseAttributes() {


}