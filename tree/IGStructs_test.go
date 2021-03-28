package tree

import "testing"

func TestValidIGComponentName(t *testing.T) {

	res := validIGComponentName("Attributes")
	if !res {
		t.Fatal("Could not resolve valid component name")
	}

}

func TestValidIGSymbol(t *testing.T) {

	res := validIGComponentSymbol("A")
	if !res {
		t.Fatal("Could not resolve valid component name")
	}

}