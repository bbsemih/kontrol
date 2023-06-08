package main

import (
	"fmt"
	"os"
)

func initCmd() error {
	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %v: %v", dir, err)
		}
	}
	//output of 'cat .git/HEAD' is: 'ref: refs/heads/main'
	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
		return fmt.Errorf("error ocured writing '.git/HEAD': %v", err)
	}
	fmt.Println("Initialized a git repository")
	return nil
}
