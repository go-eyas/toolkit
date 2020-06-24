package tcp

import (
	"encoding/binary"
	"encoding/json"
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

	pkg := util.BytesCombine(header, body)
	return pkg, nil
}

func Parser() (map[uint64][]byte, func(conn *Conn, bt []byte) ([]interface{}, error)) {
	// 解包
	var parserBuf = make(map[uint64][]byte)
	return parserBuf, func(conn *Conn, bt []byte) ([]interface{}, error) {
		preBuf, ok := parserBuf[conn.ID]
		if !ok {
			preBuf = make([]byte, 0)
			parserBuf[conn.ID] = preBuf
		}

		buf := util.BytesCombine(preBuf, bt)
		datas := make([]interface{}, 0)

		for {
			if len(buf) < 4 {
				break
			}
			header := buf[:4]
			bodyLen := binary.BigEndian.Uint32(header)
			if uint32(len(buf)) < (4 + bodyLen) {
				break
			}
			pack := buf[4 : 4+bodyLen]
			buf = buf[4+bodyLen:]
			datas = append(datas, pack)
		}
		parserBuf[conn.ID] = buf

		return datas, nil
	}
}
