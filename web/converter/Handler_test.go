package converter

import (
	"IG-Parser/core/exporter"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Default file name for error output
const errorFile = "errorOutput.error"

func TestConverterHandlerSheetsGet(t *testing.T) {

	// Initialize templates
	Init()
	// Deactivate logging
	Logging = false
	// Spin up server
	server := httptest.NewServer(http.HandlerFunc(ConverterHandlerTabular))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}
	res, err := client.Get(server.URL)
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err.Error())
	}

	if res.Status != "200 OK" {
		t.Fatal("Request returning non-200 status code: " + res.Status)
	}

	output, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		t.Fatal("Error when reading response. Error:", err2.Error())
	}

	outputString := string(output)

	// Convert output to \n linebreaks only (may contain a mix) - for OS-independent comparison
	outputString = strings.ReplaceAll(outputString, "\\r\\n", "\\n")

	// Read reference file
	content, err5 := os.ReadFile("TestConverterHandlerGoogleSheetsGet.test")
	if err5 != nil {
		t.Fatal("Error attempting to read test text input. Error:", err5.Error())
	}

	expectedOutput := string(content)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

}

/*
Tests standard GET request without any prior input.
*/
func TestConverterHandlerVisualGet(t *testing.T) {

	// Initialize templates
	Init()
	// Deactivate logging
	Logging = false
	// Spin up server
	server := httptest.NewServer(http.HandlerFunc(ConverterHandlerVisual))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}
	res, err := client.Get(server.URL)
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err.Error())
	}

	if res.Status != "200 OK" {
		t.Fatal("Request returning non-200 status code: " + res.Status)
	}

	output, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		t.Fatal("Error when reading response. Error:", err2.Error())
	}

	outputString := string(output)

	// Convert output to \n linebreaks only (may contain a mix) - for OS-independent comparison
	outputString = strings.ReplaceAll(outputString, "\\r\\n", "\\n")

	// Read reference file
	content, err5 := os.ReadFile("TestConverterHandlerVisualGet.test")
	if err5 != nil {
		t.Fatal("Error attempting to read test text input. Error:", err5.Error())
	}

	expectedOutput := string(content)

	// Compare to actual output
	if outputString != expectedOutput {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

}

/*
Tests POST request with given input for visual output.
*/
func TestConverterHandlerVisualPost(t *testing.T) {

	// Initialize templates
	Init()
	// Deactivate logging
	Logging = false
	// Spin up server
	server := httptest.NewServer(http.HandlerFunc(ConverterHandlerVisual))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}

	body := "rawStmt=Regional+Managers%2C+on+behalf+of+the+Secretary%2C+may+review%2C+reward%2C+or+sanction+approved+certified+production+and+handling+operations+and+accredited+certifying+agents+for+compliance+with+the+Act+or+regulations+in+this+part%2C+under+the+condition+that+Operations+were+non-compliant+or+violated+organic+farming+provisions+and+Manager+has+concluded+investigation.&codedStmt=A%2Cp%28Regional%29+A%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Managers%29%2C+Cex%28on+behalf+of+the+Secretary%29%2C+D%5Bstringency%3Dpermissive%5D%28may%29+I%5Bact%3Dperformance%5D%28%28review+%5BAND%5D+%28reward+%5BXOR%5D+sanction%29%29%29+Bdir%2Cp%28approved%29+Bdir1%2Cp%28certified%29+Bdir1%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28production+%5Boperations%5D%29+and+Bdir%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28handling+operations%29+and+Bdir2%2Cp%28accredited%29+Bdir2%5Brole%3Dmonitor%2Ctype%3Danimate%5D%28certifying+agents%29+Cex%5Bctx%3Dpurpose%5D%28for+compliance+with+the+%28Act+or+%5BXOR%5D+regulations+in+this+part%29%29+under+the+condition+that+%7BCac%5Bstate%5D%7BA%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28Operations%29+I%5Bact%3Dviolate%5D%28were+%28non-compliant+%5BOR%5D+violated%29%29+Bdir%5Btype%3Dinanimate%5D%28organic+farming+provisions%29%7D+%5BAND%5D+Cac%5Bstate%5D%7BA%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Manager%29+I%5Bact%3Dterminate%5D%28has+concluded%29+Bdir%5Btype%3Dactivity%5D%28investigation%29%7D%7D.&propertyTree=on&canvasHeight=2000&canvasWidth=4000"

	res, err := client.Post(server.URL, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err.Error())
	}

	if res.Status != "200 OK" {
		t.Fatal("Request returning non-200 status code: " + res.Status)
	}

	// Extract response body
	output, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		t.Fatal("Error when reading response. Error:", err2.Error())
	}

	outputString := string(output)

	// Convert output to \n linebreaks only (may contain a mix) - for OS-independent comparison
	outputString = strings.ReplaceAll(outputString, "\\r\\n", "\\n")

	// Delimiter terminating header information
	// TODO: Review whenever template changes
	headDelimiter := "submit"

	// Extract index of transaction ID
	endIdx := strings.Index(outputString, headDelimiter)
	// Extract response header part
	responseHead := outputString[:endIdx]

	// Read reference file
	content, err5 := os.ReadFile("TestConverterHandlerVisualPost.test")
	if err5 != nil {
		t.Fatal("Error attempting to read test text input. Error:", err5.Error())
	}

	expectedOutput := string(content)

	// Extract expected response header part
	expEndIdx := strings.Index(expectedOutput, headDelimiter)
	// Extract expected response header part
	expectedHead := expectedOutput[:expEndIdx]

	// Compare to actual output
	if responseHead != expectedHead {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

	// Delimiter to isolate information trailing variable Request ID
	// TODO: Review whenever template changes
	tailDelimiter := "<script src=\"/libraries"

	// Identify index of string
	tailIdx := strings.Index(outputString, tailDelimiter)
	// Extract Response tail
	responseTail := outputString[tailIdx:]

	// Extract expected response tail part
	expTailIdx := strings.Index(expectedOutput, tailDelimiter)
	// Extract expected response tail part
	expectedTail := expectedOutput[expTailIdx:]

	// Compare to actual output
	if responseTail != expectedTail {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

}

/*
Tests POST request with given input for Google Sheets output without header line
*/
func TestConverterHandlerGoogleSheetsPost(t *testing.T) {

	// Initialize templates
	Init()
	// Deactivate logging
	Logging = false
	// Spin up server
	server := httptest.NewServer(http.HandlerFunc(ConverterHandlerTabular))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}

	body := "rawStmt=Regional+Managers%2C+on+behalf+of+the+Secretary%2C+may+review%2C+reward%2C+or+sanction+approved+certified+production+and+handling+operations+and+accredited+certifying+agents+for+compliance+with+the+Act+or+regulations+in+this+part%2C+under+the+condition+that+Operations+were+non-compliant+or+violated+organic+farming+provisions+and+Manager+has+concluded+investigation.&codedStmt=A%2Cp%28Regional%29+A%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Managers%29%2C+Cex%28on+behalf+of+the+Secretary%29%2C+D%5Bstringency%3Dpermissive%5D%28may%29+I%5Bact%3Dperformance%5D%28%28review+%5BAND%5D+%28reward+%5BXOR%5D+sanction%29%29%29+Bdir%2Cp%28approved%29+Bdir1%2Cp%28certified%29+Bdir1%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28production+%5Boperations%5D%29+and+Bdir%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28handling+operations%29+and+Bdir2%2Cp%28accredited%29+Bdir2%5Brole%3Dmonitor%2Ctype%3Danimate%5D%28certifying+agents%29+Cex%5Bctx%3Dpurpose%5D%28for+compliance+with+the+%28Act+or+%5BXOR%5D+regulations+in+this+part%29%29+under+the+condition+that+%7BCac%5Bstate%5D%7BA%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28Operations%29+I%5Bact%3Dviolate%5D%28were+%28non-compliant+%5BOR%5D+violated%29%29+Bdir%5Btype%3Dinanimate%5D%28organic+farming+provisions%29%7D+%5BAND%5D+Cac%5Bstate%5D%7BA%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Manager%29+I%5Bact%3Dterminate%5D%28has+concluded%29+Bdir%5Btype%3Dactivity%5D%28investigation%29%7D%7D.&stmtId=123&igExtended=on&outputType=Google+Sheets"

	res, err := client.Post(server.URL, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err.Error())
	}

	if res.Status != "200 OK" {
		t.Fatal("Request returning non-200 status code: " + res.Status)
	}

	// Extract response body
	output, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		t.Fatal("Error when reading response. Error:", err2.Error())
	}

	outputString := string(output)

	// Convert output to \n linebreaks only (may contain a mix) - for OS-independent comparison
	outputString = strings.ReplaceAll(outputString, "\\r\\n", "\\n")

	// Delimiter terminating header information
	// TODO: Review whenever template changes
	headDelimiter := "submit"

	// Extract index of transaction ID
	endIdx := strings.Index(outputString, headDelimiter)
	// Extract response header part
	responseHead := outputString[:endIdx]

	// Read reference file
	content, err5 := ioutil.ReadFile("TestConverterHandlerGoogleSheetsPost.test")
	if err5 != nil {
		t.Fatal("Error attempting to read test text input. Error:", err5.Error())
	}

	expectedOutput := string(content)

	// Extract expected response header part
	expEndIdx := strings.Index(expectedOutput, headDelimiter)
	// Extract expected response header part
	expectedHead := expectedOutput[:expEndIdx]

	// Compare to actual output
	if responseHead != expectedHead {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

	// Delimiter to isolate information trailing variable Request ID
	// TODO: Review whenever template changes
	tailDelimiter := "<div class=\"output\">"

	// Identify index of string
	tailIdx := strings.Index(outputString, tailDelimiter)
	// Extract Response tail
	responseTail := outputString[tailIdx:]

	// Extract expected response tail part
	expTailIdx := strings.Index(expectedOutput, tailDelimiter)
	// Extract expected response tail part
	expectedTail := expectedOutput[expTailIdx:]

	// Compare to actual output
	if responseTail != expectedTail {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

}

/*
Tests POST request with given input for CSV output, without header line
*/
func TestConverterHandlerCSVPost(t *testing.T) {

	// Initialize templates
	Init()
	// Deactivate logging
	Logging = false
	// Spin up server
	server := httptest.NewServer(http.HandlerFunc(ConverterHandlerTabular))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}

	body := "rawStmt=Regional+Managers%2C+on+behalf+of+the+Secretary%2C+may+review%2C+reward%2C+or+sanction+approved+certified+production+and+handling+operations+and+accredited+certifying+agents+for+compliance+with+the+Act+or+regulations+in+this+part%2C+under+the+condition+that+Operations+were+non-compliant+or+violated+organic+farming+provisions+and+Manager+has+concluded+investigation.&codedStmt=A%2Cp%28Regional%29+A%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Managers%29%2C+Cex%28on+behalf+of+the+Secretary%29%2C+D%5Bstringency%3Dpermissive%5D%28may%29+I%5Bact%3Dperformance%5D%28%28review+%5BAND%5D+%28reward+%5BXOR%5D+sanction%29%29%29+Bdir%2Cp%28approved%29+Bdir1%2Cp%28certified%29+Bdir1%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28production+%5Boperations%5D%29+and+Bdir%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28handling+operations%29+and+Bdir2%2Cp%28accredited%29+Bdir2%5Brole%3Dmonitor%2Ctype%3Danimate%5D%28certifying+agents%29+Cex%5Bctx%3Dpurpose%5D%28for+compliance+with+the+%28Act+or+%5BXOR%5D+regulations+in+this+part%29%29+under+the+condition+that+%7BCac%5Bstate%5D%7BA%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28Operations%29+I%5Bact%3Dviolate%5D%28were+%28non-compliant+%5BOR%5D+violated%29%29+Bdir%5Btype%3Dinanimate%5D%28organic+farming+provisions%29%7D+%5BAND%5D+Cac%5Bstate%5D%7BA%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Manager%29+I%5Bact%3Dterminate%5D%28has+concluded%29+Bdir%5Btype%3Dactivity%5D%28investigation%29%7D%7D.&stmtId=123&igExtended=on&outputType=CSV+format"

	res, err := client.Post(server.URL, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err.Error())
	}

	if res.Status != "200 OK" {
		t.Fatal("Request returning non-200 status code: " + res.Status)
	}

	// Extract response body
	output, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		t.Fatal("Error when reading response. Error:", err2.Error())
	}

	outputString := string(output)

	// Convert output to \n linebreaks only (may contain a mix) - for OS-independent comparison
	outputString = strings.ReplaceAll(outputString, "\\r\\n", "\\n")

	// Delimiter terminating header information
	// TODO: Review whenever template changes
	headDelimiter := "submit"

	// Extract index of transaction ID
	endIdx := strings.Index(outputString, headDelimiter)
	// Extract response header part
	responseHead := outputString[:endIdx]

	// Read reference file
	content, err5 := os.ReadFile("TestConverterHandlerCsvPost.test")
	if err5 != nil {
		t.Fatal("Error attempting to read test text input. Error:", err5.Error())
	}

	expectedOutput := string(content)

	// Extract expected response header part
	expEndIdx := strings.Index(expectedOutput, headDelimiter)
	// Extract expected response header part
	expectedHead := expectedOutput[:expEndIdx]

	// Compare to actual output
	if responseHead != expectedHead {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

	// Delimiter to isolate information trailing variable Request ID
	// TODO: Review whenever template changes
	tailDelimiter := "<div class=\"output\">"

	// Identify index of string
	tailIdx := strings.Index(outputString, tailDelimiter)
	// Extract Response tail
	responseTail := outputString[tailIdx:]

	// Extract expected response tail part
	expTailIdx := strings.Index(expectedOutput, tailDelimiter)
	// Extract expected response tail part
	expectedTail := expectedOutput[expTailIdx:]

	// Compare to actual output
	if responseTail != expectedTail {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

}

/*
Tests POST request with given input for CSV output with explicit header row request
*/
func TestConverterHandlerCSVPostWithExplicitHeader(t *testing.T) {

	// Initialize templates
	Init()
	// Deactivate logging
	Logging = false
	// Spin up server
	server := httptest.NewServer(http.HandlerFunc(ConverterHandlerTabular))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}

	body := "rawStmt=Regional+Managers%2C+on+behalf+of+the+Secretary%2C+may+review%2C+reward%2C+or+sanction+approved+certified+production+and+handling+operations+and+accredited+certifying+agents+for+compliance+with+the+Act+or+regulations+in+this+part%2C+under+the+condition+that+Operations+were+non-compliant+or+violated+organic+farming+provisions+and+Manager+has+concluded+investigation.&codedStmt=A%2Cp%28Regional%29+A%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Managers%29%2C+Cex%28on+behalf+of+the+Secretary%29%2C+D%5Bstringency%3Dpermissive%5D%28may%29+I%5Bact%3Dperformance%5D%28%28review+%5BAND%5D+%28reward+%5BXOR%5D+sanction%29%29%29+Bdir%2Cp%28approved%29+Bdir1%2Cp%28certified%29+Bdir1%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28production+%5Boperations%5D%29+and+Bdir%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28handling+operations%29+and+Bdir2%2Cp%28accredited%29+Bdir2%5Brole%3Dmonitor%2Ctype%3Danimate%5D%28certifying+agents%29+Cex%5Bctx%3Dpurpose%5D%28for+compliance+with+the+%28Act+or+%5BXOR%5D+regulations+in+this+part%29%29+under+the+condition+that+%7BCac%5Bstate%5D%7BA%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28Operations%29+I%5Bact%3Dviolate%5D%28were+%28non-compliant+%5BOR%5D+violated%29%29+Bdir%5Btype%3Dinanimate%5D%28organic+farming+provisions%29%7D+%5BAND%5D+Cac%5Bstate%5D%7BA%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Manager%29+I%5Bact%3Dterminate%5D%28has+concluded%29+Bdir%5Btype%3Dactivity%5D%28investigation%29%7D%7D.&stmtId=123&igExtended=on&includeHeaders=on&outputType=CSV+format"

	res, err := client.Post(server.URL, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err.Error())
	}

	if res.Status != "200 OK" {
		t.Fatal("Request returning non-200 status code: " + res.Status)
	}

	// Extract response body
	output, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		t.Fatal("Error when reading response. Error:", err2.Error())
	}

	outputString := string(output)

	// Convert output to \n linebreaks only (may contain a mix) - for OS-independent comparison
	outputString = strings.ReplaceAll(outputString, "\\r\\n", "\\n")

	// Delimiter terminating header information
	// TODO: Review whenever template changes
	headDelimiter := "submit"

	// Extract index of transaction ID
	endIdx := strings.Index(outputString, headDelimiter)
	// Extract response header part
	responseHead := outputString[:endIdx]

	// Read reference file
	content, err5 := os.ReadFile("TestConverterHandlerCsvPostExplicitHeaders.test")
	if err5 != nil {
		t.Fatal("Error attempting to read test text input. Error:", err5.Error())
	}

	expectedOutput := string(content)

	// Extract expected response header part
	expEndIdx := strings.Index(expectedOutput, headDelimiter)
	// Extract expected response header part
	expectedHead := expectedOutput[:expEndIdx]

	// Compare to actual output
	if responseHead != expectedHead {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

	// Delimiter to isolate information trailing variable Request ID
	// TODO: Review whenever template changes
	tailDelimiter := "<div class=\"output\">"

	// Identify index of string
	tailIdx := strings.Index(outputString, tailDelimiter)
	// Extract Response tail
	responseTail := outputString[tailIdx:]

	// Extract expected response tail part
	expTailIdx := strings.Index(expectedOutput, tailDelimiter)
	// Extract expected response tail part
	expectedTail := expectedOutput[expTailIdx:]

	// Compare to actual output
	if responseTail != expectedTail {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

}

/*
Tests POST request with given input for CSV output, but without header row
*/
func TestConverterHandlerCSVPostWithoutHeaders(t *testing.T) {

	// Initialize templates
	Init()
	// Deactivate logging
	Logging = false
	// Spin up server
	server := httptest.NewServer(http.HandlerFunc(ConverterHandlerTabular))
	// Tear down at the end of the function
	defer server.Close()

	// Read server information
	client := http.Client{}

	body := "rawStmt=Regional+Managers%2C+on+behalf+of+the+Secretary%2C+may+review%2C+reward%2C+or+sanction+approved+certified+production+and+handling+operations+and+accredited+certifying+agents+for+compliance+with+the+Act+or+regulations+in+this+part%2C+under+the+condition+that+Operations+were+non-compliant+or+violated+organic+farming+provisions+and+Manager+has+concluded+investigation.&codedStmt=A%2Cp%28Regional%29+A%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Managers%29%2C+Cex%28on+behalf+of+the+Secretary%29%2C+D%5Bstringency%3Dpermissive%5D%28may%29+I%5Bact%3Dperformance%5D%28%28review+%5BAND%5D+%28reward+%5BXOR%5D+sanction%29%29%29+Bdir%2Cp%28approved%29+Bdir1%2Cp%28certified%29+Bdir1%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28production+%5Boperations%5D%29+and+Bdir%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28handling+operations%29+and+Bdir2%2Cp%28accredited%29+Bdir2%5Brole%3Dmonitor%2Ctype%3Danimate%5D%28certifying+agents%29+Cex%5Bctx%3Dpurpose%5D%28for+compliance+with+the+%28Act+or+%5BXOR%5D+regulations+in+this+part%29%29+under+the+condition+that+%7BCac%5Bstate%5D%7BA%5Brole%3Dmonitored%2Ctype%3Danimate%5D%28Operations%29+I%5Bact%3Dviolate%5D%28were+%28non-compliant+%5BOR%5D+violated%29%29+Bdir%5Btype%3Dinanimate%5D%28organic+farming+provisions%29%7D+%5BAND%5D+Cac%5Bstate%5D%7BA%5Brole%3Denforcer%2Ctype%3Danimate%5D%28Manager%29+I%5Bact%3Dterminate%5D%28has+concluded%29+Bdir%5Btype%3Dactivity%5D%28investigation%29%7D%7D.&stmtId=123&igExtended=on&includeHeaders=off&includeHeaders=off&outputType=CSV+format"

	res, err := client.Post(server.URL, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		t.Fatal("Error when performing HTTP request. Error:", err.Error())
	}

	if res.Status != "200 OK" {
		t.Fatal("Request returning non-200 status code: " + res.Status)
	}

	// Extract response body
	output, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		t.Fatal("Error when reading response. Error:", err2.Error())
	}

	outputString := string(output)

	// Convert output to \n linebreaks only (may contain a mix) - for OS-independent comparison
	outputString = strings.ReplaceAll(outputString, "\\r\\n", "\\n")

	// Delimiter terminating header information
	// TODO: Review whenever template changes
	headDelimiter := "submit"

	// Extract index of transaction ID
	endIdx := strings.Index(outputString, headDelimiter)
	// Extract response header part
	responseHead := outputString[:endIdx]

	// Read reference file
	content, err5 := os.ReadFile("TestConverterHandlerCsvPostWithoutHeaders.test")
	if err5 != nil {
		t.Fatal("Error attempting to read test text input. Error:", err5.Error())
	}

	expectedOutput := string(content)

	// Extract expected response header part
	expEndIdx := strings.Index(expectedOutput, headDelimiter)
	// Extract expected response header part
	expectedHead := expectedOutput[:expEndIdx]

	// Compare to actual output
	if responseHead != expectedHead {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

	// Delimiter to isolate information trailing variable Request ID
	// TODO: Review whenever template changes
	tailDelimiter := "<div class=\"output\">"

	// Identify index of string
	tailIdx := strings.Index(outputString, tailDelimiter)
	// Extract Response tail
	responseTail := outputString[tailIdx:]

	// Extract expected response tail part
	expTailIdx := strings.Index(expectedOutput, tailDelimiter)
	// Extract expected response tail part
	expectedTail := expectedOutput[expTailIdx:]

	// Compare to actual output
	if responseTail != expectedTail {
		fmt.Println("Produced output:\n", outputString)
		fmt.Println("Expected output:\n", expectedOutput)
		err6 := exporter.WriteToFile(errorFile, outputString)
		if err6 != nil {
			t.Fatal("Error attempting to write error file. Error:", err6.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to '" + errorFile + "'")
	}

}
