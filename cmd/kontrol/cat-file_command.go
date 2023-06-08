package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func catFileCmd(args []string) (err error) {

	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: kontrol cat-file <flag> <blob_hash>\n")
	}

	blobHash := args[2]
	flag := args[1]
	switch flag {
	case "-p":
		//length of blobhash in git is 40: 21be5453b240166dc5249e4878e71145ae55e126
		//while length of sha1.Size is 20
		if len(blobHash) != 2*sha1.Size {
			fmt.Errorf("not a valid object name: %v", blobHash)
		}
		//example file path: .git/objects/85/be5453b240166dc5247e4878e71145ae55e126
		path := filepath.Join(".git", "objects", blobHash[:2], blobHash[2:])
		file, err := os.Open(path)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("not a valid object name: %v", blobHash)
		}
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
		defer func() {
			e := file.Close()
			if err == nil && e != nil {
				err = fmt.Errorf("close file: %w", e)
			}
		}()
		return catFile(file)
	}
	return nil
}

// TODO
func catFile() {}

// TODO
func parseObject() {}
