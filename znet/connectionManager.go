package znet

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"zframe/config"
	"zframe/zinterface"
)

//连接管理
type ConnectionManager struct {
	//manager
	ConnManager map[uint32]zinterface.IConnection

	lock sync.RWMutex
}

func NewConnectionManager() zinterface.IConnectionManager {
	return &ConnectionManager{
		ConnManager: make(map[uint32]zinterface.IConnection),
	}
}

//add
func (c *ConnectionManager) AddConn(conn zinterface.IConnection) {

	c.lock.Lock()

	_, ok := c.ConnManager[conn.GetConnID()]
	if ok {
		//已存在
		fmt.Println("conn add failed , id=", conn.GetConnID(), "has been in manager")
		return
	} else {
		//不存在
		c.ConnManager[conn.GetConnID()] = conn
	}
	c.lock.Unlock()

	fmt.Println("current connection is ",c.GetLen(),"max is ",config.GlobalObj.MaxConnection)

}

//remove
func (c *ConnectionManager) RemoveConn(connID uint32) {

	c.lock.Lock()


	_, ok := c.ConnManager[connID]
	if ok {
		//存在
		delete(c.ConnManager, connID)

	} else {
		//不存在
		fmt.Println("conn delete failed ,id=", connID, "not in manager")
	}
	c.lock.Unlock()

	fmt.Println("current connection is ",c.GetLen(),"max is ",config.GlobalObj.MaxConnection)
}

//get by id
func (c *ConnectionManager) GetConn(connID uint32) (zinterface.IConnection, error) {

	c.lock.RLock()
	defer c.lock.RUnlock()

	conn, ok := c.ConnManager[connID]
	if ok {
		//存在
		return conn, nil

	} else {
		//不存在
		return nil, errors.New("conn get failed ,id=" + strconv.FormatInt(int64(connID), 10) + "not in manager")
	}

}

//get len
func (c *ConnectionManager) GetLen() uint32 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return uint32(len(c.ConnManager))
}

//removeall
func (c *ConnectionManager) Clear() {

	c.lock.Lock()


	for connID, conn := range c.ConnManager {

		conn.Stop()

		delete(c.ConnManager, connID)

	}
	c.lock.Unlock()
	//fmt.Println("delete all conn from manager,len = ", c.GetLen())

}
