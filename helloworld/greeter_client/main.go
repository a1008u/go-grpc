package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/domain"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/dto"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/service"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/util"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

const (
	defaultContentType = "application/json; charset=utf-8"
)



func handler(w http.ResponseWriter, r *http.Request) {
	user := domain.User{
		Name:    "taro",
		Age:     18,
	}

	w.Header().Set("content-type", defaultContentType)
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	if err := e.Encode(user); err != nil {
		fmt.Fprintf(w, "json encode err: %v\n", err)
		os.Exit(1)
	}
}

func grpchandler(hr *dto.HelloReply, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", defaultContentType)
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	if err := e.Encode(hr); err != nil {
		fmt.Fprintf(w, "json encode err: %v\n", err)
		os.Exit(1)
	}
}

func grpchandlers(hr []*dto.HelloReply, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", defaultContentType)
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	if err := e.Encode(hr); err != nil {
		fmt.Fprintf(w, "json encode err: %v\n", err)
		os.Exit(1)
	}
}

func grpcClient(w http.ResponseWriter, r *http.Request) {
	// gRPCコネクションの作成
	address := util.GetGrcpAddress()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	hr, result := service.Hello(c)
	if result {
		grpchandler(hr, w, r)
	} else {
		handler(w,r)
	}
}

func grpcClientStreamServer(w http.ResponseWriter, r *http.Request) {
	// gRPCコネクションの作成
	address := util.GetGrcpAddress()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	hr, result := service.CallServerStreaming(c)
	if result {
		grpchandlers(hr, w, r)
	} else {
		handler(w,r)
	}
}

func grpcSideStreaming(w http.ResponseWriter, r *http.Request) {
	address := util.GetGrcpAddress()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	hr, result := service.CallClientStreaming(c)
	if result {
		grpchandler(hr, w, r)
	} else {
		handler(w,r)
	}
}

func grpcStreaming(w http.ResponseWriter, r *http.Request) {
	address := util.GetGrcpAddress()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	hr, result := service.CallStreaming(c)
	if result {
		grpchandlers(hr, w, r)
	} else {
		handler(w,r)
	}

}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/grpc", grpcClient)
	http.HandleFunc("/grpc2", grpcClientStreamServer)
	http.HandleFunc("/grpc3", grpcSideStreaming)
	http.HandleFunc("/grpc4", grpcStreaming)
	http.ListenAndServe(":50051", nil)
}
