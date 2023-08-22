package util

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 加密字符串
func MD5(buf []byte) string {
	hash := md5.New()
	sum := hash.Sum(buf)
	return hex.EncodeToString(sum[:])
}
