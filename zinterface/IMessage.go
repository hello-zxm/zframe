package zinterface

//消息  抽象层
type IMessage interface {
	//消息id (类型）
	GetID() uint32

	//数据
	GetData() []byte

	//长度
	GetDataLength() uint32

	SetID(uint32)
	SetData([]byte)
	SetDataLength(uint32)
}
