package util

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
)

// MD5 加密字符串
func MD5(buf []byte) string {
	hash := md5.New()
	sum := hash.Sum(buf)
	return hex.EncodeToString(sum[:])
}

// GenerateToken 生成token
func GenerateToken(claims model.TokenClaims) (string, error) {
	//加密
	key, err := parsePrivateKey([]byte(global.Configuration.Jwt.RSA.PrivateKey))
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	return token, err
}

// ParseToken 解析token
func ParseToken(tokenStr string) (*model.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parsePublicKey([]byte(global.Configuration.Jwt.RSA.PublicKey))
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("claims invalid")
	}

	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}

	return claims, nil
}

// parsePrivateKey 解析PKCS1格式私钥
func parsePrivateKey(buf []byte) (*rsa.PrivateKey, error) {
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse private key error")
	}
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

// parsePublicKey 解析PKCS1格式公钥
func parsePublicKey(buf []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, errors.New("block nil")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}
