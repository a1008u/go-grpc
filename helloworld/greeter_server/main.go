package main

import (
	"context"
	"crypto/tls"
	"fmt"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_server/interceptor"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"time"
)

const (
	port = ":50052"
)

type server struct{}

// SayHelloメソッドを実装
func (s *server) SayHello(ctx context.Context, helloRequest *pb.HelloRequest) (*pb.HelloReply, error) {
	//　time.Sleep(time.Millisecond * 3000)

	// Client側からのメタデータ読み込み
	md, metadataAvailable := metadata.FromIncomingContext(ctx)
	if !metadataAvailable {
		return nil, status.Errorf(codes.DataLoss, "UnaryEcho: failed to get metadata")
	}
	if t, ok := md["timestamp"]; ok {
		fmt.Printf("timestamp from metadata:\n", md)
		for i, e := range t {
			fmt.Printf("====> Metadata %d. %s\n", i, e)
		}
	}

	//　gRPCの引数の確認
	if helloRequest.Name == "error" {
		log.Printf("Order ID is invalid! -> Received Name %s", helloRequest.Name)

		errorStatus := status.New(codes.InvalidArgument, "Invalid information received")
		ds, err := errorStatus.WithDetails(
			&epb.BadRequest_FieldViolation{
				Field:"Name",
				Description: fmt.Sprintf("Order Name received is not valid %s : ", helloRequest.Name),
			},
		)
		if err != nil {
			return nil, errorStatus.Err()
		}

		// Server側からのmetadataの返却
		header := metadata.New(map[string]string{"location": "JP", "timestamp": time.Now().Format(time.StampNano)})
		grpc.SendHeader(ctx, header)
		return nil, ds.Err()
	}

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

