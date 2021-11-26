package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	pb "test/grpc/grpc"
)

func main() {
	Address := "127.0.0.1:8801"
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	client := pb.NewAdminClient(conn)                                                         //  自动生成方法, 因为proto文件中service的名字是 admin
	result, err := client.AddOperationLog(context.Background(), &pb.AddOperationLogRequest{}) // 调用grpc方法, 对服务端进行通讯
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("I see code: %v \n", result.Code)
}
