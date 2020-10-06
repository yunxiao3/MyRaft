package main

import (
    "net"
    "log"
    "fmt"
    "google.golang.org/grpc"
    "golang.org/x/net/context"
)

const (
    port = "15000"
)

type server struct {
}
// 实现服务端的方法
func (s *server) RegisterServer(ctx context.Context, in *RouterInfo) (*ResultToApp, error) {
    fmt.Print(*in)
    name := "test"
    pas := "123"
    return &ResultToApp{UserName: name, Password: pas}, nil
}

func ListenFromJava() {
    fmt.Println("接收到消息")
    //监听接口
    serverLis, err := net.Listen("tcp", "127.0.0.1:"+port)

    if err != nil {
        log.Fatal("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    //注册服务
    RegisterMessageChannelServer(s, &server{})
    if err := s.Serve(serverLis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

