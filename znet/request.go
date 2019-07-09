package znet

import "zframe/zinterface"

//一次请求数据 实现

type Request struct {
	//当前链接
	Conn zinterface.IConnection

	////数据
	//Data []byte
	//
	////长度
	//DataLength int
	Msg zinterface.IMessage
}

func NewRequest(conn zinterface.IConnection, msg zinterface.IMessage) zinterface.IRequest {
	return &Request{
		Conn: conn,
		Msg:  msg,
	}
}

//func NewRequest(conn zinterface.IConnection, data []byte, len int) zinterface.IRequest {
//	return &Request{
//		Conn:       conn,
//		Data:       data,
//		DataLength: len,
//	}
//}

//当前请求链接
func (r *Request) GetConnection() zinterface.IConnection {
	return r.Conn
}

func (r *Request) GetMsg() zinterface.IMessage {
	return r.Msg
}

////当前的数据
//func (r *Request) GetData() []byte {
//	return r.Data
//}
//
////当前数据长度
//func (r *Request) GetDataLength() int {
//	return r.DataLength
//}
