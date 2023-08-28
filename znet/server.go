package znet

import (
	"errors"
	"fmt"
	"io"
	"luopei/zinx/ziface"
	"net"
	"time"
)

type Server struct {
	Name   string
	IPVer  string
	IP     string
	Port   int
	Router ziface.IRouter
}

type Connection struct {
	Conn     *net.TCPConn // 当前连接的socket tcp套接字
	ConnID   uint32       // 当前连接的ID
	isClosed bool         // 当前连接关闭状态
	// handleAPI   ziface.HandFunc // 该连接的处理方法API
	Router      ziface.IRouter // 通过路由接口处理连接请求
	ExitBufChan chan bool      // 通知连接推出的chan
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:   name,
		IPVer:  "tcp4",
		IP:     "0.0.0.0",
		Port:   6666,
		Router: nil,
	}
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		// handleAPI:   callback_api,
		Router:      router,
		ExitBufChan: make(chan bool, 1),
	}
}

// func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
// 	fmt.Println("[Conn Handle] CallBackToClient...")
// 	if _, err := conn.Write(data[:cnt]); err != nil {
// 		fmt.Println("write back buf error,", err.Error())
// 		return errors.New("CallBackToClient error")
// 	}
// 	return nil
// }

func (s *Server) Start() {
	fmt.Printf("[START] server listenner at IP: %s, Port: %d, is starting\n", s.IP, s.Port)
	go func() {
		// 获取TCP地址
		addr, err := net.ResolveTCPAddr(s.IPVer, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Printf("resolve ip addr failed, %s\n", err.Error())
			return
		}
		// 创建监听
		listenner, err := net.ListenTCP(s.IPVer, addr)
		if err != nil {
			fmt.Printf("listen tcp failed, %s\n", err.Error())
			return
		}

		fmt.Printf("start zinx server %s successed, now listenning...", s.Name)

		cID := uint32(0)

		// 接收tcp数据
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Printf("listen tcp failed, %s\n", err.Error())
				continue
			}

			// 用获取的连接新建一个用来处理链接数据的实例
			dealConn := NewConnection(conn, cID, s.Router)
			cID++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Printf("[STOP] zinx server, %s", s.Name)
}

func (s *Server) Serve() {
	s.Start()

	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("add router successed")
}

// 读取请求数据的方法，以goroutine形式运行
func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), "conn reader exit.")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Println("receive buf failed,", err.Error())
			}
			c.ExitBufChan <- true
			return
		}
		// if err := c.handleAPI(c.Conn, buf, n); err != nil {
		// 	fmt.Println("connID", c.ConnID, "handle is error.", err.Error())
		// 	c.ExitBufChan <- true
		// 	return
		// }
		req := Request{
			conn: c,
			data: buf,
		}
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	<-c.ExitBufChan
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	c.ExitBufChan <- true
	close(c.ExitBufChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
