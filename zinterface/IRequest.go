package zinterface

//抽象 一次性请求的数据

type IRequest interface {
	//当前请求链接
	GetConnection() IConnection

	////当前的数据
	//GetData() []byte
	//
	////当前数据长度
	//GetDataLength() int

	//数据
	GetMsg() IMessage
}
