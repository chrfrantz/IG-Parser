package libraries

import (
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/web/converter"
	"embed"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//go:embed d3.v7.min.js ace/ace.js
var libraryFiles embed.FS

// Default file name for error output
const errorFile = "errorOutput.error"

func TestD3LibraryRead(t *testing.T) {

	// Initialize templates
	converter.Init()
	// Spin up server
	server := httptest.NewServer(http.Handler(http.FileServer(http.FS(libraryFiles))))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}
	res, err := client.Get(server.URL + "/d3.v7.min.js")
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err.Error())
	}
	defer client.CloseIdleConnections()

	if res.Status != "200 OK" {
		t.Fatal("Request returning non-200 status code: " + res.Status)
	}

	// Should not be used in production code, but kept here for comparison purposes.
	output, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		t.Fatal("Error when reading response. Error:", err2.Error())
	}
	// Retrieved content
	outputString := string(output)

	// Read local reference file
	content, err5 := os.ReadFile("d3.v7.min.js")
	if err5 != nil {
		t.Fatal("Error attempting to read test text input. Error:", err5.Error())
	}

	expectedOutput := string(content)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := tabular.WriteToFile(errorFile, outputString, true)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

}

func TestACELibraryRead(t *testing.T) {

	// Initialize templates
	converter.Init()
	// Spin up server
	server := httptest.NewServer(http.Handler(http.FileServer(http.FS(libraryFiles))))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}
	res, err := client.Get(server.URL + "/ace/ace.js")
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err.Error())
	}
	defer client.CloseIdleConnections()

	if res.Status != "200 OK" {
		t.Fatal("Request returning non-200 status code: " + res.Status)
	}

	// Should not be used in production code, but kept here for comparison purposes.
	output, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		t.Fatal("Error when reading response. Error:", err2.Error())
	}
	// Retrieved content
	outputString := string(output)

	// Read local reference file
	content, err5 := os.ReadFile("ace/ace.js")
	if err5 != nil {
		t.Fatal("Error attempting to read test text input. Error:", err5.Error())
	}

	expectedOutput := string(content)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := tabular.WriteToFile(errorFile, outputString, true)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

}
