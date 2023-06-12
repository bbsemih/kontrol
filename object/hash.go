package object

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Hash [sha1.Size]byte

func HashFromString(s string) (Hash, error) {
	if len(s) != 2*len(Hash{}) {
		return Hash{}, fmt.Errorf("not a valid object name")
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return Hash{}, fmt.Errorf("not a valid object name")
	}

	var h Hash
	copy(h[:], b)

	return h, nil
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}
