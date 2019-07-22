package main

import (
	"context"
	pb "grpc_demo/iAutoApi"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = "iauto:50001"
	key  = "/home/liuhu/go/bin/pki/server/server.key"
	cert = "/home/liuhu/go/bin/pki/server/server.pem"
)

type server struct{}

func (s *server) GetEmployeeInfo(ctx context.Context, in *pb.Requestor) (*pb.EmployeeInfo, error) {
	out := new(pb.EmployeeInfo)
	out.Name = "LiuHu"
	out.Department = "Voice lib"
	out.Number = 3302

	return out, nil
}

func (s *server) EchoMessage(stream pb.IAutoOA_EchoMessageServer) error {
	log.Printf("recv client message")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("read eof")
			return nil
		}

		if err != nil {
			return err
		}

		log.Printf("Recv Id: %d, Mesg:%s\r\n", req.GetId(), req.GetMesg())

		if err = stream.Send(&pb.SResponsor{Id: req.GetId(), Mesg: req.GetMesg()}); err != nil {
			log.Fatalf("Invoker function send failed")
			break
		}
	}
	return nil
}

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
