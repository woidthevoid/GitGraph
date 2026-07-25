package main

import (
	"flag"
)

// Scan function that takes a folder path as argument.
func scan(path string) {
	print("scanning...")
}

// Generates graph based on email.
func stats(email string) {
	print("stats")
}

func main() {
	var folder string
	var email string
	flag.StringVar(&folder, "add", "", "Scans a path input for git repositories")
	flag.StringVar(&email, "email", "example@email.com", "Your email")
	flag.Parse()

	if folder != "" {
		scan(folder)
		return
	}

	stats(email)
}
