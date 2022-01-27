package parser

import (
	"fmt"
	"regexp"
	"testing"
)

func TestSingleComponentSyntax(t *testing.T) {

	r, err := regexp.Compile(FULL_COMPONENT_SYNTAX)
	if err != nil {
		t.Fatal("Error during compilation:", err.Error())
	}

	text := "(Bklsdjgl#k{sgv sk}lvjds) dalsjglks() Bdir,p1[ruler=governor](jglkdsjgsiovs) " +
		"Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)} " +
		"P,p343(another component values#$) " +
		"E1{ A(actor two) I1(aim1) }" +
		" text outside"

	res := r.FindAllString(text, -1)

	fmt.Println("Matching component structure (primitive and nested)")
	fmt.Println(res)
	fmt.Println("Count:", len(res))

	firstElem := "Bdir,p1[ruler=governor](jglkdsjgsiovs)"

	if res[0] != firstElem {
		t.Fatal("Wrong element matched. Should be", firstElem, ", but is "+res[0])
	}

	secondElem := "Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)}"

	if res[1] != secondElem {
		t.Fatal("Wrong element matched. Should be", secondElem, ", but is "+res[1])
	}

	thirdElem := "P,p343(another component values#$)"

	if res[2] != thirdElem {
		t.Fatal("Wrong element matched. Should be", thirdElem, ", but is "+res[2])
	}

	fourthElem := "E1{ A(actor two) I1(aim1) }"

	if res[3] != fourthElem {
		t.Fatal("Wrong element matched. Should be", fourthElem, ", but is "+res[3])
	}

	if len(res) != 4 {
		t.Fatal("Wrong number of matched elements. Should be 4, but is", len(res))
	}
}

/*
Tests for combinations within text. Note that is does not test for terminated statement combinations. That is tested in statement parsing tests.
*/
func TestComponentCombinations(t *testing.T) {

	// Note: Only used in testing; in production NESTED_COMBINATIONS_TERMINATED is used
	r, err := regexp.Compile(NESTED_COMBINATIONS)
	if err != nil {
		t.Fatal("Error during compilation:", err.Error())
	}

	text := "(Aklsdjgl#k{sgv sk}lvjds) {[]jdskgl ds()} Bdir,p1[ruler=governor](jglkdsjgsiovs) Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)}" +
		"{A(dlkgjsg) I[dgisg](kjsdglkds) [AND] (Bdir{djglksjdgkd} Cex(A(sdlgjlskd)) [XOR] A(dsgjslkj) E(gklsjgls))}" +
		"{Cac{ A(actor) I(fjhgjh) Bdir(rtyui)} [XOR] Cac{A(ertyui) I(dfghj)}}" +
		"{Cac{ A(as(dslks)a) I(adgklsjlg)} [XOR] Cac(asas) [AND] Cac12[kgkg]{lkdjgdls} [OR] A(dslgkjds)}" +
		"{Cac(andsdjsglk) [AND] A(sdjlgsl) Bdir(jslkgsjlkgds)}" +
		"{Cac(andsdjsglk) [AND] ( A(sdjlgsl) [XOR] (A(sdoidjs) [OR] A(sdjglksj)))}" +
		"((dglkdsjg [AND] jdlgksjlkgd))"

	res := r.FindAllString(text, -1)

	fmt.Println("Refined matching combinations")
	fmt.Println(res)
	fmt.Println("Count:", len(res))

	firstElem := "{A(dlkgjsg) I[dgisg](kjsdglkds) [AND] (Bdir{djglksjdgkd} Cex(A(sdlgjlskd)) [XOR] A(dsgjslkj) E(gklsjgls))}"

	if res[0] != firstElem {
		t.Fatal("Wrong element matched. Should be", firstElem, ", but is "+res[0])
	}

	secondElem := "{Cac{ A(actor) I(fjhgjh) Bdir(rtyui)} [XOR] Cac{A(ertyui) I(dfghj)}}"

	if res[1] != secondElem {
		t.Fatal("Wrong element matched. Should be", secondElem, ", but is "+res[1])
	}

	thirdElem := "{Cac{ A(as(dslks)a) I(adgklsjlg)} [XOR] Cac(asas) [AND] Cac12[kgkg]{lkdjgdls} [OR] A(dslgkjds)}"

	if res[2] != thirdElem {
		t.Fatal("Wrong element matched. Should be", thirdElem, ", but is "+res[2])
	}

	fourthElem := "{Cac(andsdjsglk) [AND] A(sdjlgsl) Bdir(jslkgsjlkgds)}"

	if res[3] != fourthElem {
		t.Fatal("Wrong element matched. Should be", fourthElem, ", but is "+res[3])
	}

	fifthElem := "{Cac(andsdjsglk) [AND] ( A(sdjlgsl) [XOR] (A(sdoidjs) [OR] A(sdjglksj)))}"

	if res[4] != fifthElem {
		t.Fatal("Wrong element matched. Should be", fifthElem, ", but is "+res[4])
	}

	if len(res) != 5 {
		t.Fatal("Wrong number of matched elements. Should be 5, but is", len(res))
	}

}
