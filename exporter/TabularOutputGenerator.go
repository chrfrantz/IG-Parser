package exporter

import (
	"IG-Parser/tree"
	"fmt"
	"log"
	"os"
	"strconv"
)

const stmtIdSeparator = "."

/*
Generates statement output in Google Sheets format
 */
func GenerateGoogleSheetsOutput(stmts [][]*tree.Node, refs map[string]int, logicalLinks []LogicalOperatorLinkage, stmtId string) string {
	output := ""

	quote := "\""
	prefix := "=SPLIT(" + quote
	linebreak := "\n"
	separator := ";"
	suffix := quote + ", \"" + separator + "\")" + linebreak
	indexSymbol := "_"
	stmtIdPrefix := "'"
	logicalOperatorStmtRefSeparator := ","

	// Generate headers
	if refs != nil && len(refs) != 0 {

		output += prefix
		output += "Statement ID" + separator
		// Iterate through component reference map
		for _, v := range tree.IGComponents {
			i := 0
			// Print headers as often as occurring in input file (i.e., one header for each column)
			for i < refs[v] {
				output += tree.IGComponentNames[v]
				// Introduce indices if multiple of the same component
				if refs[v] > 1 {
					output += indexSymbol + strconv.Itoa(i + 1)
				}
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
	for i, s := range stmts {
		//fmt.Println("Statement ", i, ": ", s)
		// Start new row
		output += prefix
		// Add statement ID for specific instance
		subStmtId := generateStatementID(stmtId, i + 1)
		output += stmtIdPrefix + subStmtId + separator
		ct := 0
		// String linking all logical operators for a given row
		logicalValue := ""
		// Iterate over component index (i.e., column)
		for v := range s {
			output += s[v].Entry

			// Check for logical operator linkage
			linksForElement := logicalLinks[v]
			// if the current checked entry corresponds to entry in linkage structure ...
			if s[v] == linksForElement.Component {
				// ... then retrieve all links
				for l, _ := range linksForElement.LinkedStatements {
					// Only operate on elements that differ from search element
					if l != s[v] {
						// Add leading operator
						logicalValue += fmt.Sprintf("%v", linksForElement.LinkedComponentOperator[l]) + "["
						// Add reference to statements
						for lks := range linksForElement.LinkedStatements[l] {
							//linksForElement.LinkedStatements[l]
							// Append actual statement id
							logicalValue += generateStatementID(stmtId, lks + 1)
							logicalValue += logicalOperatorStmtRefSeparator
						}
						// Add trailing bracket
						logicalValue += "](" + strconv.Itoa(v) + ")|"
					}
				}
			}
			ct++
			if ct < len(s) {
				output += separator
			}
			//fmt.Println("-->", s[v])
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
 */
func generateStatementID(mainID string, subStmtIndex int) string {
	return mainID + stmtIdSeparator + strconv.Itoa(subStmtIndex)
}

func WriteToFile(filename string, content string) {
	f, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	data := []byte(content)

	_, err2 := f.Write(data)

	if err2 != nil {
		log.Fatal(err2)
	}

	log.Println("Wrote file " + filename)
}
