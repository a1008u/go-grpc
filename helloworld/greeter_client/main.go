package greeter_client

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// gRPCコネクションの作成
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
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
	}
	log.Printf("Greeting: %s", r.Message)
}
