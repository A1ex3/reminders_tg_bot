package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

type Sha struct{}

func (sha *Sha) Sha256(value string) string {
	h := sha256.New()

	h.Write([]byte(value))

	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
