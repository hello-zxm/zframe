package zinterface

import "net"

//抽象链接层

type IConnection interface {
	//启动链接
	Start()

	//停止链接
	Stop()

	//获取链接ID
	GetConnID() uint32

	//获取链接的原生socket
	GetTCPConnection() *net.TCPConn

	//获取远程客户端的ip
	GetRemoteAddr() net.Addr

	//给客户端发送数据
	//Send(data []byte,count int) (int, error)
	Send(msgID uint32, data []byte) error

	//添加属性
	AddProperty(key string, i interface{})

	//获取属性
	GetProperty(key string) (interface{}, error)

	//删除属性
	RemoveProperty(key string)
}

//抽象定义 业务处理方法
//参1 socket 参2 数据  参3 数据长度
//要用户自己实现
//type HandleFunc func(request IRequest) error
