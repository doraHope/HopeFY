package util

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func CreateUUID() string {
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(id) //base64, 32*8 => 43*8
}
