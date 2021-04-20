package exporter

import "IG-Parser/tree"

/*
Flatten input structure into simple array
 */
func flatten(input [][]*tree.Node) []*tree.Node {
	output := []*tree.Node{}
	for _, v := range input {
		for _, v2 := range v {
			output = append(output, v2)
		}
	}
	return output
}
