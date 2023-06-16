package main

import (
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bbsemih/kontrol/object"
)

// 'kontrol ls-tree --name-only <tree-ish>':
// the contents of the tree object are the mode, type, object id, and file name
func lsTreeCmd(args []string) (err error) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: kontrol ls-tree <flag> <tree-ish>\n")
	}
	flag := args[1]
	treeIsh := args[2]

	switch flag {
	case "--name-only":
		//find the tree object in the .git/objects directory
		hash, err := object.HashFromString(treeIsh)
		if err != nil {
			return fmt.Errorf("hash from string: %w", err)
		}
		name := hash.String()
		tObj := filepath.Join(".git", "objects", name[:2], name[2:])
		file, err := os.Open(tObj)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("not a valid object name", tObj)
		}
		if err != nil {
			return fmt.Errorf("error occured opening file: %v", err)
		}
		defer func() {
			e := file.Close()
			if e != nil && err == nil {
				err = fmt.Errorf("error occured closing file: %v", e)
			}
		}()
		return Load(file)
	}
	return nil
}

func Load(r io.Reader) (typ string, content []byte, err error) {
	zr, err := zlib.NewReader(r)
	if err != nil {
		return "", nil, fmt.Errorf("new reader: %w", err)
	}

	defer func() {
		e := zr.Close()
		if err == nil && e != nil {
			err = fmt.Errorf("close: %w", e)
		}
	}()

	typ, content, err = object.ParseObject(zr)
	if err != nil {
		return "", nil, fmt.Errorf("parse object: %w", err)
	}
	return typ, content, nil
}
