package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func hashObjectCmd(args []string) (err error) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: kontrol hash-object <flag> <file>\n")
	}
	flag := args[1]
	filePath := args[2]
	switch flag {
	//Steps in command: 'kontrol hash-object -w <file>':
	//find the file in the .git/objects directory
	//-w is for writing the blob to disk
	//as a result of the command, a 40-char SHA is printed to stdout and the blob is written to disk
	//written to disk as a file with the first two characters of the SHA as the directory name
	case "-w":
		fileInfo, err := os.Stat(filePath)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file %v does not exist", filePath)
		}
		if err != nil {
			return fmt.Errorf("error occured reading file info: %v", err)
		}
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error occured opening file: %v", err)
		}
		defer func() {
			e := file.Close()
			if e != nil && err == nil {
				err = fmt.Errorf("error occured closing file: %v", e)
			}
		}()

		name, err := hashObject(file, "blob", fileInfo.Size())
		if err != nil {
			return fmt.Errorf("error occured hashing object: %v", err)
		}
		fmt.Println(name)
	}
	return nil
}

// result of cat .git/objects/ca/952bd58caeb705520df8574cd2cef51744f6a2 is:
// xC��blob 59xK��OR01f��O)�IUH�,�(M�K���OJ*N����/J-JM�,��J�W0�32���� we have to encode this
func hashObject(src io.Reader, typ string, size int64) (string, error) {
	var buff bytes.Buffer

	err := encodeObject(&buff, size, src)

	fileContent, err := compress(buff.Bytes())
	if err != nil {
		return "", fmt.Errorf("error occured compressing file: %v", err)
	}
	//sha1.Sum calculates the SHA-1 hash value of a given input
	sum := sha1.Sum(fileContent)
	name := hex.EncodeToString(sum[:])

	//created file name example: .git/objects/ra/552bd58c1eb703320df8574cd2cef51744f6a2
	objectPath := filepath.Join(".git", "objects", name[:2], name[2:])
	dirPath := filepath.Dir(objectPath)

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", fmt.Errorf("error occured creating directory: %v", err)
	}
	err = os.WriteFile(objectPath, fileContent, 0644)
	if err != nil {
		return "", fmt.Errorf("error occured writing file: %v", err)
	}
	return name, nil
}

func encodeObject(typ string, size int64, src io.Reader) (string, error) {
	// TODO
}

func compress(data []byte) ([]byte, error) {
}
