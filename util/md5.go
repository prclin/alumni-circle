package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
)

/*
	MD5工具函数
*/

// 计算字符串的MD5
func StringMD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// 计算文件的MD5
func FileMD5(file multipart.File) string {
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil))
}
