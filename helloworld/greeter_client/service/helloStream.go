package service

import (
	"context"
	"fmt"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/dto"
	"io"
)

func CallStreaming(c pb.GreeterClient) ([]*dto.HelloReply, bool) {
	stream, err := c.SayHelloStreaming(context.Background())
	if err != nil {
		fmt.Printf("greet bidirectional streame create error: %v\n", err)
		return nil, false
	}

	var helloReplys []*dto.HelloReply
	waitc := make(chan struct{})

	// send処理
	go func() {
		countdown := 3
		for countdown > 0{
			countdown--
			err := stream.Send(&pb.HelloRequest{Name: fmt.Sprintf("test %v", countdown)})
			if err != nil {
				fmt.Printf("during streaming send error: %v\n", err)
				// error処理の方法は??
			}
		}
		err := stream.CloseSend()
		if err != nil {
			fmt.Printf("close send error: %v\n", err)
			// error処理の方法は??
		}
	}()

	// Recv処理（全てをRecvし終わると完了となる）
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("during streaming receive error: %v\n", err)
				// error処理の方法は??
			}
			fmt.Printf("greet streaming response: %v\n", res.Message)
			helloReplys = append(helloReplys, &dto.HelloReply{Message: fmt.Sprintf("greet streaming response: %v", res.Message)})
		}
		close(waitc)
	}()

	<-waitc
	return helloReplys, true
}
