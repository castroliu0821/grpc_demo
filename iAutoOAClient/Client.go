package main

import (
	"context"
	"encoding/json"
	pb "grpc_demo/iAutoApi"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "iauto:50001"
	name    = "word"
	cert    = "/home/liuhu/go/bin/pki/server/server.pem"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Fatalf("auth failed, err : %v", err)
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Dial server failed, %v", err)
	}

	defer conn.Close()

	c := pb.NewIAutoOAClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := pb.Requestor{Type: "Hello world"}
	ret, err := c.GetEmployeeInfo(ctx, &req)
	if err != nil {
		log.Fatalf("grpc invoke interface failed, %v", err)
	}

	b, err := json.Marshal(&ret)
	if err != nil {
		log.Fatalf("to conv json failed : %v", err)
	}

	log.Printf("%s", string(b[:]))
}
