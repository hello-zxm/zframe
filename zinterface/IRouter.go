package zinterface

//抽象 路由层

type IRouter interface {

	//业务之前
	BeforeHandle(request IRequest)

	//处理的业务
	Handle(request IRequest)

	//业务之后
	AfterHandle(request IRequest)


}
