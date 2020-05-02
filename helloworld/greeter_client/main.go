package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/a1008u/go-grpc/helloworld"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/domain"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/dto"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/interceptor"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/service"
	"github.com/a1008u/go-grpc/helloworld/greeter_client/util"
	"github.com/a1008u/go-grpc/helloworld/tracer"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"strings"
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


// Unary
func grpcClient(w http.ResponseWriter, r *http.Request) {
	// gRPCコネクションの作成
	address := util.GetGrcpAddress()

	// tracer
	jaegertracer, closer, err := tracer.NewTracer("product_mgt")
	if err != nil {
		log.Println("eeeerrrrorr")
	}
	defer closer.Close()

	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		//grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				interceptor.Uac(),
				grpcMetrics.UnaryClientInterceptor(),
				grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(jaegertracer)),
			),
		),
		//grpc.WithStreamInterceptor(grpcopentracing.StreamClientInterceptor(grpcopentracing.WithTracer(jaegertracer)))
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				interceptor.Sci(),
				grpcMetrics.StreamClientInterceptor(),
				//grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(jaegertracer)),
			),
		),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	var word string
	if strings.Index(r.URL.Path, "error") != -1 {
		word = "error"
	} else {
		word = "world"
	}

	hr, result := service.Hello(c, word)

	if result {
		grpchandler(hr, w, r)
	} else {
		handler(w,r)
	}
}

func grpcClientStreamServer(w http.ResponseWriter, r *http.Request) {
	// gRPCコネクションの作成
	address := util.GetGrcpAddress()

	// tracer
	jaegertracer, closer, err := tracer.NewTracer("product_mgt")
	if err != nil {
		log.Println("eeeerrrrorr")
	}
	defer closer.Close()

	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),

		//grpc.WithStreamInterceptor(grpcopentracing.StreamClientInterceptor(grpcopentracing.WithTracer(jaegertracer)))
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				interceptor.Sci(),
				grpcMetrics.StreamClientInterceptor(),
				grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(jaegertracer)),
			),
		),
	)
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
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptor.ClientStreamInterceptor))
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
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptor.ClientStreamInterceptor))
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

var (
	grpcMetrics *grpc_prometheus.ClientMetrics
)

func main() {

	// OpenCensusの設定
	// Register stats and trace exporters to export
	// the collected data.
	view.RegisterExporter(&exporter.PrintExporter{})
	// Register the view to collect gRPC client stats.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	// prometheusの設定
	// Create a metrics registry.
	reg := prometheus.NewRegistry()
	// Create some standard client metrics.
	grpcMetrics = grpc_prometheus.NewClientMetrics()
	// Register client metrics to registry.
	reg.MustRegister(grpcMetrics)
	// Create a HTTP server for prometheus.
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: ":9094"}
	// Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	//http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.HandleFunc("/", handler)
	http.HandleFunc("/grpc", grpcClient)
	http.HandleFunc("/error", grpcClient)
	http.HandleFunc("/grpc2", grpcClientStreamServer)
	http.HandleFunc("/grpc3", grpcSideStreaming)
	http.HandleFunc("/grpc4", grpcStreaming)

	// client側のserver実行
	if err := http.ListenAndServe(":50051", nil); err != nil {
		log.Fatal("Unable to start a http server.")
	}
	//go func() {
	//	if err := http.ListenAndServe(":50054", nil); err != nil {
	//		log.Fatal("Unable to start a http server.")
	//	}
	//}()
}
