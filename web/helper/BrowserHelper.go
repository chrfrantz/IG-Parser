package helper

import (
	"os/exec"
	"runtime"
)

/*
This file contains OS-independent functionality to start web browsers upon service instantiation
when deploying locally.
*/

/*
Attempts to open a browser for a given URL on selected operating systems.
Returns error if launch fails.
Based on code provided under: https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8
*/
func OpenBrowser(url string) error {
	var cmd string
	var args []string

	// Switch depending on operating system
	switch runtime.GOOS {
	case "windows": // Windows
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin": // Mac OS
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	// Launch browser
	return exec.Command(cmd, args...).Start()
}
