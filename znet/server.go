package znet

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"zframe/config"
	"zframe/zinterface"
)

//实现层

var ServerInstance zinterface.IServer
var once sync.Once

type Server struct {
	//ip版本
	NetProtocol string

	IP string

	Port int
	//服务器名称
	Name string
	//路由  路由管理者
	RouterHandler zinterface.IRouterHandler

	//连接管理模块
	connManager zinterface.IConnectionManager

	onConnStart func(conn zinterface.IConnection)
	onConnStop  func(conn zinterface.IConnection)
}

func NewServer() zinterface.IServer {
	config.InitLoad()

	if ServerInstance == nil {

		once.Do(
			func() {
			ServerInstance = &Server{
				Name:          config.GlobalObj.Name,
				IP:            config.GlobalObj.IP,
				Port:          config.GlobalObj.Port,
				NetProtocol:   config.GlobalObj.NetProtocol,
				RouterHandler: NewRouterHandler(),
				connManager:   NewConnectionManager(),
			}
		})

	}
	return ServerInstance
}

//回执请求
//func callback(conn *net.TCPConn, data []byte,count int) error {
//func callback(request zinterface.IRequest) error {
//
//	conn := request.GetConnection().GetTCPConnection()
//	data := request.GetData()
//	count := request.GetDataLength()
//
//	_, err := conn.Write(data[:count])
//	if err != nil {
//		return err
//	}
//	return nil
//
//}

func (s *Server) Start() {

	fmt.Println("server listen at ip =", s.IP, "port=", s.Port)

	//1.
	tcpAddr, err := net.ResolveTCPAddr(s.NetProtocol, s.IP+":"+strconv.Itoa(s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr err:", err)
		return
	}
	//2.开始监听
	tcpLis, err := net.ListenTCP(s.NetProtocol, tcpAddr)
	if err != nil {
		fmt.Println("listen tcp err:", err)
		return
	}
	fmt.Println("current MaxConnection is ", config.GlobalObj.MaxConnection)
	//初始化 任务池
	s.RouterHandler.InitMissionPool()

	var id uint32 = 1

	//3.等待客户端连接
	go func() {

		for {
			//conn, err := tcpLis.Accept()
			conn, err := tcpLis.AcceptTCP()

			if err != nil {
				fmt.Println("accept err:", err)
				continue
			}
			//此时客户端已连接
			fmt.Println("客户端" + conn.RemoteAddr().String() + "连接成功")
			//连接管理模块 判断是否超出最大连接数
			if s.connManager.GetLen() >= config.GlobalObj.MaxConnection {
				fmt.Println("more than MaxConnection :", config.GlobalObj.MaxConnection)
				_ = conn.Close()
				return
			}

			//小于最大可连接数  可以连接
			currentConn := NewConnection(conn, id, s.RouterHandler)
			id++

			s.connManager.AddConn(currentConn)

			//链接开始工作   读 和 写业务
			currentConn.Start()

			//go func() {
			//	cli := "[" + conn.RemoteAddr().String() + "]"
			//	fmt.Println("客户端：", cli, "已连接")
			//	bs := make([]byte, 4096)
			//	//
			//	for {
			//		n, err := conn.Read(bs)
			//		if n == 0 {
			//			fmt.Println("客户端", cli, "断开了连接")
			//			return
			//		}
			//		if err != nil && err != io.EOF {
			//			fmt.Println("conn read err:", err)
			//			continue
			//		}
			//		fmt.Println("读到了客户端", cli, "的消息:", string(bs[:n]))
			//		n, err = conn.Write(bs[:n])
			//		if err != nil {
			//			fmt.Println("conn write err:", err)
			//			continue
			//		}
			//
			//	}
			//
			//}()

		}

	}()

}

func (s *Server) Stop() {

	//服务器停止 清空所有连接
	s.connManager.Clear()
}

func (s *Server) Run() {

	//开启

	s.Start()

	//TODO 扩展内容

	select {}

}
func (s *Server) AddRouter(msgID uint32, router zinterface.IRouter) {
	//s.RouterHandler = router

	s.RouterHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnManager() zinterface.IConnectionManager {
	return s.connManager
}

//注册连接函数
func (s *Server) RegisterOnConnStart(hookFunc func(conn zinterface.IConnection)) {
	s.onConnStart = hookFunc
}

//注册结束函数
func (s *Server) RegisterOnConnStop(hookFunc func(conn zinterface.IConnection)) {
	s.onConnStop = hookFunc
}

//执行
func (s *Server) ExecuteOnConnStart(conn zinterface.IConnection) {
	if s.onConnStart != nil {
		fmt.Println("execute on conn start func")
		s.onConnStart(conn)

	}
}
func (s *Server) ExecuteOnConnStop(conn zinterface.IConnection) {
	if s.onConnStop != nil {
		fmt.Println("execute on conn stop func")
		s.onConnStop(conn)

	}
}
