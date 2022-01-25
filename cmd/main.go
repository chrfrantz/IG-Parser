package main

import (
	"IG-Parser/exporter"
	"IG-Parser/parser"
	"IG-Parser/tree"
	"fmt"
	"log"
)

//var words = "([a-zA-Z',;]+\\s*)+"
//var wordsWithParentheses = "([a-zA-Z',;()]+\\s*)+"
//var logicalOperators = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"

func main() {

	//text := "National Organic Program's Program Manager, on behalf of the Secretary, may (inspect and [AND] review) (certified production and [AND] handling operations and [AND] accredited certifying agents) for compliance with the (Act or [XOR] regulations in this part)."

	//text := "A1,p(National Organic Program's) A1(Program Manager) I(sldkgjslkd) I(bansdkgs) " +
	//	"{Cac{A(sldkgjskl)} [OR] Cac{A(lksjgglkds)}}"
	//text := "I(suspend (a [AND] (b [OR] c)))"
	//text := "I(sustain (review [AND] (refresh [AND] drink)) srightstain)"
	//text := "I(sustain (review [AND] refresh) srightstain)"
	//text := "Cex(shared0 (left1 [XOR] left2) shared1 (mid1 [AND] mid2 [AND] mid3 [AND] mid4) shared2 (right1 [OR] right2 [OR] right3) shared3)"
	//text := "Cex(leftShared (left1 [XOR] right1) mid (left2 [AND] right2) rightShared (left3 [XOR] right3))"
	//text := "A(Frog) A(sdklgjdslk) Cex(leftShared (left1 [XOR] right1) mid (left2 [AND] right2) rightShared) Cex(banana)"
	//text := "Cex(leftShared (left1 [XOR] right1) rightShared) Cex(banana)"
	//text := "A(Owner) and A(Bidder) D(must not) I(((neglect [AND] forego) [OR] delay)) Bdir(negotiations) Cac{A(bidder) I(has initiated) Bdir(something) Cac{A(Officials) I(have warned) Bdir(about their performance), Cac{E(lsdkgjs) F(slkgdjsl) P(sldkgjsl)}}}"
	//text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2,p(another) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that {Cac[state]{A[role=monitored,type=animate](Operations) I[violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) I(has concluded) Bdir[type=activity](investigation)}}."
	//text := "The Congress finds and declares that it is the E(national policy) F([is] to (encourage [AND] assist)) the P(states) Cex{ A(states) I(to exercise) Cex(effectively) their Bdir(responsibilities) Bdir,p(in the coastal zone) Cex(through the (development [AND] implementation) of management programs to achieve wise use of the (land [AND] water) resources of the coastal zone, giving full consideration to (ecological [AND] cultural [AND] historic [AND] esthetic) values as well as the needs for compatible economic development), Cex{which E(programs) M(should) Cex(at least) F(provide for)— (A) the P1(protection) P,p1(of natural resources, including (wetlands [AND] floodplains [AND] estuaries [AND] beaches [AND] dunes [AND] barrier islands [AND] coral reefs [AND] fish and wildlife and their habitat) within the coastal zone), the P2(management) P,p2((of coastal development to minimize the loss of (life [AND] property) caused by improper development in (flood-prone [AND] storm surge [AND] geological hazard [AND] erosion-prone) areas [AND] in areas likely to be (affected by [OR] vulnerable to) (sea level rise [AND] land subsidence [AND] saltwater intrusion) [AND] by the destruction of natural protective features such as (beaches [AND] dunes [AND] wetlands [AND] barrier islands))), (C) the P3(management) P3,p(of coastal development to (improve [AND] safeguard [AND] restore) the quality of coastal waters, [AND] to protect (natural resources [AND] existing uses of those waters)), (D) P4,p1(priority) P4(consideration) P4,p2(being given to (coastal-dependent (uses [AND] orderly processes) for siting major facilities related to (national defense [AND] energy [AND] fisheries development [AND] recreation [AND] ports [AND] transportation), [AND] the location to the maximum extent practicable of new (commercial [AND] industrial) developments (in [XOR] adjacent) to areas where such development already exists)), (E) P5,p1(public) P5(access) P5,p2(to the coasts for recreation purposes), (F) P6(assistance) P6,p(in the redevelopment of (deteriorating urban (waterfronts [AND] ports) [AND] sensitive (preservation [AND] restoration) of (historic [AND] cultural [AND] esthetic) coastal features)), (G) P7(the (coordination [AND] simplification) of procedures) P7,p1(in order to ensure expedited governmental decision making for the management of coastal resources), (H) P8((continued (consultation [AND] coordination) with, [AND] the giving of adequate consideration to the views of affected Federal agencies)), (I) P9(the giving of ((timely [AND] effective) notification of , [AND] opportunities for (public [AND] local) government participation in coastal management decision making)), (J) P10(assistance) P10,p1(to support comprehensive (planning [AND] conservation [AND] management) for living marine resources) P10,p1,p1(including planning for (the siting of (pollution control [AND] aquaculture facilities) within the coastal zone [AND] improved coordination between ((State [AND] Federal) coastal zone management agencies [AND] (State [AND] wildlife) agencies))) }}"
	//text := "Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents)"
	text := "Cac(do (something [AND] (reside [OR] own) sufficient (agricultural [OR] forest) land to ...))"

	exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	exporter.SetDynamicOutput(false)

	fmt.Println("Shared mode:")
	fmt.Println(tree.SHARED_ELEMENT_INHERITANCE_MODE)

	tree.SetFlatPrinting(true)

	output, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}

	fmt.Println(output.PrintTree(nil, true, false, false, 0))

	out, err := output.PrintTree(nil, true, false, false, 0)

	exporter.WriteToFile("example.txt", out)

}
