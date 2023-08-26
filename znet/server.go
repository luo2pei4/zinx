package znet

import (
	"fmt"
	"luopei/zinx/ziface"
	"net"
	"time"
)

type Server struct {
	Name  string
	IPVer string
	IP    string
	Port  int
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:  name,
		IPVer: "tcp4",
		IP:    "0.0.0.0",
		Port:  6666,
	}
}

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

		// 接收tcp数据
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Printf("listen tcp failed, %s\n", err.Error())
				continue
			}
			go func() {
				for {
					buf := make([]byte, 512)
					n, err := conn.Read(buf)
					if err != nil {
						fmt.Printf("read from conn failed, %s\n", err.Error())
						continue
					}
					if _, err := conn.Write(buf[:n]); err != nil {
						fmt.Printf("write back to conn failed, %s\n", err.Error())
						continue
					}
				}
			}()
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
