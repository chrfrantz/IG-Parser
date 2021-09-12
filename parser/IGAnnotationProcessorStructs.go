package parser

import "IG-Parser/tree"

/*
Delimiter for component suffix entries (e.g., "1,p")
 */
const SUFFIX_DELIMITER = ","

/*
Combines two maps by copying content of second map into first and returning it. Where keys conflict, the
parameter overwrite determines whether the entry is overwritten (i.e., second map entry is used), or the original
entry retained.
 */
func CombineMaps(map1 map[*tree.Node]*tree.Node, map2 map[*tree.Node]*tree.Node, overwrite bool) map[*tree.Node]*tree.Node {
	for k, v := range map2 {
		if map1[k] != nil {
			// Overwrite existing entry only if indicated
			if overwrite {
				map1[k] = map2[k]
			}
		} else {
			// Add entry if not existing
			map1[k] = v
		}
	}
	return map1
}
