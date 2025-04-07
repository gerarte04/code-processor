package generator

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func NewSessionId() string {
    b := make([]byte, 32)

    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return ""
    }

    return base64.URLEncoding.EncodeToString(b)
}
