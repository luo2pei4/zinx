package main

import "luopei/zinx/znet"

func main() {
	s := znet.NewServer("luopei")
	s.Serve()
}
