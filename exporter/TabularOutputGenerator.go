package exporter

import (
	"IG-Parser/tree"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
Generates statement output in Google Sheets format
 */
func GenerateGoogleSheetsOutput(stmts [][]*tree.Node, refs map[string]int, logicalLinks []map[*tree.Node][]int, stmtId string) string {
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
	if refs != nil && len(refs) != 0 {

		output += prefix
		output += "Statement ID" + separator
		// Iterate through component reference map
		for _, v := range tree.IGComponents {
			i := 0
			// Print headers as often as occurring in input file (stmtCt.e., one header for each column)
			for i < refs[v] {
				output += tree.IGComponentNames[v]
				// Store symbols for columns including indices in order of occurrence for use in logical operators
				headerSymbol := v
				// Introduce indices if multiple of the same component
				if refs[v] > 1 {
					// Append suffix for header string
					output += indexSymbol + strconv.Itoa(i + 1)
					// Append suffix for cached header IDs (for logical operators)
					headerSymbol += indexSymbol + strconv.Itoa(i + 1)
				}
				// Store key for header used in logical operators
				headerSymbols = append(headerSymbols, headerSymbol)
				output += separator
				i++
			}
		}
		// Cut off last separator
		output = output[0:len(output)-len(separator)]
		// Complete line
		output += suffix
		//fmt.Println("Header: " + output)
	}

	// Generate entries
	for stmtCt, statement := range stmts {
		//fmt.Println("Statement ", stmtCt, ": ", statement)
		// Start new row
		output += prefix
		// Add statement ID for specific instance
		subStmtId := generateStatementID(stmtId, stmtCt+ 1)
		output += stmtIdPrefix + subStmtId + separator
		ct := 0
		// String linking all logical operators for a given row
		logicalValue := ""
		// Iterate over component index (i.e., column)
		for componentIdx := range statement {
			// Print element
			output += statement[componentIdx].Entry

			fmt.Println("Source node: ", statement[componentIdx])

			// Check for logical operator linkage based on index
			linksForElement := logicalLinks[componentIdx]

			// Check that entries for own component value exist
			if linksForElement[statement[componentIdx]] != nil {
				// Iterate through all component values
				for otherNode, linkedElement := range linksForElement {
					// if target node is different ...
					if otherNode != statement[componentIdx] {
						fmt.Println("Testing other node: ", otherNode, " with elements ", linkedElement)
						// find operator
						res, ops, err := tree.FindLogicalLinkage(statement[componentIdx], otherNode)
						if err.ErrorCode != tree.TREE_NO_ERROR {
							errorMsg := fmt.Sprint("Error when parsing retrieving operator linkages: ", err.ErrorMessage)
							log.Println(errorMsg)
							fmt.Errorf("&v", errorMsg)
							return ""
						}
						if res {
							fmt.Println("Node has linkage ", ops)
							// ... and append to logicalValue column string
							logicalValue += fmt.Sprint(ops)
							// Statement component identifier
							logicalValue += "." + headerSymbols[componentIdx] + "."
							// Leading bracket
							logicalValue +=	"["
							// Prepare intermediate structure to store statement references
							stmtsRefs := ""

							fmt.Println("Target node IDs: ", linkedElement)
							for lks := range linkedElement {
								//fmt.Println("Found pointer from ", statement[componentIdx] ," to ", otherNode , " as ", generateStatementID(stmtId, lks + 1))
								// Append actual statement id
								stmtsRefs += generateStatementID(stmtId, linkedElement[lks] + 1)
								if lks < len(linkedElement)-1 {
									stmtsRefs += logicalOperatorStmtRefSeparator
								}
							}

							//fmt.Println("Statement references prior to conversion: " + stmtsRefs)
							// Attempt conversion to ranges
							//stmtsRefs = convertEnumerationsToRanges(stmtsRefs)
							//fmt.Println("Statement references after conversion: " + stmtsRefs)

							// Add trailing bracket and column ref (to be reviewed)
							logicalValue += stmtsRefs + "]" + logicalOperatorSeparator
						}
						fmt.Println("Added logical relationships for value ", otherNode, ", elements: ", logicalValue)
					}
				}
			}

			// Only append separator if no more elements
			ct++
			if ct < len(statement) {
				output += separator
			}
			//fmt.Println("-->", statement[componentIdx])
		}

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
	return output
}

/*
Generate statement ID from main statement ID and index of iterated substatement
TODO: Fix that one - it does not yet work
 */
func generateStatementID(mainID string, subStmtIndex int) string {
	return mainID + stmtIdSeparator + strconv.Itoa(subStmtIndex)
}

func convertEnumerationsToRanges(enumeration string) string {

	separated := strings.Split(enumeration, logicalOperatorStmtRefSeparator)

	output := ""
	cachedRange := ""
	lastRound := ""
	lastRoundCt := 0
	for i := range separated {
		if i == 0 {
			// Prefill central field and go to next round
			cachedRange += separated[i]
			lastRound = separated[i]
			continue
		}
		/*if cachedRange == "" {
			if output != "" {
				output += logicalOperatorStmtRefSeparator
			}
			// Store first entry
			cachedRange += separated[i]
			// cache last round
			lastRound = separated[i]
			// Save the last time something was saved
			lastRoundCt = i
			continue
		}*/

		lastRoundSuffix := strings.Split(lastRound, stmtIdSeparator)[1]
		lastFloat, err := strconv.ParseFloat(lastRoundSuffix, 8)
		if err != nil {
			log.Fatal("Error converting ", lastRoundSuffix, " to number.")
		}
		fmt.Println(lastFloat)
		thisRoundSuffix := strings.Split(separated[i], stmtIdSeparator)[1]
		thisFloat, err := strconv.ParseFloat(thisRoundSuffix, 8)
		if err != nil {
			log.Fatal("Error converting ", thisRoundSuffix, " to number.")
		}
		fmt.Println(thisFloat)

		if thisFloat-lastFloat > 1 {

			// Store finalized range
			output += cachedRange
			if i-lastRoundCt < 2 {
				output += "-"
			}
			output += lastRound + logicalOperatorStmtRefSeparator
			// Overwrite cache with latest "first element"
			cachedRange = "" //separated[i]
			// Store last element
			lastRound = separated[i]
			lastRoundCt = i
		} else {
			fmt.Println("Difference too small to consider range")
		}
	}
	return output
}

/*
Writes data to given file - overwrites file if existing
 */
func WriteToFile(filename string, content string) {

	// Create file
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Prepare data
	data := []byte(content)

	// Write data
	_, err2 := f.Write(data)
	if err2 != nil {
		log.Fatal(err2)
	}
	log.Println("Wrote file " + filename)
}
