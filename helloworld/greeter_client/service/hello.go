package service

import (
	"context"
	"fmt"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/dto"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

const (
	port = ":50051"
)

func Hello(c pb.GreeterClient, word string) (*dto.HelloReply, bool) {

	// contextの準備
	// clientDeadline := time.Now().Add(time.Duration(2 * time.Second))

	// timeoutの場合
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, cancel := context.WithTimeout(context.Background(), time.Second)

	// WithTimeoutの場合：WithDeadlineをラッパーしている。（引数に時間を指定しないと無限に待ち続けるようになる。）
	// WithDeadlineの場合：リクエストのライフサイクル全体にタイムアウトを直接指定する
	// ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	// ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	// ****** Metadata : Creation *****
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"kn", "vn",
	)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)
	ctxA := metadata.AppendToOutgoingContext(mdCtx, "k1", "v1", "k1", "v2", "k2", "v3")

	// RPC using the context with new metadata.
	var header, trailer metadata.MD

	// SayHelloメソッドの呼び出し(metadataなし)
	// r, err := c.SayHello(ctx, &pb.HelloRequest{Name: word}, grpc.Header(&header), grpc.Trailer(&trailer))

	// SayHelloメソッドの呼び出し(metadataあり)
	// grpc.UseCompressor(gzip.Name) :
	// grpc.Header(&header) :
	// grpc.Trailer(&trailer) :
	r, err := c.SayHello(ctxA, &pb.HelloRequest{Name: word}, grpc.UseCompressor(gzip.Name), grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Println("could not greet: %v", err)

		errorCode := status.Code(err)
		if errorCode == codes.InvalidArgument {
			log.Printf("Invalid Argument Error : %s", errorCode)
			errorStatus := status.Convert(err)
			for _, d := range errorStatus.Details() {
				switch info := d.(type) {
				case *epb.BadRequest_FieldViolation:
					log.Printf("Request Field Invalid: %s", info)
				default:
					log.Printf("Unexpected error type: %s", info)
				}
			}
		} else {
			log.Printf("Unhandled error : %s ", errorCode)
		}
		return nil, false
	}

	// Serverからのmetadataを確認する
	if t, ok := header["timestamp"]; ok {
		log.Printf("timestamp from header:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Printf("+++++++++++++++++++++++++++++++++++++++++++++")
		log.Printf("timestamp expected but doesn't exist in header")
		log.Printf("+++++++++++++++++++++++++++++++++++++++++++++")
	}
	if l, ok := header["location"]; ok {
		log.Printf("location from header:\n")
		for i, e := range l {
			fmt.Printf(" %d. %s\n", i, e)
		}
	} else {
		log.Printf("+++++++++++++++++++++++++++++++++++++++++++++")
		log.Printf("location expected but doesn't exist in header")
		log.Printf("+++++++++++++++++++++++++++++++++++++++++++++")
	}


	hr := &dto.HelloReply{Message:r.Message}
	log.Printf("Greeting: %s", r.Message)
	return hr, true
}
