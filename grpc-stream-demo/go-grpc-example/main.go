package main

import (
	"crypto/tls"
	"crypto/x509"
	pb "go-grpc-example/proto"
	service "go-grpc-example/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	// 加载证书配置
	// credentials.NewServerTLSFromFile: 根据服务端输入的证书文件和密钥构造 TLS 凭证
	cert, err := tls.LoadX509KeyPair("cert/server.pem", "cert/server-key.pem")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair error: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("cert/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err: %v\n", err)
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	// 构造一个TLS gRPC服务对象
	// grpc.Creds(): 返回一个 ServerOption,用于设置服务器连接的凭据.
	// 用于 grpc.NewServer(opt ...ServerOption) 为 gRPC Server 设置连接选项.
	rpcServer := grpc.NewServer(grpc.Creds(c))

	// 通过gRPC插件生成的RegisterStreamGetUsersServiceServer函数注册我们实现的StreamUsersService服务
	pb.RegisterStreamGetUsersServiceServer(rpcServer, new(service.StreamUsersService))

	// 监听端口
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalln("net.Listen error:", err)
	}

	// 在监听的端口上启动服务
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatalln("Start server failed, error:", err)
	}
}
