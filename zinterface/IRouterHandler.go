package zinterface

//路由管理 抽象层
type IRouterHandler interface {

	//
	GetRouterHandler() map[uint32]IRouter

	//添加路由
	AddRouter(msgID uint32, router IRouter)

	//执行路由业务
	DoRouter(request IRequest)

	//初始化任务池
	InitMissionPool()

	//将消息添加到池子
	AddMsgToPool(request IRequest)
}
