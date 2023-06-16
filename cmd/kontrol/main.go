package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: kontrol <command> [<args>...]\n")
		os.Exit(1)
	}

	var err error

	switch command := os.Args[1]; command {
	case "init":
		err = initCmd()
	case "cat-file":
		err = catFileCmd(os.Args[1:])
	//compute object id and optionally creates a blob from a file
	//'kontrol hash-object -w <file>' -w is for writing the blob to disk
	case "hash-object":
		err = hashObjectCmd(os.Args[1:])
	case "ls-tree":
		err = lsTreeCmd(os.Args[1:])
	case "commit-tree":
		err = commitTreeCmd(os.Args[1:])
	default:
		err = fmt.Errorf("unknown command %s", command)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
}
