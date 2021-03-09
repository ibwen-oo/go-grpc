package server

import (
	"context"
	pb "go-grpc-example/proto"
	"io"
	"log"
	"time"
)

// 定义服务结构体
type StreamUsersService struct{}

// 实现rpc方法 普通模式
func (s *StreamUsersService) GetUserScore(context.Context, *pb.StreamUserRequest) (*pb.StreamUserResponse, error) {
	users := []*pb.UserInfo{
		{UserId: 1, UserScore: 100},
		{UserId: 2, UserScore: 200},
		{UserId: 3, UserScore: 300},
	}
	response := &pb.StreamUserResponse{Users: users}
	return response, nil
}

// 实现rpc方法 服务端流
func (s *StreamUsersService) GetUsersScoreByServer(r *pb.StreamUserRequest, stream pb.StreamGetUsersService_GetUsersScoreByServerServer) error {
	// 构造数据(用户积分)
	var score int32 = 100
	users := make([]*pb.UserInfo, 0)
	for i, user := range r.Users {
		user.UserScore = score
		score++
		users = append(users, user)
		// 每查到两条用户积分就返回给客户端
		if (i+1) % 2 == 0 && i > 0 {
			// send 方法发送数据到客户端
			err := stream.Send(&pb.StreamUserResponse{Users: users})
			if err != nil {
				log.Fatalln("steam.Send error:", err)
			}
			// 清空users切片,重新构造数据
			users = (users)[0:0]
			time.Sleep(time.Second*1)
		}
	}
	// 如果还有遗留数据,返回给客户端.
	if len(users) > 0 {
		err := stream.Send(&pb.StreamUserResponse{Users: users})
		if err != nil {
			log.Fatalln("steam.Send error:", err)
		}
	}

	return nil
}

// 实现rpc方法 客户端流
func (s *StreamUsersService) GetUserScoreByClientStream(stream pb.StreamGetUsersService_GetUserScoreByClientStreamServer) error {
	var score int32 = 100
	users := make([]*pb.UserInfo, 0)

	for {
		r, err := stream.Recv()
		if err == io.EOF {  // 接收完了
			return stream.SendAndClose(&pb.StreamUserResponse{Users: users})
		}
		if err != nil {
			return err
		}

		for _, user := range r.Users {
			user.UserScore = score
			score++
			users = append(users, user)
		}
	}
}

// 实现rpc方法 双向流
func (s *StreamUsersService) GetUserScoreByTWF(stream pb.StreamGetUsersService_GetUserScoreByTWFServer) error {
	var score int32 = 100
	users := make([]*pb.UserInfo, 0)

	for {
		r, err := stream.Recv()
		if err == io.EOF {  // 接收完了
			return nil
		}
		if err != nil {
			return err
		}

		for _, user := range r.Users {
			user.UserScore = score
			score++
			users = append(users, user)
		}
		// 返回响应
		err = stream.Send(&pb.StreamUserResponse{Users: users})
		if err != nil {
			log.Println("stream send data failed,", err)
		}
		// 每次发送之后清空测试数据,生产中不用
		users = users[0:0]
	}
}