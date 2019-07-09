package znet

import (
	"zframe/zinterface"
)

//实现路由

type Router struct {
}

//业务之前
func (r *Router) BeforeHandle(request zinterface.IRequest) {

}

//处理的业务
func (r *Router) Handle(request zinterface.IRequest) {

}

//业务之后
func (r *Router) AfterHandle(request zinterface.IRequest) {

}
