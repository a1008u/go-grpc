package service

import (
	"context"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/dto"
	"log"
	"os"
	"time"
)

const (
	defaultName = "world"
	port = ":50051"
)

func Hello(c pb.GreeterClient) (*dto.HelloReply, bool) {
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
