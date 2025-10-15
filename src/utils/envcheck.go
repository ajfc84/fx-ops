package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func CheckEnvironment() {
	var (
		socket,
		home string
	)
	switch runtime.GOOS {
	case "linux":
		socket = "/var/run/docker.sock"
		home = os.Getenv("HOME")
	case "windows":
		socket = `//./pipe/docker_engine`
		home = os.Getenv("USERPROFILE")
	default:
		panic("unsupported OS: " + runtime.GOOS)
	}

	if _, err := os.Stat(socket); err != nil {
		fmt.Printf("ERROR: Docker socket not found at %s\n", socket)
		os.Exit(1)
	}

	sshDir := filepath.Join(home, ".ssh")
	if _, err := os.Stat(sshDir); err != nil {
		fmt.Printf("ERROR: SSH directory not found at %s\n", sshDir)
		os.Exit(1)
	}
}
