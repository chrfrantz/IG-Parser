package exporter

import (
	"IG-Parser/tree"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Separator for main statement ID (e.g., 123) and suffix for introduced substatement (e.g., .1, i.e., 123.1)
const stmtIdSeparator = "."
// Separator for logical operators in enumerations of statement references (e.g., OR[650.1,650.2, ...])
const logicalOperatorStmtRefSeparator = ","
// Symbol separating component symbol and indices (e.g., Bdir vs. Bdir_1, Bdir_2, etc.)
const indexSymbol = "_"
// Statement ID prefix to ensure interpretation as text field (does not remove trailing zeroes)
const stmtIdPrefix = "'"
// Separator for logical operator expressions (e.g., OR[650.1,650.2]|AND[123.1,123.2])
const logicalOperatorSeparator = ";"

/*
Generates statement output in Google Sheets format.
Input:
- Atomic statements with corresponding node references [statement][node references]
- Map with with component name as key and corresponding number of columns in input stmts (i.e., same component can have
  values across multiple columns)
- References to entries for given nodes as indicated by logical operators, and used to produce corresponding linkages
  (e.g., AND[row1, row2, etc.])
- ID to be used as prefix for generation of substatement IDs (e.g., ID 5 produces substatements 5.1, 5.2, etc.)
 */
func GenerateGoogleSheetsOutput(stmts [][]*tree.Node, componentFrequency map[string]int, logicalLinks []map[*tree.Node][]string, stmtId string) (string, tree.ParsingError) {
	output := ""

	// Quote to terminate input string for Google Sheets interpretation
	quote := "\""
	// Line prefix for Google Sheets
	prefix := "=SPLIT(" + quote
	// Linebreak at the end of each entry
	linebreak := "\n"
	// Column separator used for Sheets output
	separator := ";"
	// Line suffix for Google Sheets
	suffix := quote + ", \"" + separator + "\")" + linebreak


	// Caches column header symbols by component index for reuse in logical operator construction
	headerSymbols := []string{}

	// Generate headers
	if componentFrequency != nil && len(componentFrequency) != 0 {

		output += prefix
		output += "Statement ID" + separator

		// Iterate through header frequencies and create header row
		output, headerSymbols = generateHeaderRow(output, componentFrequency, separator)

		// Complete line
		output += suffix
		//fmt.Println("Header: " + output)
	}

	// Default error during parsing
	errorVal := tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}

	// Generate entries
	for stmtCt, statement := range stmts {
		//fmt.Println("Statement ", stmtCt, ": ", statement)
		// Start new row
		output += prefix
		// Add statement ID for specific instance
		subStmtId := generateStatementIDint(stmtId, stmtCt + 1)
		output += stmtIdPrefix + subStmtId + separator
		ct := 0
		// String linking all logical operators for a given row
		logicalValue := ""
		// Iterate over component index (i.e., column)
		for componentIdx := range statement {
			// Append element value as output for given cell
			output += statement[componentIdx].Entry

			fmt.Println("Source node: ", statement[componentIdx])

			// Now generate logical links expression corresponding to particular entry (component index in statement instance)
			logicalValue, errorVal = generateLogicalLinksExpressionForGivenComponentValue(logicalValue, statement,
				componentIdx, headerSymbols, logicalLinks, stmtId)
			if errorVal.ErrorCode != tree.PARSING_NO_ERROR {
				return output, errorVal
			}

			// Only append separator if no more elements in the statement (i.e., no further components)
			ct++
			if ct < len(statement) {
				output += separator
			}
			//fmt.Println("-->", statement[componentIdx])
		}

		// Append the logical expression at the end of each row
		if logicalValue != "" {
			// Append separator
			output += separator
			// Append logical combination column
			output += logicalValue
			// Reset for next round
			logicalValue = ""
		}
		// Append suffix to each row
		output += suffix
	}
	return output, errorVal
}

/*
Generates IG 2.0 header row and appends it to given string based on component frequency input. It further returns a slice
containing header information.
 */
func generateHeaderRow(stringToAppendTo string, componentFrequency map[string]int, separator string) (string, []string) {

	// Header symbols to be returned for later use
	headerSymbols := []string{}
	// Iterate through component reference map
	for _, symbol := range tree.IGComponentSymbols {
		i := 0
		// Print headers as often as occurring in input file (stmtCt.e., one header for each column)
		for i < componentFrequency[symbol] {
			stringToAppendTo += tree.IGComponentSymbolNameMap[symbol]
			// Store symbols for columns including indices in order of occurrence for use in logical operators
			headerSymbol := symbol
			// Introduce indices if multiple of the same component
			if componentFrequency[symbol] > 1 {
				// Append suffix for header string
				stringToAppendTo += indexSymbol + strconv.Itoa(i + 1)
				// Append suffix for cached header IDs (for logical operators)
				headerSymbol += indexSymbol + strconv.Itoa(i + 1)
			}
			// Store key for header used in logical operators
			headerSymbols = append(headerSymbols, headerSymbol)
			stringToAppendTo += separator
			i++
		}
	}
	// Cut off last separator
	stringToAppendTo = stringToAppendTo[0:len(stringToAppendTo)-len(separator)]
	// Return generated string as well as symbol map
	return stringToAppendTo, headerSymbols
}

/*
Generates logical expression string for given component entry for a given statement.
It relies on the expression string as input, alongside the statement of concern, as well as component index.
In addition a slice of all header symbols is required (to generate reference to columns in logical expressions),
as well as the logical links for a given component value. Finally, the statement ID is used to generate the corresponding
substatement IDs used in the link references.
It returns the link for the particular table entry.
 */
func generateLogicalLinksExpressionForGivenComponentValue(logicalExpressionString string, statement []*tree.Node,
	componentIdx int, headerSymbols []string, logicalLinks []map[*tree.Node][]string, stmtId string) (string, tree.ParsingError) {
	// Check for logical operator linkage based on index
	linksForElement := logicalLinks[componentIdx]
	fmt.Println("Links for element: ", linksForElement)

	// Node key array (maintaining order of iteration)
	nodesKeys := []*tree.Node{}

	if len(linksForElement) > 0 {
		// Retrieve keys to determine order of iteration
		var firstKey *tree.Node
		for nd := range linksForElement {
			// Assign first key
			firstKey = nd
			// Then break out - since that is enough to get entire tree
			break
			// ALTERNATIVE: Sorting based on alphabet
			//nodesKeys = append(nodesKeys, nd)
		}
		// Sort by retrieving leaves for the given tree
		leaves := firstKey.GetSyntheticRootNode().GetLeafNodes()
		if len(leaves) > 0 {
			nodesKeys = leaves[0]
		} else {
			fmt.Println("No component keys to iterate over for logical relationships")
		}

		// ALTERNATIVE: Sorting based on alphabet by interface
		//sort.Sort(tree.ByEntry(nodesKeys))

		fmt.Println("Sorted keys: ", nodesKeys)
	}

	// Check that entries for own component value exist
	if linksForElement[statement[componentIdx]] != nil {
		// Iterate through all component values based on ordered keys
		for _, nodesKey := range nodesKeys {
			// Extract node
			otherNode := nodesKey
			// Extract references attached to node
			linkedElement := linksForElement[nodesKey]

			// NOTE: OLD iteration directly on elements leads to inconsistent iteration order - LEFT ONLY FOR DOCUMENTATION
			//for otherNode, linkedElement := range linksForElement {

			// if target node is different ...
			if otherNode != statement[componentIdx] {
				fmt.Println("Testing other node: ", otherNode, " with elements ", linkedElement)
				// find operator
				res, ops, err := tree.FindLogicalLinkage(statement[componentIdx], otherNode)
				if err.ErrorCode != tree.TREE_NO_ERROR {
					errorMsg := fmt.Sprint("Error when parsing retrieving operator linkages: ", err.ErrorMessage)
					log.Println(errorMsg)
					fmt.Errorf("%v", errorMsg)
					return "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_LOGICAL_EXPRESSION_GENERATION}
				}
				if res {
					fmt.Println("Collapsing adjacent AND operators ...")
					// Collapse adjacent AND operators
					ops = tree.CollapseAdjacentOperators(ops, []string{tree.AND})

					fmt.Println("Node has linkage ", ops)
					// ... and append to logical expression column string
					logicalExpressionString += fmt.Sprint(ops)
					// Statement component identifier
					logicalExpressionString += "." + headerSymbols[componentIdx] + "."
					// Leading bracket
					logicalExpressionString +=	"["
					// Prepare intermediate structure to store statement references
					stmtsRefs := ""

					fmt.Println("Target node IDs: ", linkedElement)
					for lks := range linkedElement {
						//fmt.Println("Found pointer from ", statement[componentIdx] ," to ", otherNode , " as ", generateStatementID(stmtId, lks + 1))
						// Append actual statement id
						stmtsRefs += generateStatementIDString(stmtId, linkedElement[lks])
						if lks < len(linkedElement)-1 {
							stmtsRefs += logicalOperatorStmtRefSeparator
						}
					}

					// Add trailing bracket and column ref (to be reviewed)
					logicalExpressionString += stmtsRefs + "]" + logicalOperatorSeparator
				}
				fmt.Println("Added logical relationships for value ", otherNode, ", elements: ", logicalExpressionString)
			}
		}
	}
	// Return generated logical expression for given component
	return logicalExpressionString, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Generate statement ID from main statement ID and index of iterated substatement
 */
func generateStatementIDString(mainID string, subStmtIndex string) string {
	return mainID + stmtIdSeparator + subStmtIndex
}

/*
Generate statement ID from main statement ID and index of iterated substatement
*/
func generateStatementIDint(mainID string, subStmtIndex int) string {
	return mainID + stmtIdSeparator + strconv.Itoa(subStmtIndex)
}

/*
Writes data to given file - overwrites file if existing
 */
func WriteToFile(filename string, content string) error {

	// Create file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Prepare data
	data := []byte(content)

	// Write data
	_, err2 := f.Write(data)
	if err2 != nil {
		return err2
	}
	log.Println("Wrote file " + filename)

	// No error
	return nil
}
