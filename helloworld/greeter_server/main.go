package main

import (
	"context"
	"fmt"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_server/interceptor"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

const (
	port = ":50052"
)

type server struct{}

// SayHelloメソッドを実装
func (s *server) SayHello(ctx context.Context, helloRequest *pb.HelloRequest) (*pb.HelloReply, error) {
	//　time.Sleep(time.Millisecond * 3000)
	log.Printf("Received: %v", helloRequest.Name)
	return &pb.HelloReply{Message: "Hello " + helloRequest.Name}, nil
}

// one to many
// goroutine
func (s *server) SayHelloServerSideStreaming (req *pb.HelloRequest, stream pb.Greeter_SayHelloServerSideStreamingServer)  error {
	coundown := 10
	for coundown > 0 {
		result := fmt.Sprintf("Hello  %s, %d.", req.Name, coundown)
		res := &pb.HelloReply{Message: result}

		// clientにresを送っている。（userの数だけ => clinetのRecvメソッドで受信される）
		coundown--
		err := stream.Send(res)
		if err != nil {
			fmt.Printf("during streaming send error: %v\n", err)
			return err
		}
	}
	return nil
}

func (s *server) SayHelloClientSideStreaming(stream pb.Greeter_SayHelloClientSideStreamingServer) error{
	result := "Hello!"
	log.Print("-------------- weak---")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			r := &pb.HelloReply{
				Message: result,
			}
			return stream.SendAndClose(r)
		}
		if err != nil {
			fmt.Printf("during streaming receive error: %v\n", err)
			return err
		}
		result += fmt.Sprintf("Hello! %s", req.Name)
		log.Printf("-------------- weak-- %s", result)
	}
}

func (s *server) SayHelloStreaming(stream pb.Greeter_SayHelloStreamingServer) error {
	var count int
	for {
		count++
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("during streaming receive error: %v\n", err)
			return err
		}
		err = stream.Send(&pb.HelloReply{
			Message: fmt.Sprintf("Hello! req is %s, count is %v", req.Name, count),
		})
		if err != nil {
			fmt.Printf("during streaming send error: %v\n", err)
			return err
		}
	}
}

func main() {
	// リッスン処理
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// サーバ起動(interceptorも一緒に設定しています。)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor),
		grpc.StreamInterceptor(interceptor.ServerStreamInterceptor))
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

