package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "grpc-client/proto"
	"io/ioutil"
	"log"
)

func main() {
	// 加载证书配置
	cert, err := tls.LoadX509KeyPair("cert/client.pem", "cert/client-key.pem")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("cert/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		RootCAs:      certPool,
	})

	// 与gRPC服务端建立连接
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	// 通过gRPC插件生成的NewProdServiceClient函数基于建立的连接构造ProdServiceClient对象
	// 返回的ProdServiceClient对象其实是一个ProdServiceClient接口对象,
	// 通过接口定义的方法就可以调用服务端对应的gRPC服务提供的方法.
	prodServiceClient := pb.NewProdServiceClient(conn)

	// ProdServiceClient对象调用接口定义的方法.
	response, err := prodServiceClient.GetProdStock(context.Background(), &pb.ProdRequest{ProdId: 10})
	if err != nil {
		log.Fatalln("prodServiceClient.GetProdStock failed, error:", err)
	}

	fmt.Println(response.ProdStock)
}
