package main

import (
	"context"
	pb "grpc_demo/iAutoApi"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct{}

func (s *server) GetEmployeeInfo(ctx context.Context, in *pb.Requestor) (*pb.EmployeeInfo, error) {
	out := new(pb.EmployeeInfo)
	out.Name = "LiuHu"
	out.Department = "Voice lib"
	out.Number = 3302

	return out, nil
}

const (
	port = "iauto:50001"
	key  = "/home/liuhu/go/bin/pki/server/server.key"
	cert = "/home/liuhu/go/bin/pki/server/server.pem"
)

func main() {
	log.Printf("start server...")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("create listen failed, %v", err)
	}

	trans, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatalf("create auth failed, err:%v", err)
	}

	srv := grpc.NewServer(grpc.Creds(trans))
	pb.RegisterIAutoOAServer(srv, &server{})
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("bind listen failed, err:%v", err)
	}
}
