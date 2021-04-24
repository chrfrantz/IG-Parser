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
// Separator for referenced statements in cell (e.g., multiple references to activation conditions, i.e., {65}.1,{65}.2)
const componentStmtRefSeparator = ","
// Symbol separating component symbol and indices (e.g., Bdir vs. Bdir_1, Bdir_2, etc.)
const indexSymbol = "_"
// Statement ID prefix to ensure interpretation as text field (does not remove trailing zeroes)
const stmtIdPrefix = "'"
// Separator for logical operator expressions (e.g., OR[650.1,650.2]|AND[123.1,123.2])
const logicalOperatorSeparator = ";"
// Left bracket for logical combination expressions
const logicalCombinationLeft = "["
// Right bracket for logical combination expressions
const logicalCombinationRight = "]"
// Left brace surrounding identifier for component-level nested statements
const componentNestedLeft = "{"
// Right brace surrounding identifier for component-level nested statements
const componentNestedRight = "}"
// Column identifier for Statement ID
const stmtIdColHeader = "Statement ID"
// Column identifier for logical linkage of components
const logLinkColHeaderComps = "Logical Linkage (Components)"
// Column identifier for logically linked statements (not just components)
const logLinkColHeaderStmts = "Logical Linkage (Statements)"

// Structure referencing ID (based on input ID), along with (nested) statement to be decomposed
type IdentifiedStmt struct {
	ID string // Generated ID as provided in output
	NestedStmt *tree.Node // Single statement wrapped in Node
}

/*
Generates array of statement maps corresponding to identified elements format. Includes parsing of nested statements.
Input:
- Atomic statements with corresponding node references [statement][node references]
- Map with with component name as key and corresponding number of columns in input stmts (i.e., same component can have
  values across multiple columns)
- References to entries for given nodes as indicated by logical operators, and used to produce corresponding linkages
  (e.g., AND[row1, row2, etc.])
- ID to be used as prefix for generation of substatement IDs (e.g., ID 5 produces substatements 5.1, 5.2, etc.)
Output:
- Array of statement entry maps (i.e., values for each component in given statement, i.e., [statement]map[component]componentValue)
- Array of header symbols (used for component linkage references)
- Array of header symbols names (for human-readable header construction)
 */
func generateTabularStatementOutput(stmts [][]*tree.Node, componentFrequency map[string]int, logicalLinks []map[*tree.Node][]string, stmtId string) ([]map[string]string, []string, []string, tree.ParsingError) {
	//output := ""

	// Quote to terminate input string for Google Sheets interpretation
	//quote := "\""
	// Line prefix for Google Sheets
	//prefix := "=SPLIT(" + quote
	// Linebreak at the end of each entry
	//linebreak := "\n"
	// Column separator used for Sheets output
	separator := ";"
	// Line suffix for Google Sheets
	//suffix := quote + ", \"" + separator + "\")" + linebreak

	// Caches column header symbols by component index for reuse in logical operator construction
	headerSymbols := []string{}
	// Caches column header names associated with symbols for human-readable header construction
	headerSymbolsNames := []string{}

	// Generate headers
	if componentFrequency != nil && len(componentFrequency) != 0 {

		//output += prefix
		//output += stmtIdColHeader + separator

		// Iterate through header frequencies and create header row
		_, headerSymbols, headerSymbolsNames = generateHeaderRow("", componentFrequency, separator)

		// Complete line
		//output += suffix
		//fmt.Println("Header: " + output)
	}

	// Default error during parsing
	errorVal := tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}

	// Structures for nested statements

	// Map containing inverse index for all observed statements to IDs
	componentNestedStmtsMap := make(map[*tree.Node]string)
	// Nested statement index
	nestedStatementIdx := 1
	// Statements nested on components - to be processed last
	componentNestedStmts := make([]IdentifiedStmt, 0)

	// Map of entries to be returned at the end
	entriesMap := make([]map[string]string, 0)

	// Generate entries
	for stmtCt, statement := range stmts {

		// Individual entry
		entryMap := make(map[string]string)

		//fmt.Println("Statement ", stmtCt, ": ", statement)
		// Start new row
		//output += prefix
		// Add statement ID for specific instance
		subStmtId := generateStatementIDint(stmtId, stmtCt + 1)
		// Add statement ID to entryMap
		entryMap[stmtIdColHeader] = subStmtId
		// Add to the output
		//output += stmtIdPrefix + subStmtId + separator
		//ct := 0
		// String linking all logical operators for a given row
		logicalValue := ""
		// Iterate over component index (i.e., column)
		for componentIdx := range statement {
			// Append element value as output for given cell
			if statement[componentIdx].IsEmptyNode() {
				// Empty entry - don't add anything
				fmt.Println("Found empty node")
			} else if statement[componentIdx].HasPrimitiveEntry() {
				// Regular leaf entry (i.e., component)
				// Save entry value into entryMap for given statement and component column
				entryMap[headerSymbols[componentIdx]] = statement[componentIdx].Entry.(string)
				// Add to output
				//output += statement[componentIdx].Entry.(string)
			} else {
				// Nested statements are stored for later processing, but assigned IDs and references added to calling row
				fmt.Println("Found complex entry (nested statement): " + fmt.Sprint(statement[componentIdx]))
				// Check for statement combination (i.e., node combination)

				// Add entry to array (assuming single nested statement)
				entryVals := []*tree.Node{statement[componentIdx]}

				// Check if combination contained; if so, flatten, and override
				if entryVals[0].IsCombination() {
					fmt.Println("Detected statement combination")
					// If combination of statements, retrieve all elements
					stmts := entryVals[0].GetLeafNodes()
					// Flatten array and override entry values for iteration
					entryVals = tree.Flatten(stmts)
				} else {
					fmt.Println("Detected individual nested statement")
				}

				// Iterate over all nested statements
				for _, entryVal := range entryVals {

					idToReferenceInCell := ""

					// Retrieve ID of already identified statements ...
					if nestedStmtID, ok := componentNestedStmtsMap[entryVal]; ok {
						// Prepare reference to be saved
						idToReferenceInCell = nestedStmtID
						// and attach to output
						//output += nestedStmtID
					} else {
						// ... else create new one
						// Generate ID for component-level nested statement
						nestedStmtId := componentNestedLeft +
							stmtId +
							componentNestedRight +
							stmtIdSeparator +
							strconv.Itoa(nestedStatementIdx)

						// Added component-level nested statement
						componentNestedStmts = append(componentNestedStmts,
							IdentifiedStmt{nestedStmtId, entryVal})
						// Add newly identified nested statement to lookup index
						componentNestedStmtsMap[entryVal] = nestedStmtId
						// Increase index for component-level nested statements (for next round)
						nestedStatementIdx++
						// Prepare reference to to-be component-level nested statements to output
						idToReferenceInCell = nestedStmtId
						//output += nestedStmtId
						fmt.Println("Parsing: Added nested statement (ID:", nestedStmtId, ", Val:", entryVal, ")")
					}

					// Save entry into entryMap for calling row
					if entryMap[headerSymbols[componentIdx]] != "" {
						// Add separator if already an entry
						entryMap[headerSymbols[componentIdx]] += componentStmtRefSeparator
					}
					// Append current value in any case
					entryMap[headerSymbols[componentIdx]] += idToReferenceInCell


				}
			}
			fmt.Println("Source/calling node: ", statement[componentIdx])


			// Process component-level logical linkage

			// Now generate logical links expression corresponding to particular entry (component index in statement instance)
			logicalValue, errorVal = generateLogicalLinksExpressionForGivenComponentValue(logicalValue, statement,
				componentIdx, headerSymbols, logicalLinks, stmtId)
			if errorVal.ErrorCode != tree.PARSING_NO_ERROR {
				return nil, nil, nil, errorVal
			}

			// Only append separator if no more elements in the statement (i.e., no further components)
			//ct++
			/*if ct < len(statement) {
				output += separator
			}*/
			//fmt.Println("-->", statement[componentIdx])
		}

		// Append the logical expression at the end of each row
		if logicalValue != "" {
			// Append separator
			//output += separator
			// Add to entryMap
			entryMap[logLinkColHeaderComps] = logicalValue
			// Append logical combination column
			//output += logicalValue
			// Reset for next round
			logicalValue = ""
		}
		// Add to entries map for statement to map for all statements (collection for return)
		entriesMap = append(entriesMap, entryMap)
		// Append suffix to each row
		//output += suffix
	}

	fmt.Println("Component-level nested statements to be decomposed: " + fmt.Sprint(componentNestedStmts))
	for _, val := range componentNestedStmts {

		fmt.Println("Nested Statement to parse, ID:", val.ID, ", Stmt:", val.NestedStmt)

		// Parse individual nested statements on component level
		_, nestedMap, nestedHeaders, nestedHeadersNames, err := GenerateGoogleSheetsOutputFromParsedStatement(val.NestedStmt.Entry.(tree.Statement), val.ID, "")
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return nil, nil, nil, errorVal
		}

		// Add linkages between statements (statement-level combinations)

		// Determine linkages to fellow nested statements
		stmtLinksString, err := generateLogicalLinksExpressionForStatements(val.NestedStmt, componentNestedStmts)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return nil, nil, nil, errorVal
		}

		// Add identified linkages to nestedMap (i.e., for all atomic statements)
		for i, _ := range nestedMap {
			nestedMap[i][logLinkColHeaderStmts] = stmtLinksString
		}

		// Add Logical linkage header if not already existing
		nestedHeaders = addElementIfNotExisting(logLinkColHeaderStmts, nestedHeaders)
		nestedHeadersNames = addElementIfNotExisting(logLinkColHeaderStmts, nestedHeadersNames)

		// Add nested entries to top-level list
		entriesMap = append(entriesMap, nestedMap...)

		// Merge headers to consider nested ones
		headerSymbols = tree.MergeSlices(headerSymbols, nestedHeaders, indexSymbol)
		// Merge header names to consider nested ones
		headerSymbolsNames = tree.MergeSlices(headerSymbolsNames, nestedHeadersNames, indexSymbol)
		// Append component-level nested statement to output
		//output += nestedOutput
	}

	// Organise headers

	// Move Statement ID to first position
	headerSymbols = moveElementToFirstPosition(stmtIdColHeader, headerSymbols, true)
	headerSymbolsNames = moveElementToFirstPosition(stmtIdColHeader, headerSymbolsNames, true)
	// Add statement logical linkages to second-last position
	headerSymbols = moveElementToLastPosition(logLinkColHeaderStmts, headerSymbols, true)
	headerSymbolsNames = moveElementToLastPosition(logLinkColHeaderStmts, headerSymbolsNames, true)
	// Add component logical linkages to last position
	headerSymbols = moveElementToLastPosition(logLinkColHeaderComps, headerSymbols, true)
	headerSymbolsNames = moveElementToLastPosition(logLinkColHeaderComps, headerSymbolsNames, true)

	return entriesMap, headerSymbols, headerSymbolsNames, errorVal
}

/*
Resolves all logical linkages to other statements and returns those as compound logical expression (e.g., [AND][{65}.1],[AND][{65}.2])
 */
func generateLogicalLinksExpressionForStatements(sourceStmt *tree.Node, allNestedStmts []IdentifiedStmt) (string, tree.ParsingError) {
	logicalExpressionString := ""

	// Iterate over all nested statements
	for _, targetEntry := range allNestedStmts {

		// If statement is not the same
		if targetEntry.NestedStmt != sourceStmt {
			targetID := targetEntry.ID
			targetStmt := targetEntry.NestedStmt

			// Retrieve linkage
			res, ops, err := tree.FindLogicalLinkage(sourceStmt, targetStmt)
			if err.ErrorCode != tree.TREE_NO_ERROR {
				errorMsg := fmt.Sprint("Error when parsing retrieving operator linkages: ", err.ErrorMessage)
				log.Println(errorMsg)
				return "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_LOGICAL_EXPRESSION_GENERATION}
			}
			if res {
				fmt.Println("Collapsing adjacent AND operators ...")
				// Collapse adjacent AND operators
				ops = tree.CollapseAdjacentOperators(ops, []string{tree.AND})

				fmt.Println("Node has linkage ", ops)

				if logicalExpressionString != "" {
					logicalExpressionString += logicalOperatorStmtRefSeparator
				}
				// ... and append to logical expression column string
				logicalExpressionString += fmt.Sprint(ops)
				// Leading bracket
				logicalExpressionString += logicalCombinationLeft

				fmt.Println("Target node IDs: ", targetID)
				// Add trailing bracket and column ref (to be reviewed)
				logicalExpressionString += targetID + logicalCombinationRight
			}

		}
	}
	// Return generated string
	return logicalExpressionString, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Generates Google Sheets output from map of categorized statement elements, as well as header columns as indices.
Optionally writes to file.
 */
func GenerateGoogleSheetsOutput(statementMap []map[string]string, headerCols []string, headerColsNames []string, filename string) (string, tree.ParsingError) {
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

	// Generate header column row based on names
	output := prefix
	for _, v := range headerColsNames {
		output += v + separator
	}
	output += suffix

	// Generate all entry rows
	for _, entry := range statementMap {
		// Create new row with Google Sheets syntax and leading ' to ensure text interpretation of ID
		output += prefix + stmtIdPrefix
		// Reconstruct based on header column order
		for _, header := range headerCols {
			if entry[header] == "" {
				// if entry for given header is empty, add space
				output += " " + separator
			} else {
				// else add entry value
				output += entry[header] + separator
			}
		}
		// Trim last separator
		output = output[:len(output)-1]
		// Append Google Sheets-specific suffix to complete row
		output += suffix
	}

	// Write file
	if filename != "" {
		err := WriteToFile(filename, output)
		if err != nil {
			return output, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_WRITE, ErrorMessage: err.Error()}
		}
	}

	return output, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Generates Google Sheets tabular output for a given parsed statement, with a given statement ID.
Generates all substatements and logical combination linkages in Google Sheets output format.
Additionally returns array of statement entries, header symbols and corresponding header symbol names.
If filename is provided, the result is printed to the corresponding file.
 */
func GenerateGoogleSheetsOutputFromParsedStatement(statement tree.Statement, stmtId string, filename string) (string, []map[string]string, []string, []string, tree.ParsingError) {
	log.Println("Step: Extracting leaf arrays")
	// Retrieve leaf arrays from generated tree (alongside frequency indications for components)
	leafArrays, componentRefs := statement.GenerateLeafArrays()

	log.Println("Generated leaf arrays: ", leafArrays, " component: ", componentRefs)

	log.Println("Step: Generate permutations of leaf arrays (atomic statements)")
	// Generate all permutations of logically-linked components to produce statements
	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", nil, nil, nil, err
	}

	log.Println("Generated permutations: ", res)

	log.Println("Step: Generate logical operators for atomic statements")
	// Extract logical operator links
	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	log.Println("Links:", links)

	log.Println("Step: Generate tabular output")
	// Export in Google Sheets format
	statementMap, statementHeaders, statementHeaderNames, err := generateTabularStatementOutput(res, componentRefs, links, stmtId)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", nil, nil, nil, err
	}

	// Create Google Sheets output based on generated map, alongside header names as output
	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeaderNames, filename)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return output, statementMap, statementHeaders, statementHeaderNames, err
	}

	// Outfile will only be written if filename is specified
	/*if filename != "" {
		log.Println("Step: Write output to file")
		// Write to file
		errWrite := WriteToFile(filename, output)
		if errWrite != nil {
			// Wrap into own error, alongside generated (but not written) output
			return output, nil, nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_WRITE, ErrorMessage: errWrite.Error()}
		}
	}*/

	return output, statementMap, statementHeaders, statementHeaderNames, err
}

/*
Generates IG 2.0 header row and appends it to given string based on component frequency input. It further returns a slice
containing header information.
 */
func generateHeaderRow(stringToAppendTo string, componentFrequency map[string]int, separator string) (string, []string, []string) {

	// Header symbols to be returned for later use (used in logical operators)
	headerSymbols := []string{}
	// Header symbol names to be returned for column header construction
	headerSymbolsNames := []string{}
	// Iterate through component reference map
	for _, symbol := range tree.IGComponentSymbols {
		i := 0
		// Print headers as often as occurring in input file (stmtCt.e., one header for each column)
		for i < componentFrequency[symbol] {
			// Store header name for column name construction
			headerSymbolsName := tree.IGComponentSymbolNameMap[symbol]
			// Store symbols for columns including indices in order of occurrence for use in logical operators
			headerSymbol := symbol
			// Introduce indices if multiple of the same component
			if componentFrequency[symbol] > 1 {
				// Append suffix for header string
				headerSymbolsName += indexSymbol + strconv.Itoa(i + 1)
				// Append suffix for cached header IDs (for logical operators)
				headerSymbol += indexSymbol + strconv.Itoa(i + 1)
			}
			// Store key for header used in logical operators
			headerSymbols = append(headerSymbols, headerSymbol)
			headerSymbolsNames = append(headerSymbolsNames, headerSymbolsName)
			// Append full header names to string
			stringToAppendTo += headerSymbolsName
			stringToAppendTo += separator
			i++
		}
	}
	// Cut off last separator
	stringToAppendTo = stringToAppendTo[0:len(stringToAppendTo)-len(separator)]
	// Return generated string as well as symbol map and mapped names
	return stringToAppendTo, headerSymbols, headerSymbolsNames
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
		if firstKey != nil {
			leaves := firstKey.GetSyntheticRootNode().GetLeafNodes()
			if len(leaves) > 0 {
				nodesKeys = leaves[0]
			} else {
				fmt.Println("No component keys to iterate over for logical relationships")
			}
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
				if len(linkedElement) > 0 {
					fmt.Println("Testing other node: ", otherNode, " with elements ", linkedElement)
					// find operator
					res, ops, err := tree.FindLogicalLinkage(statement[componentIdx], otherNode)
					if err.ErrorCode != tree.TREE_NO_ERROR {
						errorMsg := fmt.Sprint("Error when parsing retrieving operator linkages: ", err.ErrorMessage)
						log.Println(errorMsg)
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
						logicalExpressionString += logicalCombinationLeft
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
						logicalExpressionString += stmtsRefs + logicalCombinationRight + logicalOperatorSeparator
					}
					fmt.Println("Added logical relationships for value", otherNode, ", elements:", logicalExpressionString)
				} else {
					fmt.Println("Did not find target links for", otherNode, "- did not add logical operator links for component")
				}
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
