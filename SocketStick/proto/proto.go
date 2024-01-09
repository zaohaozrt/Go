package proto

import (
	"bufio"
	"bytes"           //bytes实现了操作字节切片([]byte)的函数。它类似于字符串包的功能。
	"encoding/binary" //数字和字节序列之间的转换以及变量的编码和解码。
)

// 字符串加包头编码

func Encode(message string) ([]byte, error) {
	//读取消息字节长度，转换为int32类型
	var length = int32(len(message)) //四个字节
	var pkg = new(bytes.Buffer)      //建立一个空字节缓冲区,缓冲区内容可追加
	//写消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	//写入消息实体
	// err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	// if err != nil {
	// 	return nil, err
	// }
	// return pkg.Bytes(), nil //从字节缓冲区里面把字节读出来
	pkg.Write([]byte(message)) //等同上面
	return pkg.Bytes(), nil
}

// 解码消息
func Decode(reader *bufio.Reader) (string, error) {
	//读取消息长度
	//在不推进读取器的情况下返回下一个n字节
	lengthByte, _ := reader.Peek(4)           //窥视,读取前4个字节
	lengthBuff := bytes.NewBuffer(lengthByte) //建立一个字节缓冲区，用lengthByte初始化
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length) //读取包头信息
	if err != nil {
		return "", err
	}
	// Buffered返回缓冲中现有的可读取的字节数。
	// 判断内容存取是否有错误
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	//读取真正的消息
	pack := make([]byte, int(4+length)) //每次只读一次个包的信息，所以能解决TCP粘连
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil

}
