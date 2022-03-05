package converter

import (
	"IG-Parser/exporter"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConverterHandlerSheetsGet(t *testing.T) {

	// Initialize templates
	Init()
	// Spin up server
	server := httptest.NewServer(http.HandlerFunc(ConverterHandlerSheets))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}
	res, err := client.Get(server.URL)
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err)
	}

	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Error when reading response. Error:", err.Error())
	}

	outputString := string(output)

	// Read reference file
	content, err2 := ioutil.ReadFile("TestConverterHandlerSheetsGet.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	expectedOutput := string(content)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err2 := exporter.WriteToFile("errorOutput.error", outputString)
		if err2 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}
