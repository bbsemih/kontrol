package object

import (
	"bufio"
	"io"
	"strconv"
)

func ParseObject(r io.Reader) (string, []byte, error) {
	br := bufio.NewReader(r)
	typ, err := br.ReadString(' ')
	if err != nil {
		return "", nil, err
	}
	typ = typ[:len(typ)-1]
	sizeStr, err := br.ReadString('\000')
	if err != nil {
		return "", nil, err
	}

	sizeStr = sizeStr[:len(sizeStr)-1] // cut '\000'
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return "", nil, err
	}
	content := make([]byte, size)
	_, err = io.ReadFull(br, content)
	if err != nil {
		return "", nil, err
	}
	return typ, content, nil
}
