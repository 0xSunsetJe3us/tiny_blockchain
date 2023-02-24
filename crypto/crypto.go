// Package crypto provides encrypts methods.
package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"strings"
)

func NewHashHandler(hashType string) func() hash.Hash {
	switch strings.ToLower(hashType) {
	case "md5":
		return md5.New
	case "sha1":
		return sha1.New
	case "sha256":
		return sha256.New
	case "sha512":
		return sha512.New
	default:
		return sha1.New
	}
}
