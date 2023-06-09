package main

import (
	"bufio"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func catFileCmd(args []string) (err error) {

	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: kontrol cat-file <flag> <blob_hash>\n")
	}

	blobHash := args[2]
	flag := args[1]
	switch flag {
	//TODO: case "-t":
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

// since the content of the file is compressed by zlib compression algorithm, we need to decompress it
func catFile(r io.Reader) (err error) {
	//read the compressed file (cf = compressed file)
	cf, err := zlib.NewReader(r)
	if err != nil {
		return fmt.Errorf("create zlib reader: %w", err)
	}
	//close the zlib reader
	defer func() {
		e := cf.Close()
		if err == nil && e != nil {
			err = fmt.Errorf("close zlib reader: %w", e)
		}
	}()
	//parse the content of the zlib reader
	err = parseObject(cf)
	if err != nil {
		return fmt.Errorf("parse object: %w", err)
	}
	return nil
}

func parseObject(r io.Reader) (err error) {
	br := bufio.NewReader(r)

	//ReadString(delimiter) reads until the first occurrence of delim in the input
	//returning a string containing the data up to and including the delimiter
	typ, err := br.ReadString(' ')
	if err != nil {
		return err
	}

	typ = typ[:len(typ)-1]
	if typ != "blob" {
		return fmt.Errorf("unsupported type %q", typ)
	}

	//'\000' represents the null character
	//example for sizeStr may be: 4242\000
	sizeStr, err := br.ReadString('\000')
	if err != nil {
		return err
	}
	sizeStr = sizeStr[:len(sizeStr)-1]

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return err
	}
	_, err = io.CopyN(os.Stdout, br, size)
	if err != nil {
		return err
	}
	return nil
}
