package hash

import (
	"crypto/md5"
	"encoding/hex"
)

type MD5Hash struct{}

func NewMD5() MD5Hash {
	return MD5Hash{}
}

func (t MD5Hash) generate(msg []byte) string {
	hash := md5.Sum(msg)
	return hex.EncodeToString(hash[:])
}
