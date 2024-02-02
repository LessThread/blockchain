package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

// 全局错误处理
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// int64转字节类型
func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// 检查文件存在性
func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}
