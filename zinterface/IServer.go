package zinterface

//抽象层

//
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Run()
	//添加路由
	AddRouter(msgID uint32, router IRouter)
	//连接管理模块
	GetConnManager() IConnectionManager

	//注册连接函数
	RegisterOnConnStart(hookFunc func(conn IConnection))
	//注册结束函数
	RegisterOnConnStop(hookFunc func(conn IConnection))
	//执行
	ExecuteOnConnStart(conn IConnection)
	ExecuteOnConnStop(conn IConnection)
}
