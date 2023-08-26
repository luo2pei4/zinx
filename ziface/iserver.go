package ziface

import "net"

type IServer interface {
	Start()
	Stop()
	Serve()
}

// 定义连接接口
type IConnection interface {
	Start()            // 启动连接
	Stop()             // 关闭连接
	GetConnID() uint32 // 获取当前连接ID
}

// 定义一个统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
