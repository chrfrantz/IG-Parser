package helper

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
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
Returned function allows specification of suffix appended to filename.
 */
func SaveOutput(filename string) (func(string) error, error) {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		folderEnd := strings.LastIndex(filename, "/")
		if folderEnd != -1 {
			err = os.MkdirAll(filename[:folderEnd], 0700)
			if err != nil {
				log.Println("Failed to create folder " + filename[:folderEnd] +
					", Error: " + err.Error())
			} else {
				log.Println("Created folder " + filename[:folderEnd])
			}
		}
	}

	// Reassign for flexibility
	outfile := filename

	log.Println("Log file: " + outfile)
	// Open file with read/write access; create if it does not exist, and clear file if it exists (overwrite),
	// Also, set general read/write permissions
	f, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	// Save existing stdout and stderr for later reassignment
	origStdout := os.Stdout
	origStderr := os.Stderr

	// Create MultiWriter writing both to outfile and output
	mw := io.MultiWriter(origStdout, f)

	// Get pipe reader and writer for redirection
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

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
	return func(suffix string) error {
		// Close writer, then block on exit channel; allows MultiWriter to finish operation before terminating
		err = w.Close()

		// Reset stdout and stderr
		os.Stdout = origStdout
		os.Stderr = origStderr

		// Write any error to standard out to review (only here, since the file may not be writable)
		if err != nil {
			fmt.Println("Error during closing writer: " + err.Error())
			return err
		}

		// Send exit signal to channel
		<-exit

		// Finally, close the outfile after writing has finished
		err = f.Close()
		// Write any error to standard out to review
		if err != nil {
			fmt.Println("Error during closing log file: " + err.Error())
			return err
		}

		// Rename file if suffix is provided
		if suffix != "" {
			err = os.Rename(f.Name(), f.Name()+suffix)
			// Write any error to standard out to review
			if err != nil {
				fmt.Println("Error during renaming of log file: " + err.Error())
				return err
			}
		}

		// No error
		return nil

	}, nil

}
