package znet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"zframe/zinterface"
)

type PackageUti struct {
}

func NewPackageUti() *PackageUti {
	return &PackageUti{}
}

//得到数据包头的长度
//id|len 的总长度 取决于这两个的数据类型 当前用的都为uint64 所以为8+8
func (p *PackageUti) GetHeadLen() int {
	return 8
}

//
////或者分开写 id len分别写
////获取包头存放的id数据的长度
//func (p *PackageUti) GetHeadIDLen() int {
//	return 8
//}
//
////获取包头存放len(data)数据的长度
//func (p *PackageUti) GetHeadDataLen() int {
//	return 8
//}

//打包   len(data)|id|data
//Message类型为封装好的消息数据 包含消息id ，打包方法为：将消息的id 长度 和消息 打成二进制数据发送
func (p *PackageUti) Pack(msg zinterface.IMessage) ([]byte, error) {

	buf := bytes.NewBuffer([]byte{})


	//len(data)
	err := binary.Write(buf, binary.LittleEndian, msg.GetDataLength())
	if err != nil {
		fmt.Println("pack data len err:", err)
		return nil, err
	}
	//id
	err = binary.Write(buf, binary.LittleEndian, msg.GetID())
	if err != nil {
		fmt.Println("pack data id err:", err)
		return nil, err
	}
	//data
	err = binary.Write(buf, binary.LittleEndian, msg.GetData())
	if err != nil {
		fmt.Println("pack data err:", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

//拆包
//将读到的二级制数据id|len(data)| 拆包为Message类型 暂时无data 只有头数据
func (p *PackageUti) UnPack(data []byte) (zinterface.IMessage, error) {

	msg := &Message{}

	buf := bytes.NewBuffer(data)


	err := binary.Read(buf, binary.LittleEndian, &msg.DataLength)
	if err != nil {
		fmt.Println("unpack data length err:", err)
		return nil, err
	}

	//id
	err = binary.Read(buf, binary.LittleEndian, &msg.ID)
	if err != nil {
		fmt.Println("unpack data id err:", err)
		return nil, err
	}
	//fmt.Println("msg id = ", msg.ID)

	return msg, nil
}
