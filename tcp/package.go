package tcp

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/go-eyas/toolkit/util"
)

// 一个默认的私有协议实现

// 协议组成
// 4bt(自定义数据长度) + 任意bt(json字符串数据)
// json 格式 {"cmd": "test", "data": {}}

// 打包
func Packer(data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	bodyLen := uint32(len(body))
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, bodyLen)
	// return header, nil
	// header := bytes.NewBuffer(make([]byte, 4))
	// _ = binary.Write(header, binary.BigEndian, bodyLen)
	pkg := util.BytesCombine(header, body)
	return pkg, nil
}

// 解包
var parserBuf = make(map[*Conn][]byte)

func Parser(conn *Conn, bt []byte) (interface{}, error) {
	preBuf, ok := parserBuf[conn]
	if !ok {
		preBuf = make([]byte, 0)
		parserBuf[conn] = preBuf
	}
	buf := util.BytesCombine(preBuf, bt)
	if len(buf) < 4 {
		parserBuf[conn] = util.BytesCombine(parserBuf[conn], bt)
		return nil, errors.New("half pack")
	}
	header := buf[:4]
	bodyLen := binary.BigEndian.Uint32(header)
	if uint32(len(buf)) < (4 + bodyLen) {
		parserBuf[conn] = util.BytesCombine(parserBuf[conn], bt)
		return nil, errors.New("half pack")
	}
	body := buf[4 : 4+bodyLen]
	parserBuf[conn] = buf[4+bodyLen:]
	return body, nil
}
