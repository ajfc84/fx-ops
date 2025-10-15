package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// AskPassword prompts securely (no echo).
func AskPassword(prompt string) string {
	fmt.Printf("%s: ", prompt)
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // newline after input
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to read password: %v\n", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(pw))
}

// AskInput prompts for regular input (with echo).
func AskInput(prompt string) string {
	fmt.Printf("%s: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
