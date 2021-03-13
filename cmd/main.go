package main

import (
	igTree "IG-Parser"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

var words = "([a-zA-Z',;]+\\s?)+"
var logicalOperators = "(" + igTree.AND + "|" + igTree.OR + "|" + igTree.XOR + ")"

func main() {

	//text := "National Organic Program's Program Manager, on behalf of the Secretary, may (inspect and [AND] review) (certified production and [AND] handling operations and [AND] accredited certifying agents) for compliance with the (Act or [XOR] regulations in this part)."

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), D(may) I(inspect and [AND] review) Bdir(certified production and [AND] handling operations and [AND] accredited certifying agents) Cex(for compliance with the (Act or [XOR] regulations in this part))."


	s := parseStatement(text)

	fmt.Println(s.String())

}

func parseStatement(text string) igTree.Statement {
	s := igTree.Statement{}

	a := parseAttributes(text)
	switch len(a) {
	case 1:		s.Attributes = igTree.ComponentLeafNode(a[0][0], igTree.ATTRIBUTES)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(a)) + " items.")
	default: 	log.Println("No Attributes found")
	}

	d := parseDeontic(text)
	switch len(d) {
	case 1:		s.Deontic = igTree.ComponentLeafNode(d[0][0], igTree.DEONTIC)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(d)) + " items.")
	default: 	log.Println("No Deontic found")
	}

	i := parseAim(text)
	switch len(i) {
	case 1:		s.Aim = igTree.ComponentLeafNode(i[0][0], igTree.AIM)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(i)) + " items.")
	default: 	log.Println("No Aim found")
	}

	bdir := parseDirectObject(text)
	switch len(bdir) {
	case 1:		s.DirectObject = igTree.ComponentLeafNode(bdir[0][0], igTree.DIRECT_OBJECT)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bdir)) + " items.")
	default: 	log.Println("No Direct Object found")
	}

	bdirp := parseDirectObjectProperty(text)
	switch len(bdirp) {
	case 1:		s.DirectObjectProperty = igTree.ComponentLeafNode(bdirp[0][0], igTree.DIRECT_OBJECT_PROPERTY)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bdirp)) + " items.")
	default: 	log.Println("No Direct Object Property found")
	}

	bind := parseIndirectObject(text)
	switch len(bind) {
	case 1:		s.IndirectObject = igTree.ComponentLeafNode(bind[0][0], igTree.INDIRECT_OBJECT)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bind)) + " items.")
	default: 	log.Println("No Indirect Object found")
	}

	bindp := parseIndirectObjectProperty(text)
	switch len(bindp) {
	case 1:		s.IndirectObjectProperty = igTree.ComponentLeafNode(bindp[0][0], igTree.INDIRECT_OBJECT_PROPERTY)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bindp)) + " items.")
	default: 	log.Println("No Indirect Object Property found")
	}

	return s

	/*
	e := parseConstitutedEntity(text)[0]
	s.ConstitutedEntity = igTree.ComponentLeafNode(e, igTree.CONSTITUTED_ENTITY)

	ep := parseConstitutedEntityProperty(text)[0]
	s.ConstitutedEntityProperty = igTree.ComponentLeafNode(ep, igTree.CONSTITUTED_ENTITY_PROPERTY)

	f := parseConstitutingFunction(text)[0]
	s.ConstitutiveFunction = igTree.ComponentLeafNode(f, igTree.CONSTITUTIVE_FUNCTION)

	p := parseConstitutingProperties(text)[0]
	s.ConstitutingProperties = igTree.ComponentLeafNode(p, igTree.CONSTITUTING_PROPERTIES)

	pp := parseConstitutingPropertiesProperty(text)[0]
	s.ConstitutingPropertiesProperty = igTree.ComponentLeafNode(pp, igTree.CONSTITUTING_PROPERTIES_PROPERTY)

	/* TODO: Review for complex statements */
	/*cac := parseActivationCondition(text)[0]
	s.ActivationConditionSimple = igTree.ComponentLeafNode(cac, igTree.ACTIVATION_CONDITION)

	cex := parseExecutionConstraint(text)[0]
	s.ExecutionConstraintSimple = igTree.ComponentLeafNode(cex, igTree.EXECUTION_CONSTRAINT)

	return s
	*/
}

func parseAttributes(text string) [][]string {
	return parse(igTree.ATTRIBUTES, text)
}

func parseAttributesProperty(text string) [][]string {
	return parse(igTree.ATTRIBUTES_PROPERTY, text)
}

func parseDeontic(text string) [][]string {
	return parse(igTree.DEONTIC, text)
}

func parseAim(text string) [][]string {
	return parse(igTree.AIM, text)
}

func parseDirectObject(text string) [][]string {
	return parse(igTree.DIRECT_OBJECT, text)
}

func parseDirectObjectProperty(text string) [][]string {
	return parse(igTree.DIRECT_OBJECT_PROPERTY, text)
}

func parseIndirectObject(text string) [][]string {
	return parse(igTree.INDIRECT_OBJECT, text)
}

func parseIndirectObjectProperty(text string) [][]string {
	return parse(igTree.INDIRECT_OBJECT_PROPERTY, text)
}

func parseConstitutedEntity(text string) [][]string {
	return parse(igTree.CONSTITUTED_ENTITY, text)
}

func parseConstitutedEntityProperty(text string) [][]string {
	return parse(igTree.CONSTITUTED_ENTITY_PROPERTY, text)
}

func parseConstitutingFunction(text string) [][]string {
	return parse(igTree.CONSTITUTIVE_FUNCTION, text)
}

func parseConstitutingProperties(text string) [][]string {
	return parse(igTree.CONSTITUTING_PROPERTIES, text)
}

func parseConstitutingPropertiesProperty(text string) [][]string {
	return parse(igTree.CONSTITUTING_PROPERTIES_PROPERTY, text)
}

func parseActivationCondition(text string) [][]string {
	return parse(igTree.ACTIVATION_CONDITION, text)
}

func parseExecutionConstraint(text string) [][]string {
	return parse(igTree.EXECUTION_CONSTRAINT, text)
}

func parse(component string, text string) [][]string {
	log.Println("Invoking parsing of component " + component)

	r, _ := regexp.Compile(component + "\\(" + words + "(\\[" + logicalOperators + "\\]\\s" + words + ")*\\)")

	/*for k,v := range r.FindAllStringSubmatch(text, -1){
		fmt.Println(k)
		fmt.Println(v[0])
	}*/

	//fmt.Println()

	return r.FindAllStringSubmatch(text, -1)

}