package service

import (
	"context"
	"fmt"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/dto"
	"io"
)

/*
１つのrequestに対してserverが複数のresponseをstreamingとして返す形式
*/
func CallServerStreaming(c pb.GreeterClient) ([]*dto.HelloReply, bool) {

	var helloReplys []*dto.HelloReply

	req := &pb.HelloRequest{
		Name: "clinet to server",
	}

	// streamはgreetServiceGreetServerSideStreamingClient（protoで作成）です。
	stream, err := c.SayHelloServerSideStreaming(context.Background(), req)
	if err != nil {
		fmt.Printf("greet server side stream create error: %v\n", err)
		return nil, false
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("during streaming receive error: %v\n", err)
			return nil, false
		}
		helloReplys = append(helloReplys, &dto.HelloReply{Message:res.Message})
		fmt.Printf("greet server side streaming response: %v\n", res.Message)
	}
	return helloReplys, true
}
