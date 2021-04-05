package helper

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

/*
Generates random unique 8-character ID and filename including current datetime in YYYYMMDDHHMMSS format and ID.
 */
func GenerateUniqueIdAndFilename() (string, string) {
	// Get unique ID
	id := GenerateRandomID(8)
	// Construct filename, and return ID and filename
	return id, GetDateTimeString() + "-" + id + ".log"
}

/*
Generates random string ID of given length
 */
func GenerateRandomID(length int) string {

	// Define relevant character set
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Instantiate random number generator
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Draw from string charset
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

/*
Return datetime string in YYYYMMDD-HHMMSS format
 */
func GetDateTimeString() string {
	// Define datetime layout
	const layout = "20060102-150405"
	t := time.Now()
	return t.Format(layout)
}

/*
Creates output redirection for stdout and stderr to file (and console).
Restores original output association after call of returned function.
 */
func SaveOutput(filename string) func() {

	outfile := filename
	// Open file with read/write access; create if it does not exist, and clear file if it exists (overwrite),
	// Also, set general read/write permissions
	f, _ := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	// Save existing stdout and stderr for later reassignment
	origStdout := os.Stdout
	origStderr := os.Stderr

	// Create MultiWriter writing both to outfile and output
	mw := io.MultiWriter(origStdout, f)

	// Get pipe reader and writer for redirection
	r, w, _ := os.Pipe()

	// Redirect stdout and stderr to pipe writer
	os.Stdout = w
	os.Stderr = w

	// Also redirect log output
	log.SetOutput(mw)

	// Create channel to manage asynchronous writing of output
	exit := make(chan bool)

	go func() {
		// Copy all reads from pipe to MultiWriter - to ensure it is written to file and console
		_,_ = io.Copy(mw, r)
		// Set exit signal to terminate goroutine
		exit <- true
	}()

	// Function to be called once logging finishes (either direct call or defer)
	return func() {
		// Close writer, then block on exit channel; allows MultiWriter to finish operation before terminating
		_ = w.Close()
		// Reset stdout and stderr
		os.Stdout = origStdout
		os.Stderr = origStderr
		<-exit
		// Finally, close the outfile after writing has finished
		_ = f.Close()
	}

}
