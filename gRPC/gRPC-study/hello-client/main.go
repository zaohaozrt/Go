package main

import (
	"context"
	"fmt"
	pb "gRPC-study/hello-client/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientTokenAuth struct {
}

// token
func (c ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appId":  "root",
		"appKey": "123456",
	}, nil
}

// 不使用安全认证
func (c ClientTokenAuth) RequireTransportSecurity() bool {
	return false
}
func main() {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())) //证书验证（无）
	opts = append(opts, grpc.WithPerRPCCredentials(new(ClientTokenAuth)))         //token验证

	//连接到server端，此处禁用安全传输，没有加密和验证
	//基于tcp连接建立RPC连接
	conn, err := grpc.Dial("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	defer conn.Close()

	//通过RPC连接找到对应服务
	client := pb.NewSayHelloClient(conn)
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequet{RequestName: "wsz"})
	fmt.Println(resp.GetResponseMgs())
}
