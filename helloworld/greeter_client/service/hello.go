package service

import (
	"context"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/dto"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/util"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const (
	defaultName = "world"
	port = ":50051"
)

func Hello() (*dto.HelloReply, bool) {
	// gRPCコネクションの作成
	address := util.GetGrcpAddress()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, false
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 引数の準備
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	// contextの準備
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SayHelloメソッドの呼び出し
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return nil, false
	}

	hr := &dto.HelloReply{Message:r.Message}

	log.Printf("Greeting: %s", r.Message)
	return hr, true
}
