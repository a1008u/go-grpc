package service

import (
	"context"
	"fmt"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/dto"
	"os"
)

func CallClientStreaming(c pb.GreeterClient) (*dto.HelloReply, bool) {
	stream, err := c.SayHelloClientSideStreaming(context.Background())
	if err != nil {
		fmt.Printf("greet client side stream create error: %v\n", err)
		os.Exit(1)
	}

	countdown := 10
	for countdown > 0 {
		countdown--
		err := stream.Send(&pb.HelloRequest{Name: fmt.Sprintf("test %v", countdown)})
		if err != nil {
			fmt.Printf("during streaming send error: %v\n", err)
			return nil, false
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("close and receive response error: %v\n", err)
		return nil, false
	}
	fmt.Printf("greet client side streaming response: %v\n", res.Message)
	return &dto.HelloReply{Message:res.Message}, true
}
