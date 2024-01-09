package main

import (
	"context"
	"errors"
	"fmt"
	pb "gRPC-study/hello-client/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequet) (*pb.HelloResponse, error) {
	//获取元数据,检验token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("未传输token")
	}
	var appId string
	var appKey string
	if v, ok := md["appid"]; ok {
		appId = v[0]
	}
	if v, ok := md["appkey"]; ok {
		appKey = v[0]
	}
	if appId != "root" || appKey != "123456" {
		return nil, errors.New("token 不正确")
	}

	return &pb.HelloResponse{ResponseMgs: "hello" + req.RequestName}, nil
}

func main() {
	//开启端口
	listen, _ := net.Listen("tcp", ":9090")
	//创建gRPC
	grpcServer := grpc.NewServer()
	//在grpc服务端中去注册自己编写的服务
	pb.RegisterSayHelloServer(grpcServer, &server{})
	//启动服务,grpc基于tcp协议，protobuf用来进行数据序列化反序列化
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("启动服务失败")
		return
	}
}
