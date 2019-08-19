# zframe
go frame for tcp
                                                    made by @ZhaoXiming

type TestRouter struct {
	znet.Router
}

//
func (r *TestRouter) BeforeHandle(request zinterface.IRequest) {
	//fmt.Println("before handle")
}

//
func (r *TestRouter) Handle(request zinterface.IRequest) {
    // get data
    //data := request.GetMsg().GetData()
    // send data
	err := request.GetConnection().Send(2, []byte("hello zframe"))

	if err != nil {
		fmt.Println("send err:", err)
		return
	}

	//data := request.GetMsg().GetData()

	//fmt.Println("收到消息=", string(data))

	fmt.Println("handle")
}

//
func (r *TestRouter) AfterHandle(request zinterface.IRequest) {
	//fmt.Println("after handle")

}


func onStart(conn zinterface.IConnection) {
	fmt.Println(conn.GetConnID(), "开始了连接")

	conn.AddProperty("one", "one")
	conn.AddProperty("two", 2)

}
func onStop(conn zinterface.IConnection) {
	fmt.Println(conn.GetConnID(), "结束了连接")

	v1, _ := conn.GetProperty("one")
	v2, _ := conn.GetProperty("two")
	fmt.Println(v1)
	fmt.Println(v2)
}

func main() {

    //create a new server
	s := znet.NewServer()

    //add router with type
	s.AddRouter(1, &TestRouter{})

    //register on start func
	s.RegisterOnConnStart(onStart)
	//register on stop func
	s.RegisterOnConnStop(onStop)

    //run the server
	s.Run()

}