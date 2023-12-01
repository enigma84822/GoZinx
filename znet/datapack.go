package znet

import (
	"GoZinx/utils"
	"GoZinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

// 封包、拆包
type DataPack struct {
}

// 获取包的头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

// 封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen 写进buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 讲MsgId 写进buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 讲data数据 写进buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// ioReader
	dataBuff := bytes.NewBuffer(binaryData)

	// 解压
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	// 判断数据长度是否超出最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}
	return msg, nil
}

func NewDataPack() *DataPack {
	return &DataPack{}
}
