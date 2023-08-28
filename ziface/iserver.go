package ziface

import "net"

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(router IRouter)
}

// 定义连接接口
type IConnection interface {
	Start()                         // 启动连接
	Stop()                          // 关闭连接
	GetTCPConnection() *net.TCPConn // 获取TCP连接
}

// 定义一个统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
