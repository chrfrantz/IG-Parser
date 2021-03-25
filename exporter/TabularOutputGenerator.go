package exporter

import (
	"IG-Parser/tree"
	"log"
	"os"
	"strconv"
)

/*
Generates statement output in Google Sheets format
 */
func GenerateGoogleSheetsOutput(stmts [][]*tree.Node, refs map[string]int, stmtId string) string {
	output := ""

	quote := "\""
	prefix := "=SPLIT(" + quote
	linebreak := "\n"
	separator := ";"
	suffix := quote + ", \"" + separator + "\")" + linebreak
	indexSymbol := "_"
	stmtIdPrefix := "'"
	stmtIdSeparator := "."

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
		// Add statement ID
		output += stmtIdPrefix + stmtId + stmtIdSeparator + strconv.Itoa(i + 1) + separator
		ct := 0
		for v := range s {
			output += s[v].Entry
			ct++
			if ct < len(s) {
				output += separator
			}
			//fmt.Println("-->", s[v])
		}
		output += suffix
	}

	return output
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
