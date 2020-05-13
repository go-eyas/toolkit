package util

import (
	"bytes"
	"io"
	"io/ioutil"
)

// ByteToReader 将字节转换成读取流
func ByteToReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}

// ByteToReadCloser 将字节转换成一次性的读取流
func ByteToReadCloser(b []byte) io.ReadCloser {
	return ioutil.NopCloser(ByteToReader(b))
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
