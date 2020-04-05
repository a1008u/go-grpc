package main

import (
	"context"
	pb "github.com/a1008u/go-grpc/helloworld"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50052"
)

type server struct{}

// SayHelloメソッドを実装
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// リッスン処理
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// サーバ起動
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

