package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"zframe/config"
	"zframe/zinterface"
)

//具体的TCP连接模块
type Connection struct {
	//当前链接的原生socket
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前链接状态
	IsClosed bool

	////当前链接绑定的业务处理方法
	//HandleAPI zinterface.HandleFunc
	RouterHandler zinterface.IRouterHandler

	//通知writer发送消息的通道
	toWriterChan chan []byte

	//reader通知writer退出的通道
	writerExitCh chan bool

	//属性
	propertyMap map[string]interface{}

	propertyLock sync.RWMutex
}

//初始化链接方法

func NewConnection(conn *net.TCPConn, connID uint32, handler zinterface.IRouterHandler) zinterface.IConnection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		IsClosed: false,
		//HandleAPI: callback,
		RouterHandler: handler,
		toWriterChan:  make(chan []byte),
		writerExitCh:  make(chan bool),
		propertyMap:   make(map[string]interface{}),
	}

}

func (c *Connection) StartRead() {

	//fmt.Println("id=", c.ConnID, "start read")

	//bs := make([]byte, 4096)

	defer func() {
		//fmt.Println("id=", c.ConnID, "conn read exit")
		c.Stop()
	}()
	pack := NewPackageUti()
	for {

		//读数据

		headLen := pack.GetHeadLen() //头长度

		buf := make([]byte, headLen) //头数据

		n, err := io.ReadFull(c.Conn, buf) //读满buf的长度
		if n == 0 {
			fmt.Println("id=", c.ConnID, "的客户端"+c.GetRemoteAddr().String()+"断开了链接")
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("id=", c.ConnID, "read err:", err)
			continue
		}

		imsg, err := pack.UnPack(buf)
		if err != nil {
			fmt.Println("id=", c.ConnID, "unpack err:", err)
			return
		}
		var data []byte
		if imsg.GetDataLength() > 0 {
			//进行数据操作
			data = make([]byte, imsg.GetDataLength())
			n, err := io.ReadFull(c.Conn, data)
			if n == 0 {
				fmt.Println("id=", c.ConnID, "的客户端"+c.GetRemoteAddr().String()+"断开了链接")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("id=", c.ConnID, "read err:", err)
				continue
			}

		}
		imsg.SetData(data)
		//读完成  传递给请求	(request)

		//req := NewRequest(c, bs, n)
		req := NewRequest(c, imsg.(*Message))

		//添加到任务池
		if config.GlobalObj.MissionPoolCount > 0 {
			c.RouterHandler.AddMsgToPool(req)
		} else {
			//直接执行
			go c.RouterHandler.DoRouter(req)
		}

		//if err != nil {
		//	fmt.Println("id=", c.ConnID, "handle回执给客户端消息错误:", err)
		//	continue
		//}

	}
}

//写给客户端
func (c *Connection) StartWrite() {

	//defer fmt.Println("id=", c.ConnID, "conn write exit")

	for {
		select {
		case data := <-c.toWriterChan:
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("id=", c.ConnID, "write err:", err)
				return
			}
		case <-c.writerExitCh:
			return
		}
	}

}

//启动链接
func (c *Connection) Start() {
	//执行hook On start
	go ServerInstance.ExecuteOnConnStart(c)

	fmt.Println("id=", c.ConnID, "start conn")
	//读 业务
	go c.StartRead()

	// 写 业务
	go c.StartWrite()

}

//停止链接
func (c *Connection) Stop() {

	if c.IsClosed {
		fmt.Println("id=", c.ConnID, "conn has been stopped")
		return
	}

	//关闭后 执行的函数 wait hook exe done
	ServerInstance.ExecuteOnConnStop(c)

	c.IsClosed = true

	//通知writer退出
	c.writerExitCh <- true

	_ = c.Conn.Close()

	//从连接管理模块 移除当前连接 conn
	ServerInstance.GetConnManager().RemoveConn(c.ConnID)

	//回收通道
	close(c.writerExitCh)
	close(c.toWriterChan)
	fmt.Println("id=", c.ConnID, "conn stop")
}

//获取链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取conn的原生socket套接字
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取远程客户端的ip地址
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

////发送数据给对方客户端
//func (c *Connection) Send(data []byte, count int) (int, error) {
//	n, err := c.Conn.Write(data[:count])
//	if err != nil {
//		return 0, err
//	}
//	return n, err
//}

//发送数据给对方客户端
func (c *Connection) Send(msgID uint32, data []byte) error {

	if c.IsClosed {
		fmt.Println("conn has been closed")
		return errors.New("conn id = " + strconv.Itoa(int(c.ConnID)) + " closed when send msg")
	}
	//封装message
	msg := NewMessage(msgID, data)

	pack := NewPackageUti()

	packData, err := pack.Pack(msg)
	if err != nil {
		fmt.Println("msg id =", msgID, "conn send : pack data err:", err)
		return err
	}

	//n, err := c.Conn.Write(packData)
	//if err != nil {
	//	return 0, err
	//}
	c.toWriterChan <- packData
	return nil
}

//添加属性
func (c *Connection) AddProperty(key string, i interface{}) {
	//_, err := c.GetProperty(key)
	//if err != nil {
	//	c.propertyLock.Lock()
	//	defer c.propertyLock.Unlock()
	//	c.propertyMap[key] = i
	//}else{
	//	fmt.Println()
	//}
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.propertyMap[key] = i

}

//获取属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	v, ok := c.propertyMap[key]
	if ok {
		return v, nil
	} else {
		return nil, errors.New(fmt.Sprint("conn id = ", c.ConnID, "do not have the property which key = ", key))
	}
}

//删除属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.propertyMap, key)
}
