package main

import (
	"context"
	"encoding/json"
	pb "grpc_demo/iAutoApi"
	"io"
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

func initilize() (*grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Fatalf("auth failed, err : %v", err)
	}

	return grpc.Dial(address, grpc.WithTransportCredentials(creds))
}

func main() {

	conn, err := initilize()
	if err != nil {
		log.Panicf("dial failed, %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := pb.NewIAutoOAClient(conn)

	// case : invoke GetEmployeeInfo
	ret, err := c.GetEmployeeInfo(ctx, &pb.Requestor{Type: "Hello world"})
	if err != nil {
		log.Fatalf("grpc invoke interface failed, %v", err)
	}

	b, err := json.Marshal(&ret)
	if err != nil {
		log.Fatalf("Invoke fatal error, %s", err)
	}
	log.Printf("Ret: %s", string(b[:]))

	// case : invoke EchoMessage
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cStream, err := c.EchoMessage(ctx)
	if err != nil {
		log.Fatalf("Invoke fatal error, %s", err)
	}

	waitc := make(chan bool, 1)
	go func() { // send stream closure package
		var message = map[int32]string{
			1: "First Message",
			2: "Second Message",
			3: "Third Message",
			4: "Fourth Message",
		}

		for id, mesg := range message {
			log.Printf("send :%id, %s", id, mesg)
			cStream.Send(&pb.SRequestor{Id: id, Mesg: mesg})
		}
	}()

	go func() { // recv stream closure package
		for {
			resp, err := cStream.Recv()
			if err == io.EOF {
				log.Printf("read eof")
				break
			}
			if err != nil {
				log.Fatalf("Invalid recv")
				break
			}

			log.Printf("Echo message: %s", resp.GetMesg())
			if resp.GetId() == 4 {
				waitc <- true
				break
			}
		}
	}()

	<-waitc
	log.Printf("Exit program")
	log.Printf("%s", string(b[:]))
}
