package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "grpc-client/proto"
	"io"
	"io/ioutil"
	"log"
)

// normalClient 普通模式客户端
func normalClient(client pb.StreamGetUsersServiceClient) error {
	// 构造请求数据
	users :=[]*pb.UserInfo{
		{UserId: 1},
		{UserId: 2},
		{UserId: 3},
	}
	req := &pb.StreamUserRequest{Users: users}
	// 发送请求
	resp, err := client.GetUserScore(context.Background(), req)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("response:", resp)
	return nil
}

// serverStream 服务端流请求客户端
func serverStream(client pb.StreamGetUsersServiceClient) error {
	// 构造请求数据
	var i int32
	req := pb.StreamUserRequest{}
	for i = 0; i < 11; i++ {
		req.Users = append(req.Users, &pb.UserInfo{UserId: i+1})
	}

	// 发送请求
	stream, err := client.GetUsersScoreByServer(context.Background(), &req)
	if err != nil {
		log.Fatalln(err)
	}

	// 循环接收数据
	for {
		// Recv 方法阻塞等待服务端返回的数据
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(response)
	}
	return nil
}

// clientStream 客户端流请求客户端
func clientStream(client pb.StreamGetUsersServiceClient) error {
	stream, err := client.GetUserScoreByClientStream(context.Background())
	if err != nil {
		return err
	}

	// 构造请求数据
	var i int32
	var uid int32 = 1
	req := pb.StreamUserRequest{}

	for x := 1; x < 4; x++ {
		for i = 0; i < 5; i++ { // 假设这是一个耗时的过程.
			req.Users = append(req.Users, &pb.UserInfo{UserId: uid})
			uid++
		}
		// 每次请求5个用户的信息,共发三次请求
		err := stream.Send(&req)
		if err != nil {
			log.Println(err)
		}
	}
	// 所有请求发送完成之后,等待服务端返回响应
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	fmt.Println("response:", resp)
	return nil
}

// bidirectionalStream 双向流请求客户端
func bidirectionalStream(client pb.StreamGetUsersServiceClient) error {
	stream, err := client.GetUserScoreByTWF(context.Background())
	if err != nil {
		return err
	}
	// 构造请求数据
	var i int32
	var uid int32 = 1
	req := pb.StreamUserRequest{}

	for x := 1; x < 4; x++ {
		for i = 0; i < 5; i++ { // 假设这是一个耗时的过程.
			req.Users = append(req.Users, &pb.UserInfo{UserId: uid})
			uid++
		}
		// 每次请求5个用户的信息,共发三次请求
		err := stream.Send(&req)
		if err != nil {
			log.Println(err)
		}
		// 重置请求,准备下次请求数据
		req = pb.StreamUserRequest{}

		// 等待响应
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		fmt.Println("response:", resp)
	}

	stream.CloseSend()
	return nil
}

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

	// 构造客户端
	client := pb.NewStreamGetUsersServiceClient(conn)

	/*
	// 普通客户端请求
	err = normalClient(client)
		if err != nil {
			log.Println(err)
	}
	*/


	/*
	// 服务端流请求测试
	err = serverStream(client)
	if err != nil {
		log.Println(err)
	}
	 */

	/*
	// 客户端流请求测试
	err = clientStream(client)
	if err != nil {
		log.Println(err)
	}
	*/

	// 双向流客户端测试
	err = bidirectionalStream(client)
	if err != nil {
		log.Println(err)
	}
}
