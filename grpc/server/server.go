package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "test/grpc/grpc"
)

type server struct {
}

func (s *server) AddOperationLog(ctx context.Context, in *pb.AddOperationLogRequest) (*pb.AddOperationLogReply, error) {
	log.Printf("in param: %+v", in)
	return &pb.AddOperationLogReply{
		Code: 1,
	}, nil
}

func main() {
	service := server{}
	conn, err := net.Listen("tcp", ":8801")
	if err != nil {
		log.Fatalf("listen port failed, err: %v", err)
	}
	log.Printf("init grpc")
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("close connect failed, err: %v", err)
		}
	}()
	server := grpc.NewServer()
	pb.RegisterAdminServer(server, &service)
	log.Printf("rpc server start...")
	if err := server.Serve(conn); err != nil {
		log.Printf("rpc server start failed, err: %v", err)
		return
	}
}
