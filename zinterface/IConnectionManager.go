package zinterface

//连接管理 抽象层
type IConnectionManager interface {
	//add
	AddConn(conn IConnection)

	//remove
	RemoveConn(connID uint32)

	//get by id
	GetConn(connID uint32) (IConnection, error)

	//get len
	GetLen() uint32

	//removeall
	Clear()
}
