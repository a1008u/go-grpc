package main

import (
	"encoding/json"
	"fmt"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/domain"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/dto"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/service"
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

func grpcClient(w http.ResponseWriter, r *http.Request) {
	hr, result := service.Hello()
	if result {
		grpchandler(hr, w, r)
	} else {
		handler(w,r)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/grpc", grpcClient)
	http.ListenAndServe(":50051", nil)
}
