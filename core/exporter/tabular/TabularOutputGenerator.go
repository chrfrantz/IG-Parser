package tabular

import (
	"IG-Parser/core/parser"
	"IG-Parser/core/shared"
	"IG-Parser/core/tree"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/*
This file contains the tabular output generation functionality, including
associated constants relevant for output formatting.
*/

// Separator for main statement ID (e.g., 123) and suffix for introduced substatement (e.g., .1, i.e., 123.1)
const stmtIdSeparator = "."

// Separator for logical operators in enumerations of statement references (e.g., OR[650.1,650.2, ...])
const logicalOperatorStmtRefSeparator = ","

// Separator for referenced statements in cell (e.g., multiple references to activation conditions, i.e., {65}.1,{65}.2)
const componentStmtRefSeparator = ","

// Symbol separating component symbol and indices (e.g., Bdir vs. Bdir_1, Bdir_2, etc.)
const indexSymbol = "_"

// Statement ID prefix to ensure interpretation as text field in Google Sheets (does not remove trailing zeroes)
const stmtIdPrefix = "'"

// Separator for logical operator expressions (e.g., OR[650.1,650.2]|AND[123.1,123.2])
const logicalOperatorSeparator = ";"

// Suffix separator for extrapolated statements (e.g., Base ID (1.1) with operator . extrapolated statements would be 1.1.1, 1.1.2, etc.)
const extrapolatedStatementSuffixSeparator = "."

// Left bracket for logical combination expressions
const logicalCombinationLeft = parser.LEFT_BRACKET

// Right bracket for logical combination expressions
const logicalCombinationRight = parser.RIGHT_BRACKET

// Left brace surrounding identifier for component-level nested statements
const componentNestedLeft = parser.LEFT_BRACE

// Right brace surrounding identifier for component-level nested statements
const componentNestedRight = parser.RIGHT_BRACE

// Column identifier for Statement ID
const stmtIdColHeader = "Statement ID"

// Column identifier for Original Statement input
const stmtOriginalStatementHeader = "Original Statement"

// Column identifier for IG Script-encoded statement
const stmtIgScriptHeader = "IG Script Encoding"

// Column identifier for logical linkage of components
const logLinkColHeaderComps = "Logical Linkage (Components)"

// Column identifier for logically linked statements (not just components)
const logLinkColHeaderStmts = "Logical Linkage (Statements)"

// Default separator used for header row generation
var CellSeparator = "|"

// Default separator for multiple items within cell
const cellValueSeparator = ","

// Structure referencing ID (based on input ID), along with (nested) statement to be decomposed
type IdentifiedStmt struct {
	ID         string     // Generated ID as provided in output
	NestedStmt *tree.Node // Single statement wrapped in Node
}

/*
Generates array of statement maps corresponding to identified elements format. Includes parsing of nested statements.
Consider the specification of INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT variable to indicate whether shared elements
are to be included in output.
Input:
  - Atomic statements with corresponding node references [statement][node references]
  - Input statement annotations (i.e., of statement, not components)
  - String representing linkage to other statements (not components) produced based on extrapolated component pairs
  - Map with component name as key and corresponding number of columns in input stmts (i.e., same component can have
    values across multiple columns)
  - References to entries for given nodes as indicated by logical operators, and used to produce corresponding linkages
    (e.g., AND[row1, row2, etc.])
  - ID to be used as prefix for generation of substatement IDs (e.g., ID 5 produces substatements 5.1, 5.2, etc.)
  - headerSeparator used for generation of header row (e.g., ";")
  - outputType allows for the specification of target output type to introduce necessary preprocessing as part of the matrix generation (e.g., prefixing quotes).
    Valid output types are defined in TabularOutputGeneratorConfig (e.g., #OUTPUT_TYPE_GOOGLE_SHEETS, etc.)
  - printHeaders indicates whether header row is included in output.

Output:
- Array of statement entry maps (i.e., values for each component in given statement, i.e., [statement]map[component]componentValue)
- Array of header symbols (used for component linkage references)
- Array of header symbols names (for human-readable header construction)
*/
func generateStatementMatrix(stmts [][]*tree.Node, annotations interface{}, stmtLogicalLinks string, componentFrequency map[string]int, logicalLinks []map[*tree.Node][]string, stmtId string, headerSeparator string, outputType string, printHeaders bool) ([]map[string]string, []string, []string, tree.ParsingError) {

	if headerSeparator == "" {
		return nil, nil, nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_MISSING_SEPARATOR_VALUE,
			ErrorMessage: "Value for separator symbol is invalid."}
	}

	// Caches column header symbols by component index for reuse in logical operator construction
	headerSymbols := []string{}
	// Caches column header names associated with symbols for human-readable header construction
	headerSymbolsNames := []string{}

	sepErr := tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}

	if ProduceDynamicOutput() {
		// Generate headers based on parsed statement input
		if componentFrequency != nil && len(componentFrequency) != 0 {
			// Iterate through header frequencies and create header row
			_, headerSymbols, headerSymbolsNames, sepErr = generateHeaderRow("", componentFrequency, headerSeparator)
			if sepErr.ErrorCode != tree.PARSING_NO_ERROR {
				return nil, nil, nil, sepErr
			}
		}
	} else {
		// Generate static headers not taking frequencies of components into account
		for k, v := range componentFrequency {
			if v != 1 {
				log.Println("Found component frequency > 1 for component", k)
			}
		}

		Println("Providing output based on fixed structure")

		// Iterate through header frequencies and create header row
		_, headerSymbols, headerSymbolsNames, sepErr = generateHeaderRow("", GetStaticTabularOutputSchema(), headerSeparator)
		if sepErr.ErrorCode != tree.PARSING_NO_ERROR {
			return nil, nil, nil, sepErr
		}
	}

	Println("Generated Header Symbols: ", headerSymbols)
	Println("Generated Header Symbol Names: ", headerSymbolsNames)

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

		Println("Statement (to be parsed in tabular form), ID:", stmtCt, ":", statement)

		// Create new entry with individual ID

		// Add statement ID for specific instance
		subStmtId := stmtId
		// Only generate subIDs if indeed multiple statements (i.e., some form of nesting)
		if len(stmts) > 1 {
			subStmtId = generateStatementIDint(stmtId, stmtCt+1)
		}
		// Add statement ID to entryMap
		entryMap[stmtIdColHeader] = subStmtId
		// String linking all logical operators for a given row
		logicalValue := ""

		// Include statement-level annotations if activated and existing in input
		if include_ANNOTATIONS && annotations != nil {
			entryMap[tree.STATEMENT_ANNOTATION] = annotations.(string)
		}
		// Iterate over component index (i.e., column) covering conventional components
		for componentIdx := range statement {
			// Append element value as output for given cell
			if statement[componentIdx].IsEmptyOrNilNode() {
				// Empty entry - don't add anything
				Println("Found empty node for component", fmt.Sprint(statement[componentIdx].GetComponentName()))
			} else if statement[componentIdx].HasPrimitiveEntry() {
				Println("Found primitive entry in component", fmt.Sprint(statement[componentIdx].GetComponentName()), ", Entry: ", statement[componentIdx])
				// Regular leaf entry (i.e., component) - unless IG Core coding is intended (which implies collapsed nested components)

				// Provide default values for left and right elements (potentially used hereafter)
				leftString := ""
				rightString := ""
				// If shared elements are to be included (based on configuration, extract those ...)
				if IncludeSharedElementsInTabularOutput() {
					// Prepare left and right shared elements by stringifying
					leftString = shared.StringifySlices(statement[componentIdx].GetSharedLeft())
					// ... but don't append whitespace just yet - depends on matching of shared strings later
					rightString = shared.StringifySlices(statement[componentIdx].GetSharedRight())
					if rightString != "" {
						// Add preceding whitespace
						rightString = " " + rightString
					}
				}

				// Prepare value for entry
				entryVal := strings.Builder{}
				// Indicates whether values within cell should be comma-separated
				skipSeparator := false

				// CHECK FOR CELL SEPARATORS

				// Check for preceding shared elements, and suppress left element if needed
				if leftString != "" {
					if ProduceDynamicOutput() {
						// Dynamic variant
						// Check whether value exists in cell
						if len(entryMap[headerSymbols[componentIdx]]) > 0 &&
							strings.HasSuffix(entryMap[headerSymbols[componentIdx]], leftString) {
							// Suppress left shared element if identical with shared right one on existing value
							// but add whitespace to link to previous value
							entryVal.WriteString(" ")
							entryVal.WriteString(statement[componentIdx].Entry.(string))
							entryVal.WriteString(rightString)
							// Skip comma separation
							skipSeparator = true
						} else {
							// Regular sharedLeft, whitespace + value sharedRight concatenation
							entryVal.WriteString(leftString)
							entryVal.WriteString(" ")
							entryVal.WriteString(statement[componentIdx].Entry.(string))
							entryVal.WriteString(rightString)
						}
					} else {
						// Static variant
						// Check whether value exists in cell
						if len(entryMap[statement[componentIdx].GetComponentName()]) > 0 &&
							strings.HasSuffix(entryMap[statement[componentIdx].GetComponentName()], leftString) {
							// Suppress left shared element if identical with shared right one on existing value
							// but add whitespace to link to previous value
							entryVal.WriteString(" ")
							entryVal.WriteString(statement[componentIdx].Entry.(string))
							entryVal.WriteString(rightString)
							// Skip comma separation
							skipSeparator = true
						} else {
							// Regular sharedLeft, whitespace + value sharedRight concatenation
							entryVal.WriteString(leftString)
							entryVal.WriteString(" ")
							entryVal.WriteString(statement[componentIdx].Entry.(string))
							entryVal.WriteString(rightString)
						}
					}
				} else {
					// Create regular entry (without left shared value, since that will be empty)
					entryVal.WriteString(statement[componentIdx].Entry.(string))
					entryVal.WriteString(rightString)
				}

				// Determine whether cell separation is used
				effectiveCellSeparator := ""
				if !skipSeparator {
					effectiveCellSeparator = cellValueSeparator
				}

				// SPECIAL SYMBOLS (APPLICATION-SPECIFIC) --> SYMBOLS THAT REQUIRE SUBSTITUTION
				// Substitute symbols before producing output (e.g., " with ')
				// TODO: Review for further symbols
				entryValStr := shared.EscapeSymbolsForExport(entryVal.String())

				// HANDLE OUTPUT-SPECIFIC MODIFICATIONS

				// TODO: Google Sheets specific - consider adaption to support further formats
				if outputType == OUTPUT_TYPE_GOOGLE_SHEETS {
					// Duplicate leading ' for proper Google Sheets parsing
					if len(entryValStr) > 0 && entryValStr[0:1] == "'" {
						entryValStr = "'" + entryValStr
					}
				}

				// ADDING ACTUAL ENTRY

				if ProduceDynamicOutput() {
					// Dynamic variant
					// Save entry value into entryMap for given statement and component column
					if len(entryMap[headerSymbols[componentIdx]]) > 0 {
						// Add separator for cell values
						b := strings.Builder{}
						b.WriteString(entryMap[headerSymbols[componentIdx]])
						b.WriteString(effectiveCellSeparator)
						b.WriteString(entryValStr)
						entryMap[headerSymbols[componentIdx]] = b.String()
					} else {
						// First value, hence no separator needed
						entryMap[headerSymbols[componentIdx]] = entryValStr
					}
				} else {
					// Static variant
					Println("Looking for component entry for", statement[componentIdx].GetComponentName())
					// Save entry for a given field matched based on node's component type
					if len(entryMap[statement[componentIdx].GetComponentName()]) > 0 {
						// Add separator for cell values
						b := strings.Builder{}
						b.WriteString(entryMap[statement[componentIdx].GetComponentName()])
						b.WriteString(effectiveCellSeparator)
						b.WriteString(entryValStr)
						entryMap[statement[componentIdx].GetComponentName()] = b.String()
					} else {
						// First value, hence no separator needed
						entryMap[statement[componentIdx].GetComponentName()] = entryValStr
					}
				}
				Println("Added entry ", entryValStr)

				// PRIVATE NODES

				// For static output, consider private nodes
				if !ProduceDynamicOutput() && statement[componentIdx].HasPrivateNodes() {
					for _, privateNodeValue := range statement[componentIdx].PrivateNodeLinks {

						// If iterated private node is of type statement (single nested statement),
						// embed into node slice to enforce corresponding processing ...
						if reflect.TypeOf(privateNodeValue.Entry) == reflect.TypeOf(&tree.Statement{}) {
							tempNode := []*tree.Node{&tree.Node{Entry: privateNodeValue.Entry.(*tree.Statement)}}
							// Reassign node-slice-embedded statement to iterated item
							privateNodeValue.Entry = tempNode
						}

						// Nested private properties
						if reflect.TypeOf(privateNodeValue.Entry) == reflect.TypeOf([]*tree.Node{}) {
							for _, v := range privateNodeValue.Entry.([]*tree.Node) {

								if ProduceIGExtendedOutput() {
									// IG Extended --> expand into nested statements referenced here, and appended as part of other nested statements
									stmtRef := ""

									// Check for existing private nodes ...
									existing := entryMap[privateNodeValue.GetComponentName()+tree.REF_SUFFIX]
									if len(existing) > 0 {
										// ... and append if necessary
										existing += cellValueSeparator
									}

									// Generate (and cache) new statement ID for nested statement
									componentNestedStmtsMap, nestedStatementIdx, componentNestedStmts, stmtRef = generateNewNestedStatementID(v, componentNestedStmtsMap, nestedStatementIdx, componentNestedStmts, stmtId)

									// Add nested statement ID to existing value
									existing += stmtRef

									// Assign ID to reference field
									entryMap[privateNodeValue.GetComponentName()+tree.REF_SUFFIX] = existing
								} else {
									// IG Core --> print flat string representation of nested statement

									// Check for existing private nodes ...
									existing := entryMap[privateNodeValue.GetComponentName()]
									if len(existing) > 0 {
										// ... and append if necessary
										existing += cellValueSeparator
									}

									// Print flat content
									existing += v.StringFlat()

									// (Re)Assign to entry to be output
									entryMap[privateNodeValue.GetComponentName()] = existing
								}

							}
							// processing of registered nested statement will occur later

							// non-nested private nodes ...
						} else if reflect.TypeOf(privateNodeValue.Entry) == reflect.TypeOf("") {

							// Check for existing private nodes ...
							existing := entryMap[privateNodeValue.GetComponentName()]
							if len(existing) > 0 {
								// ... and append if necessary
								existing += cellValueSeparator
							}

							// Add actual primitive value
							existing += privateNodeValue.Entry.(string)

							// (Re)Assign to entry to be output
							entryMap[privateNodeValue.GetComponentName()] = existing
						}
					}
					Println("Added private nodes to given output node")
				}

				// ANNOTATIONS

				// For static output, consider annotations (if activated)
				if !ProduceDynamicOutput() && IncludeAnnotations() && statement[componentIdx].HasAnnotations() {

					// Check for existing annotations ...
					existing := entryMap[statement[componentIdx].GetComponentName()+tree.ANNOTATION]
					if len(existing) > 0 {
						// ... and append if necessary
						existing += cellValueSeparator
					}
					// Add actual value
					existing += statement[componentIdx].GetAnnotations().(string)
					// (Re)Assign to entry to be output
					entryMap[statement[componentIdx].GetComponentName()+tree.ANNOTATION] = existing
				}

				Println("Entry (after adding primitive entry):", entryMap)
			} else {
				// Nested statements are stored for later processing, but assigned IDs and references added to calling row
				Println("Found complex entry (nested statement) in component: " + fmt.Sprint(statement[componentIdx].GetComponentName()))
				Println("Complex entry (nested statement):" + fmt.Sprint(statement[componentIdx]))
				Println("Complex entry's annotations:" + fmt.Sprint(statement[componentIdx].Annotations))
				// Check for statement combination (i.e., node combination)

				// Add entry to array (assuming single nested statement)
				entryVals := []*tree.Node{statement[componentIdx]}

				// Check if combination contained; if so, flatten, and override
				if entryVals[0].IsCombination() {
					Println("Detected statement combination")
					// If combination of statements, retrieve all elements
					combStmts := entryVals[0].GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)
					// Flatten array and override entry values for iteration
					entryVals = tree.Flatten(combStmts)
					Println("Flattened combination:", entryVals)
				} else {
					Println("Detected individual nested statement")
				}

				// Iterate over all nested statements
				for _, entryVal := range entryVals {

					// Keep track whether entry is last in iteration -
					// which is relevant for inclusion of logical operators in combined statements
					last := false
					if entryVal == entryVals[len(entryVals)-1] {
						last = true
					}

					idToReferenceInCell := ""

					if ProduceIGExtendedOutput() {

						// Generate (and cache) new statement ID for nested statement
						componentNestedStmtsMap, nestedStatementIdx, componentNestedStmts, idToReferenceInCell = generateNewNestedStatementID(entryVal, componentNestedStmtsMap, nestedStatementIdx, componentNestedStmts, stmtId)
					}

					if ProduceDynamicOutput() {
						// Dynamic version
						// Save entry into entryMap for calling row
						if entryMap[headerSymbols[componentIdx]] != "" &&
							// Suppress separator if preceding element is a statement
							!strings.HasSuffix(entryMap[headerSymbols[componentIdx]], logicalCombinationRight+" ") {

							// Add separator if already an entry
							entryMap[headerSymbols[componentIdx]] += componentStmtRefSeparator
						}

						if ProduceIGExtendedOutput() {
							// Add nested statement reference (IG Extended)
							entryMap[headerSymbols[componentIdx]] += idToReferenceInCell
						} else {
							// IG Core output without nesting

							// Append flat string representation of nested statements
							entryMap[headerSymbols[componentIdx]] += entryVal.StringFlat()

							// Add logical operator if not last entry (and parent not empty otherwise)
							if !last && entryVal.Parent != nil {
								b := strings.Builder{}
								b.WriteString(" ")
								b.WriteString(logicalCombinationLeft)
								b.WriteString(entryVal.Parent.LogicalOperator)
								b.WriteString(logicalCombinationRight)
								b.WriteString(" ")
								entryMap[headerSymbols[componentIdx]] += b.String()
							}
						}
					} else {
						// Static version
						Println("Linking substatement ID to component", statement[componentIdx],
							"Name: ", statement[componentIdx].GetComponentName())
						// Save entry into entryMap for calling row
						if entryMap[statement[componentIdx].GetComponentName()+tree.REF_SUFFIX] != "" &&
							// Suppress separator if preceding element is a statement
							!strings.HasSuffix(entryMap[statement[componentIdx].GetComponentName()+tree.REF_SUFFIX], logicalCombinationRight+" ") {

							// Add separator if already an entry
							entryMap[statement[componentIdx].GetComponentName()+tree.REF_SUFFIX] += componentStmtRefSeparator
						}

						if ProduceIGExtendedOutput() {
							// Add nested statement reference (IG Extended)

							// Append current value in any case
							entryMap[statement[componentIdx].GetComponentName()+tree.REF_SUFFIX] += idToReferenceInCell
						} else {
							// IG Core output without nesting

							// Append flat string representation of nested statements
							entryMap[statement[componentIdx].GetComponentName()+tree.REF_SUFFIX] += entryVal.StringFlat()

							// Add logical operator if not last entry (and parent not empty otherwise)
							if !last && entryVal.Parent != nil {
								b := strings.Builder{}
								b.WriteString(" ")
								b.WriteString(logicalCombinationLeft)
								b.WriteString(entryVal.Parent.LogicalOperator)
								b.WriteString(logicalCombinationRight)
								b.WriteString(" ")
								entryMap[statement[componentIdx].GetComponentName()+tree.REF_SUFFIX] += b.String()

							}
						}
					}
				}
			}
			Println("Source/calling node (for nested statement): ", statement[componentIdx])

			// Process component-level logical linkage

			// Now generate logical links expression corresponding to particular entry (component index in statement instance)
			logicalValue, errorVal = generateLogicalLinksExpressionForGivenComponentValue(logicalValue, statement,
				componentIdx, headerSymbols, logicalLinks, stmtId)
			if errorVal.ErrorCode != tree.PARSING_NO_ERROR {
				return nil, nil, nil, errorVal
			}
			Println("Logical expressions for atomic statement:", logicalValue)
		}

		Println("Entries (complete row - before adding logical operators and nested statements):", entryMap)

		// Append the logical expression at the end of each row
		if logicalValue != "" {
			// Add to entryMap
			entryMap[logLinkColHeaderComps] = logicalValue
			// Reset for next round
			logicalValue = ""
		}
		// Append the logical expression for linkage to other statements (extrapolated statements) at the end of each top-level statement row
		if stmtLogicalLinks != "" {
			// Add to EntryMap
			entryMap[logLinkColHeaderStmts] = stmtLogicalLinks
		}
		// Add to entries map for statement to map for all statements (collection for return)
		entriesMap = append(entriesMap, entryMap)
	}

	Println("Component-level nested statements to be decomposed: " + fmt.Sprint(componentNestedStmts))
	for _, val := range componentNestedStmts {

		Println("Nested Statement to parse, ID:", val.ID, ", Annotations:", val.NestedStmt.Annotations, ", Stmt:", val.NestedStmt)

		Println("Parsing nested statement ...")
		// Parse individual nested statements on component level in order to attach those to main output
		nestedTabularResult := GenerateTabularOutputFromParsedStatement(val.NestedStmt, nil, "", "", val.NestedStmt.Annotations, val.ID, "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, headerSeparator, outputType, printHeaders, ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
		if nestedTabularResult.Error.ErrorCode != tree.PARSING_NO_ERROR {
			return nil, nil, nil, nestedTabularResult.Error
		}

		// Add linkages between statements (statement-level combinations)

		// Determine linkages to fellow nested statements
		stmtLinksString, err := generateLogicalLinksExpressionForStatements(val.NestedStmt, componentNestedStmts)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return nil, nil, nil, err
		}
		Println("Logical links to other statements: ", stmtLinksString)

		// Add identified linkages to nestedMap (i.e., for all atomic statements on *same level* (don't overwrite lower levels!)
		if stmtLinksString != "" {
			// Check for nesting prefix first (to prevent overriding statements with higher-level nesting
			// (e.g., Bdir{ A(actor) Cac{ A(actor2) ...}; here only atomic statements with the same nesting level
			// (i.e., Bdir{ A(actor) } should be appended
			// This is identified based on the uniform prefix for such statements that will all conclude with "}." followed by unique ID
			// Example: Top level 123.1 --> no nesting, but same level; {123.1}.1 --> first nesting level; {{123.1}.1}.1 --> second nesting level
			nestedIdPrefix := ""
			if strings.Contains(val.ID, "}.") {
				nestedIdPrefix = val.ID[:strings.LastIndex(val.ID, "}.")]
			}
			for i := range nestedTabularResult.StatementMap {
				// If prefix exists and the prefix applies also to the atomic statement of concern ...
				if nestedIdPrefix != "" && strings.HasPrefix(nestedTabularResult.StatementMap[i][stmtIdColHeader], nestedIdPrefix) {
					// ... then append ...
					if nestedTabularResult.StatementMap[i][logLinkColHeaderStmts] != "" {
						nestedTabularResult.StatementMap[i][logLinkColHeaderStmts] =
							nestedTabularResult.StatementMap[i][logLinkColHeaderStmts] + componentStmtRefSeparator + stmtLinksString
					} else {
						// ... or assign outright
						nestedTabularResult.StatementMap[i][logLinkColHeaderStmts] = stmtLinksString
					}
					// else do not assign at all (since statements are likely on lower nesting level
				} else if nestedIdPrefix == "" {
					// ... else if no prefix exists (i.e., no nesting), simply assign to all statements (must be on same level)
					nestedTabularResult.StatementMap[i][logLinkColHeaderStmts] = stmtLinksString
				}
			}
		}

		// Add Logical linkage header if not already existing
		nestedTabularResult.HeaderSymbols = addElementIfNotExisting(logLinkColHeaderStmts, nestedTabularResult.HeaderSymbols)
		nestedTabularResult.HeaderNames = addElementIfNotExisting(logLinkColHeaderStmts, nestedTabularResult.HeaderNames)

		// Add nested entries to top-level list
		entriesMap = append(entriesMap, nestedTabularResult.StatementMap...)

		// Merge headers to consider nested ones
		headerSymbols = tree.MergeSlices(headerSymbols, nestedTabularResult.HeaderSymbols, indexSymbol)
		// Merge header names to consider nested ones
		headerSymbolsNames = tree.MergeSlices(headerSymbolsNames, nestedTabularResult.HeaderNames, indexSymbol)
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
Generate new nested statement ID for given nested input statement, and adds it to index of nested statements.
Returns extended structures containing new statement.
*/
func generateNewNestedStatementID(entryVal *tree.Node, componentNestedStmtsMap map[*tree.Node]string, nestedStatementIdx int, componentNestedStmts []IdentifiedStmt, stmtId string) (map[*tree.Node]string, int, []IdentifiedStmt, string) {

	idToReferenceInCell := ""

	// Retrieve ID of already identified statements ...
	if nestedStmtID, ok := componentNestedStmtsMap[entryVal]; ok {
		// Prepare reference to be saved
		idToReferenceInCell = nestedStmtID
	} else {
		// ... else create new one
		// Generate ID for component-level nested statement
		b := strings.Builder{}
		b.WriteString(componentNestedLeft)
		b.WriteString(stmtId)
		b.WriteString(componentNestedRight)
		b.WriteString(stmtIdSeparator)
		b.WriteString(strconv.Itoa(nestedStatementIdx))
		nestedStmtId := b.String()
		Println("Generated ID for nested statement:", nestedStmtId)
		// Add component-level nested statement
		componentNestedStmts = append(componentNestedStmts,
			IdentifiedStmt{nestedStmtId, entryVal})
		// Add newly identified nested statement to lookup index
		componentNestedStmtsMap[entryVal] = nestedStmtId
		// Increase index for component-level nested statements (for next round)
		nestedStatementIdx++
		// Prepare reference to to-be component-level nested statements to output
		idToReferenceInCell = nestedStmtId
		Println("Parsing: Added nested statement (ID:", nestedStmtId, ", Annotations:", entryVal.Annotations, ", Val:", entryVal, ")")
	}

	return componentNestedStmtsMap, nestedStatementIdx, componentNestedStmts, idToReferenceInCell

}

/*
Resolves all logical linkages to other statements and returns those as compound logical expression (e.g., [AND][{65}.1],[AND][{65}.2])
*/
func generateLogicalLinksExpressionForStatements(sourceStmt *tree.Node, allNestedStmts []IdentifiedStmt) (string, tree.ParsingError) {
	builder := strings.Builder{}

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
				if CollapseOperators() {
					Println("Collapsing adjacent AND, bAND and wAND operators ...")
					// Collapse adjacent AND operators
					ops = tree.CollapseAdjacentOperators(ops, []string{tree.AND, tree.SAND_BETWEEN_COMPONENTS, tree.SAND_WITHIN_COMPONENTS})
				}

				Println("Node has linkage ", ops)

				if builder.String() != "" {
					builder.WriteString(logicalOperatorStmtRefSeparator)
				}
				// ... and append to logical expression column string
				builder.WriteString(fmt.Sprint(ops))
				// Leading bracket
				builder.WriteString(logicalCombinationLeft)

				Println("Target node IDs: ", targetID)
				// Add trailing bracket and column ref (to be reviewed)
				builder.WriteString(targetID)
				builder.WriteString(logicalCombinationRight)
			}

		}
	}
	// Return generated string
	return builder.String(), tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Generates final tabular output based on input statement matrices, alongside header information and optionally prints it to file.

Input includes
statement matrix (i.e., parsed institutional statements decomposed into matrix structure),
Original statement input (for consideration in output - parameterized via printOriginalStatement),
IG Script input (for consideration in output - parameterized via printIgScript),
header symbols (ordered IG component symbols associated with matrix columns),
header names (ordered and associated with header symbols),
row prefix (prefix for each row - to accommodate specific output formats),
stmtIdPrefix (prefix for statement ID to ensure parsing of output as text),
row suffix (suffix for each row - to accommodate specific output formats),
separator used to separate individual cells per row,
filename the output should be printed to (should be "" if no output is to be printed),
printHeaders to indicate whether the header row is to be included in output,
printOriginalStatement to indicate the inclusion of the original statement input in the generated output (for options see tabular.ORIGINAL_STATEMENT_INCLUSION_OPTIONS).
printIgScript to indicate the inclusion of the original IG Script input in the generated output (for options see tabular.IG_SCRIPT_INCLUSION_OPTIONS).

Returns string containing flat output as well as potential parsing error
*/
func printTabularOutput(statementMap []map[string]string, originalStatement string, igScriptInput string, headerCols []string, headerColsNames []string, rowPrefix string, stmtIdPrefix string, rowSuffix string, separator string, filename string, overwrite bool, printHeaders bool, printOriginalStatement string, printIgScript string) (string, tree.ParsingError) {

	// Prepare builder
	builder := strings.Builder{}

	if printHeaders {
		// Generate header column row based on names
		builder.WriteString(rowPrefix)
		// Iterate through actual content
		for _, v := range headerColsNames {
			builder.WriteString(v)
			builder.WriteString(separator)
			// Immediately write optional Original Statement and IG Script headers (if input is not empty)
			if v == stmtIdColHeader {
				// Include column for Original Statement output if some form of output is selected
				if printOriginalStatement != ORIGINAL_STATEMENT_OUTPUT_NONE {
					// Column for Original Statement content
					builder.WriteString(stmtOriginalStatementHeader)
					builder.WriteString(separator)
				}
				// Include column for IG Script output if some form of output is selected
				if printIgScript != IG_SCRIPT_OUTPUT_NONE {
					// Column for IG Script content
					builder.WriteString(stmtIgScriptHeader)
					builder.WriteString(separator)
				}
			}
		}
		builder.WriteString(rowSuffix)
	}

	// Generate all entry rows
	for i, entry := range statementMap {
		// Create new row with given syntax and potential statement ID prefix (e.g., ' to ensure text interpretation of ID).
		builder.WriteString(rowPrefix)
		builder.WriteString(stmtIdPrefix)
		// Reconstruct based on header column order
		for _, header := range headerCols {
			if entry[header] == "" {
				// if entry for given header is empty, add space
				builder.WriteString(" ")
				builder.WriteString(separator)
			} else {
				// else add entry value
				builder.WriteString(entry[header])
				builder.WriteString(separator)
			}
			// Add Original Statement column entry (if not deactivated)
			if header == stmtIdColHeader && printOriginalStatement != ORIGINAL_STATEMENT_OUTPUT_NONE {
				switch printOriginalStatement {
				case ORIGINAL_STATEMENT_OUTPUT_FIRST_ENTRY:
					// Only print for the first actual entry (index 0)
					if i == 0 {
						// Write full entry with content in first entry row
						builder.WriteString(originalStatement)
						builder.WriteString(separator)
					} else {
						// ... else just add whitespace and separator
						builder.WriteString(" ")
						builder.WriteString(separator)
					}
				case ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES:
					// Print in every row
					builder.WriteString(originalStatement)
					builder.WriteString(separator)
				default:
					log.Println("Invalid Original Statement output specification ('" + printOriginalStatement + "'). Output suppressed.")
				}
			}
			// Add IG Script column entry (if not deactivated)
			if header == stmtIdColHeader && printIgScript != IG_SCRIPT_OUTPUT_NONE {
				switch printIgScript {
				case IG_SCRIPT_OUTPUT_FIRST_ENTRY:
					// Only print for the first actual entry (index 0)
					if i == 0 {
						// Write full entry with content in first entry row
						builder.WriteString(igScriptInput)
						builder.WriteString(separator)
					} else {
						// ... else just add whitespace and separator
						builder.WriteString(" ")
						builder.WriteString(separator)
					}
				case IG_SCRIPT_OUTPUT_ALL_ENTRIES:
					// Print in every row
					builder.WriteString(igScriptInput)
					builder.WriteString(separator)
				default:
					log.Println("Invalid IG Script output specification ('" + printIgScript + "'). Output suppressed.")
				}
			}
		}
		// Append format-specific row suffix
		builder.WriteString(rowSuffix)
	}

	// Write file
	if filename != "" {
		err := WriteToFile(filename, builder.String(), overwrite)
		if err != nil {
			return builder.String(), tree.ParsingError{ErrorCode: tree.PARSING_ERROR_WRITE, ErrorMessage: err.Error()}
		}
	}

	return builder.String(), tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Generates CSV output from map of categorized statement elements, IG Script input (included in output if parameterized via printIgScriptInput),
as well as header columns (symbols and names) for output generation.
Further requires column header names for output generation, alongside specification of separator symbol.
Optionally writes to file (if filename is provided), with option to overwrite existing files,
as well as options to print header lines, and IG Script input as part of the output (for options see tabular.IG_SCRIPT_INCLUSION_OPTIONS).
*/
func generateCSVOutput(statementMap []map[string]string, originalStatement string, igScriptInput string, headerCols []string, headerColsNames []string, separator string, filename string, overwrite bool, printHeaders bool, printOriginalStatement string, printIgScriptInput string) (string, tree.ParsingError) {

	// Linebreak at the end of each entry
	suffix := "\n"

	// Delegate actual printing
	return printTabularOutput(statementMap, originalStatement, igScriptInput, headerCols, headerColsNames, "", stmtIdPrefix, suffix, separator, filename, overwrite, printHeaders, printOriginalStatement, printIgScriptInput)
}

/*
Generates Google Sheets output from map of categorized statement elements, IG Script input (for consideration in output - parameterized via printOriginalStatement and printIgScriptInput),
as well as header columns (symbols and names) for output generation.
Further requires column header names for output generation, alongside specification of separator symbol.
Optionally writes to file (if filename is provided), with option to overwrite existing files,
as well as options to print header lines, and IG Script input as part of the output (for options see tabular.IG_SCRIPT_INCLUSION_OPTIONS).
*/
func generateGoogleSheetsOutput(statementMap []map[string]string, originalStatement string, igScriptInput string, headerCols []string, headerColsNames []string, separator string, filename string, overwrite bool, printHeaders bool, printOriginalStatement string, printIgScriptInput string) (string, tree.ParsingError) {

	// Quote to terminate input string for Google Sheets interpretation
	quote := "\""
	// Line prefix for Google Sheets
	bpre := strings.Builder{}
	bpre.WriteString("=SPLIT(")
	bpre.WriteString(quote)
	prefix := bpre.String()
	// Linebreak at the end of each entry
	linebreak := "\n"
	// Line suffix for Google Sheets (e.g., "; "|")" )
	bsuf := strings.Builder{}
	bsuf.WriteString(quote)
	bsuf.WriteString("; \"")
	bsuf.WriteString(separator)
	bsuf.WriteString("\")")
	bsuf.WriteString(linebreak)
	suffix := bsuf.String()

	// Delegate actual printing
	return printTabularOutput(statementMap, originalStatement, igScriptInput, headerCols, headerColsNames, prefix, stmtIdPrefix, suffix, separator, filename, overwrite, printHeaders, printOriginalStatement, printIgScriptInput)
}

/*
Generates IDs for extrapolated statements by appending a suffix identifying the separate statements to a given base ID.
The returned map can be used for dynamic resolution of IDs and associated nodes (e.g., to generate logical linkages).
Relies on #extrapolatedStatementSuffixSeparator as separator between base ID and suffix.
*/
func generateExtrapolatedStatementIDs(stmts []*tree.Node, baseId string) map[*tree.Node]string {

	stmtIds := map[*tree.Node]string{}

	for i, v := range stmts {
		// Generate unique derived ID per statement by appending and incrementing suffix
		stmtIds[v] = baseId + extrapolatedStatementSuffixSeparator + strconv.Itoa(i+1)
	}

	return stmtIds
}

/*
Generates logical linkages string from a given statement to all other extrapolated statements to be used in tabular output.
Takes potentially prepopulated logical link string, source node, node array of all nodes (filters against source node) as well
as the pregenerated map of Node-to-ID mappings (generated using #generateExtrapolatedStatementIDs).
*/
func generateLogicalLinkageForExtrapolatedStatements(logicalExpressionString string, source *tree.Node, stmts []*tree.Node, IDs map[*tree.Node]string) (string, tree.ParsingError) {

	//logLinkString := ""
	builder := strings.Builder{}
	// Switch to assess whether first operator string has already been appended (to ensure correct separation of additional logical operator expressions)
	logicalStringInitiated := false
	if len(logicalExpressionString) > 0 {
		// since the operator string is not empty, it must already contain decomposed operator strings - thus requiring separation of further entries
		logicalStringInitiated = true
	}

	for _, v := range stmts {
		if v != source {

			// Extract logical linkages from node tree
			res, ops, err := tree.FindLogicalLinkage(source, v)
			if err.ErrorCode != tree.TREE_NO_ERROR {
				errorMsg := fmt.Sprint("Error when parsing retrieving operator linkages: ", err.ErrorMessage)
				Println(errorMsg)
				return "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_LOGICAL_EXPRESSION_GENERATION}
			}
			if res {
				if CollapseOperators() {
					Println("Collapsing adjacent AND, bAND and wAND operators ...")
					// Collapse adjacent AND operators
					ops = tree.CollapseAdjacentOperators(ops, []string{tree.AND, tree.SAND_BETWEEN_COMPONENTS, tree.SAND_WITHIN_COMPONENTS})
				}
				Println("Node", IDs[source], "has linkage", ops, "to", IDs[v])
				if logicalStringInitiated {
					// Append logical operator separator if logical operator linkage already exists
					builder.WriteString(logicalOperatorSeparator)
				} else {
					// Any further printing of logical operators for given component will lead to addition of separator
					logicalStringInitiated = true
				}
				// ... and append to logical expression column string
				builder.WriteString(fmt.Sprint(ops))
				// Internal separator between logical operator and target statement ID
				builder.WriteString(".")
				// Leading bracket
				builder.WriteString(logicalCombinationLeft)
				// Write ID for given target node
				builder.WriteString(IDs[v])
				// Trailing bracket
				builder.WriteString(logicalCombinationRight)
			}
		}
	}
	// Return generated logical expression for given component
	return builder.String(), tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Generates combined tabular output for given statements in node array.
Uses #GenerateTabularOutputFromParsedStatement function internally.
*/
func GenerateTabularOutputFromParsedStatements(stmts []*tree.Node, annotations interface{}, originalStatement string, igScriptInput string, stmtId string, filename string, overwrite bool, aggregateImplicitLinkages bool, separator string, outputFormat string, printHeaders bool, printOriginalStatement string, printIgScriptInput string) []TabularOutputResult {

	// Instance holding parsed output
	results := []TabularOutputResult{}

	// Remove potential line breaks from original and IG Script input
	originalStatement = CleanInput(originalStatement, separator)
	igScriptInput = CleanInput(igScriptInput, separator)

	for i, stmtNode := range stmts {
		Println("Processing output for node entry ", i)

		// Extract potential hierarchy
		topLevelStmts := stmtNode.GetTopLevelStatementNodes()

		// Prepare return structure
		var res TabularOutputResult

		overwriteFile := overwrite
		printHeadersInFile := printHeaders

		for j, topLevelStmt := range topLevelStmts {

			// Suppress repeated headers if appending extrapolated statements to output file
			if j > 0 {
				// temporary override if component pair combinations are printed (chained printing)
				overwriteFile = false
				printHeadersInFile = false
			}

			// single node: simply parse node in isolation
			res = GenerateTabularOutputFromParsedStatement(topLevelStmt, topLevelStmts, originalStatement, igScriptInput, annotations, stmtId, filename, overwriteFile, aggregateImplicitLinkages, separator, outputFormat, printHeadersInFile, printOriginalStatement, printIgScriptInput)
			if res.Error.ErrorCode != tree.PARSING_NO_ERROR {
				Println("Error during output generation for single statement. Statement ignored from output (Statement node: " + stmtNode.String() + ")")

			}
			results = append(results, res)
		}
	}

	return results
}

/*
Generates Google Sheets tabular output for a given parsed statement, with a given statement ID.
Generates all substatements and logical combination linkages in specified output format (e.g., Google Sheets, CSV).
Takes original statement input to potentially include as separate column (parameterized via printOriginalStatement).
Takes original IG Script input to potentially include as separate column (parameterized via printIgScriptInput).
Additionally returns array of statement entries, header symbols and corresponding header symbol names wrapped in generic return structure.
Allows for specification of statement-level annotations passed to output.
Allows for specification of separator to delimit generated flat file output.
Allows for specification of output file type (e.g., Google Sheets, CSV) based on constants #OUTPUT_TYPE_GOOGLE_SHEETS or #OUTPUT_TYPE_CSV.
If filename is provided, the result is printed to the corresponding file.
It is further necessary to indicate whether files should be overwritten or appended to
If printHeaders is true, the header row will be included in output.
printOriginalStatement indicates the inclusion of Original Statement in output (for options see tabular.ORIGINAL_STATEMENT_INCLUSION_OPTIONS).
printIgScriptInput indicates the inclusion of IG Script in output (for options see tabular.IG_SCRIPT_INCLUSION_OPTIONS).
*/
func GenerateTabularOutputFromParsedStatement(node *tree.Node, allStmts []*tree.Node, originalStatement string, igScriptInput string, annotations interface{}, stmtId string, filename string, overwrite bool, aggregateImplicitLinkages bool, separator string, outputFormat string, printHeaders bool, printOriginalStatement string, printIgScriptInput string) TabularOutputResult {

	// Prepopulate derived IDs for uniform access
	derivedIDs := map[*tree.Node]string{}
	derivedIDs[node] = stmtId

	// String holding all logical linkages to other statements
	logicalLinkageStmts := ""
	// Determine linkages to other statements (if extrapolating is needed)
	if len(allStmts) > 1 {

		// Overwrite if multiple statements
		derivedIDs = generateExtrapolatedStatementIDs(allStmts, stmtId)

		// Generate strings indicating linkages between statements
		linkString, err := generateLogicalLinkageForExtrapolatedStatements("", node, allStmts, derivedIDs)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			Println("Error when generating extrapolated statement linkages: ", err)
			return TabularOutputResult{}
		}
		// Append links for all statements
		logicalLinkageStmts += linkString
	}

	var stmt *tree.Statement

	Println(" Step: Extracting leaf arrays")
	// Extract statement from node - varies whether it is ...
	switch reflect.TypeOf(node.Entry) {
	case reflect.TypeOf(&tree.Statement{}):
		// first-order statement component
		stmt = node.Entry.(*tree.Statement)
	case reflect.TypeOf([]*tree.Node{}):
		// component-level nested pair combination
		stmt = node.Entry.([]*tree.Node)[0].Entry.(*tree.Statement)
	default:
		result := TabularOutputResult{}
		result.Error = tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNKNOWN_INPUT_TYPE, ErrorMessage: "Parsing failed for unknown input type " + reflect.TypeOf(node.Entry).String() +
			". Please report this error (alongside Request ID and input statement) to the developer for further investigation."}
		return result
	}
	// Retrieve leaf arrays from generated tree (alongside frequency indications for components)
	leafArrays, componentRefs := stmt.GenerateLeafArrays(aggregateImplicitLinkages)

	Println(" Generated leaf arrays: ", leafArrays, " component: ", componentRefs)

	// Prepare return structure
	result := TabularOutputResult{Output: "", StatementMap: nil, HeaderSymbols: nil, HeaderNames: nil}

	Println(" Step: Generate permutations of leaf arrays (atomic statements)")
	// Generate all permutations of logically-linked components to produce statements
	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		result.Error = err
		return result
	}

	Println(" Generated permutations: ", res)

	Println(" Step: Generate logical operators for atomic statements")
	// Extract logical operator links
	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	Println(" Links:", links)

	Println(" Step: Generate tabular output")

	// Prepare export to tabular output (including pre-generated annotations and logical linkage to other statements)
	result.StatementMap, result.HeaderSymbols, result.HeaderNames, result.Error = generateStatementMatrix(res, annotations, logicalLinkageStmts, componentRefs, links, derivedIDs[node], separator, outputFormat, printHeaders)
	if result.Error.ErrorCode != tree.PARSING_NO_ERROR {
		return result
	}

	// Default output
	result.Output = ""

	switch outputFormat {
	case OUTPUT_TYPE_NONE:
		// No output generated (useful for internal use such as parsing of nested statements) - simply return matrices and empty flat output
		return result
	case OUTPUT_TYPE_GOOGLE_SHEETS:
		// Create Google Sheets output based on generated map, alongside header names as output
		result.Output, result.Error = generateGoogleSheetsOutput(result.StatementMap, originalStatement, igScriptInput, result.HeaderSymbols, result.HeaderNames, separator, filename, overwrite, printHeaders, printOriginalStatement, printIgScriptInput)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return result
		}
	case OUTPUT_TYPE_CSV:
		// Create CSV output based on generated map, alongside header names as output
		result.Output, result.Error = generateCSVOutput(result.StatementMap, originalStatement, igScriptInput, result.HeaderSymbols, result.HeaderNames, separator, filename, overwrite, printHeaders, printOriginalStatement, printIgScriptInput)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return result
		}
	default:
		result = TabularOutputResult{Output: "", StatementMap: nil, HeaderSymbols: nil, HeaderNames: nil, Error: tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_OUTPUT_TYPE, ErrorMessage: "Invalid output type specified. Should be Google Sheets or CSV."}}
	}

	return result
}

/*
Generates IG 2.0 header row and appends it to given string based on component frequency input. It further returns a slice
containing header information.
*/
func generateHeaderRow(stringToAppendTo string, componentFrequency map[string]int, separator string) (string, []string, []string, tree.ParsingError) {

	if separator == "" {
		return "", nil, nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_MISSING_SEPARATOR_VALUE,
			ErrorMessage: "Value for separator symbol is invalid."}
	}

	// Builder to append to input
	builder := strings.Builder{}
	// Add original string
	builder.WriteString(stringToAppendTo)

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
				headerSymbolsName += indexSymbol + strconv.Itoa(i+1)
				// Append suffix for cached header IDs (for logical operators)
				headerSymbol += indexSymbol + strconv.Itoa(i+1)
			}
			// Store key for header used in logical operators
			headerSymbols = append(headerSymbols, headerSymbol)
			headerSymbolsNames = append(headerSymbolsNames, headerSymbolsName)
			// Append full header names to string
			builder.WriteString(headerSymbolsName)
			builder.WriteString(separator)
			i++
		}
	}
	// Convert to string
	stringToAppendTo = builder.String()
	// Cut off last separator
	stringToAppendTo = stringToAppendTo[0 : len(stringToAppendTo)-len(separator)]
	// Return generated string as well as symbol map and mapped names
	return stringToAppendTo, headerSymbols, headerSymbolsNames, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Generates logical expression string for given component entry for a given statement.
It relies on the expression string as input, alongside the statement of concern, as well as component index.
In addition, a slice of all header symbols is required (to generate reference to columns in logical expressions),
as well as the logical links for a given component value. Finally, the statement ID is used to generate the corresponding
substatement IDs used in the link references.
It returns the link for the particular table entry.
*/
func generateLogicalLinksExpressionForGivenComponentValue(logicalExpressionString string, statement []*tree.Node,
	componentIdx int, headerSymbols []string, logicalLinks []map[*tree.Node][]string, stmtId string) (string, tree.ParsingError) {
	// Check for logical operator linkage based on index
	linksForElement := logicalLinks[componentIdx]
	Println("Links for element: ", linksForElement)

	// Node key array (maintaining order of iteration)
	nodesKeys := []*tree.Node{}

	// Switch to assess whether first operator string has already been appended (to ensure correct separation of additional logical operator expressions)
	logicalStringInitiated := false
	if len(logicalExpressionString) > 0 {
		// since the operator string is not empty, it must already contain decomposed operator strings - thus requiring separation of further entries
		logicalStringInitiated = true
	}

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
			leaves := [][]*tree.Node{}
			if tree.AGGREGATE_IMPLICIT_LINKAGES {
				// Retrieve actual root node, not just the one that sits below synthetic linkage
				leaves = firstKey.GetRootNode().GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)
				Println("Root:", firstKey.GetRootNode())
			} else {
				// Retrieve all nodes up to synthetic linkage
				leaves = firstKey.GetNodeBelowSyntheticRootNode().GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)
				Println("Synthetic Root:", firstKey.GetNodeBelowSyntheticRootNode())
			}

			if len(leaves) > 0 {
				nodesKeys = leaves[0]
			} else {
				Println("No component keys to iterate over for logical relationships")
			}
		} else {
			Println("No component keys to iterate over for logical relationships")
		}

		// ALTERNATIVE: Sorting based on alphabet by interface
		//sort.Sort(tree.ByEntry(nodesKeys))

		Println("Sorted keys: ", nodesKeys)
	}

	// Builder for string generation
	builder := strings.Builder{}
	builder.WriteString(logicalExpressionString)

	// Check that entries for own component value exist
	if linksForElement[statement[componentIdx]] != nil {
		// Iterate through all component values based on ordered keys
		for _, nodesKey := range nodesKeys {
			// Extract node
			otherNode := nodesKey
			// Extract references attached to node
			linkedElement := linksForElement[nodesKey]

			// if target node is different ...
			if otherNode != statement[componentIdx] {

				if len(linkedElement) > 0 {
					Println("Testing other node: ", otherNode, " with elements ", linkedElement)
					// find operator
					res, ops, err := tree.FindLogicalLinkage(statement[componentIdx], otherNode)
					if err.ErrorCode != tree.TREE_NO_ERROR {
						errorMsg := fmt.Sprint("Error when parsing retrieving operator linkages: ", err.ErrorMessage)
						log.Println(errorMsg)
						return "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_LOGICAL_EXPRESSION_GENERATION}
					}
					if res {
						if CollapseOperators() {
							Println("Collapsing adjacent AND, bAND and wAND operators ...")
							// Collapse adjacent AND operators
							ops = tree.CollapseAdjacentOperators(ops, []string{tree.AND, tree.SAND_BETWEEN_COMPONENTS, tree.SAND_WITHIN_COMPONENTS})
						}

						Println("Node has linkage ", ops)
						if logicalStringInitiated {
							// Append logical operator separator if logical operator linkage already exists
							builder.WriteString(logicalOperatorSeparator)
						} else {
							// Any further printing of logical operators for given component will lead to addition of separator
							logicalStringInitiated = true
						}
						// ... and append to logical expression column string
						builder.WriteString(fmt.Sprint(ops))
						// Statement component identifier
						if ProduceDynamicOutput() {
							// Based on index or parsed input nodes
							builder.WriteString(".")
							builder.WriteString(headerSymbols[componentIdx])
							builder.WriteString(".")
						} else {
							// Based on name of current element
							builder.WriteString(".")
							builder.WriteString(statement[componentIdx].GetComponentName())
							builder.WriteString(".")
						}
						// Leading bracket
						builder.WriteString(logicalCombinationLeft)
						// Prepare intermediate structure to store statement references
						stmtsRefs := strings.Builder{}

						Println("Target node IDs: ", linkedElement)
						for lks := range linkedElement {
							//Println("Found pointer from ", statement[componentIdx] ," to ", otherNode , " as ", generateStatementID(stmtId, lks + 1))
							// Append actual statement id
							stmtsRefs.WriteString(generateStatementIDString(stmtId, linkedElement[lks]))
							if lks < len(linkedElement)-1 {
								stmtsRefs.WriteString(logicalOperatorStmtRefSeparator)
							}
						}

						// Add statement reference and trailing bracket
						builder.WriteString(stmtsRefs.String())
						builder.WriteString(logicalCombinationRight)
					}
					Println("Added logical relationships for value", otherNode, ", elements:", logicalExpressionString)
				} else {
					Println("Did not find target links for", otherNode, "- did not add logical operator links for component")
				}
			}
		}
	}
	// Return generated logical expression for given component
	return builder.String(), tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Generate statement ID from main statement ID and index of iterated substatement
*/
func generateStatementIDString(mainID string, subStmtIndex string) string {
	b := strings.Builder{}
	b.WriteString(mainID)
	b.WriteString(stmtIdSeparator)
	b.WriteString(subStmtIndex)
	return b.String()
}

/*
Generate statement ID from main statement ID and index of iterated substatement
*/
func generateStatementIDint(mainID string, subStmtIndex int) string {
	b := strings.Builder{}
	b.WriteString(mainID)
	b.WriteString(stmtIdSeparator)
	b.WriteString(strconv.Itoa(subStmtIndex))
	return b.String()
}

/*
Writes data to given file - overwrites or appends to file as specified.
*/
func WriteToFile(filename string, content string, overwrite bool) error {

	// Open file, create if not existing, and append if existing
	var f *os.File
	var err error
	if overwrite {
		// Open file in overwrite mode
		f, err = os.Create(filename)
	} else {
		// Open file in append mode
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	if err != nil {
		return err
	}

	// Defer closing of file
	defer func() error {
		err := f.Close()
		if err != nil {
			log.Println("Error when writing file", filename, "Error:", err.Error())
			return err
		}
		return nil
	}()

	// Prepare data
	data := []byte(content)

	// Write data
	_, err2 := f.Write(data)
	if err2 != nil {
		return err2
	}
	log.Println("Wrote to file " + filename)

	// No error
	return nil
}
