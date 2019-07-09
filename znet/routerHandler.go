package znet

import (
	"fmt"
	"zframe/config"
	"zframe/zinterface"
)

//路由管理/分发
type RouterHandler struct {
	//key:路由id(类型) value:路由
	RouterHandler map[uint32]zinterface.IRouter

	//任务池
	//按照连接顺序(conn id) 加入任务池，保证同一连接的请求为队列形式
	//后期要改
	MissionPool []chan zinterface.IRequest

	//任务池长度
	PoolSize uint32
}

func NewRouterHandler() zinterface.IRouterHandler {
	return &RouterHandler{
		RouterHandler: make(map[uint32]zinterface.IRouter),
		MissionPool:   make([]chan zinterface.IRequest, config.GlobalObj.MissionPoolCount),
		PoolSize:      config.GlobalObj.MissionPoolCount,
	}
}

func (r *RouterHandler) GetRouterHandler() map[uint32]zinterface.IRouter {
	return r.RouterHandler
}

//添加路由
func (r *RouterHandler) AddRouter(msgID uint32, router zinterface.IRouter) {
	_, ok := r.RouterHandler[msgID]
	if ok {
		fmt.Println("msg id=", msgID, "has been in router handler")
		return
	}
	r.RouterHandler[msgID] = router
}

//执行路由业务
func (r *RouterHandler) DoRouter(request zinterface.IRequest) {
	router, ok := r.RouterHandler[request.GetMsg().GetID()]
	if !ok {
		fmt.Println("msg id=", request.GetMsg().GetID(), "was not in router handler")
		return
	}
	//fmt.Println(request.GetConnection().GetConnID(), "th mission is going to execute")
	router.BeforeHandle(request)
	router.Handle(request)
	router.AfterHandle(request)
}

//初始化任务池  开启服务 初始化任务池
func (r *RouterHandler) InitMissionPool() {
	fmt.Println("mission pool size=", r.PoolSize)
	for i := uint32(0); i < r.PoolSize; i++ {

		r.MissionPool[i] = make(chan zinterface.IRequest,config.GlobalObj.MissionItemCount)
		fmt.Println(i+1, "th mission is starting")
		//等待消息 收到消息后执行路由业务
		go func(item chan zinterface.IRequest) {

			for {

				select {
				case req := <-item:
					r.DoRouter(req)
				}
			}

		}(r.MissionPool[i])

	}
}

//将消息添加到池子
func (r *RouterHandler) AddMsgToPool(request zinterface.IRequest) {

	//连接id
	connID := request.GetConnection().GetConnID()%r.PoolSize - 1 //conn id 从1开始的  要减1

	r.MissionPool[connID] <- request

}
