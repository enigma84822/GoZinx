package znet

import (
	"GoZinx/ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //连接的管理集合
	connLock    sync.RWMutex                  // 读写锁
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 讲conn加入到connManager里
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Printf("connID=%d add to ConnManager successfully! conn num=%d\n", conn.GetConnID(), connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 讲conn加入到connManager里
	delete(connMgr.connections, conn.GetConnID())
	fmt.Printf("connID=%d remove from ConnManager successfully! conn num=%d\n", conn.GetConnID(), connMgr.Len())

}

func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found!\n")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		conn.Stop()

		delete(connMgr.connections, connID)
	}
	fmt.Printf("Clear All Connections succ! conn num=%d\n", connMgr.Len())
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}
