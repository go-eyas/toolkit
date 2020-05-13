package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"

	"golang.org/x/crypto/bcrypt"
)

func BcryptHash(src string) string {
	bt, _ := bcrypt.GenerateFromPassword([]byte(src), bcrypt.MinCost)
	return string(bt)
}

func BcryptVerify(hash, src string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(src)) == nil
}

func AesEncrypt(src, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	src = _padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return src
}

func AesDecrypt(src []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = _unpadding(src)
	return src
}

func _padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func _unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}
