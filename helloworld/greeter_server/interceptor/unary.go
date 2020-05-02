package interceptor

import (
	"context"
	"github.com/a1008u/go-grpc/helloworld/tracer"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"log"
)

// Server - Unary Interceptor ---------------------------------------
func UnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	// 前処理：preprocessing
	log.Printf("======= [Server Interceptor]  FullMethod is %s ===== Serer is %s", info.FullMethod, info.Server)

	// Tracer
	// initialize jaegertracer
	jaegertracer, closer, err := tracer.NewTracer("product_mgt")
	if err != nil {
		log.Fatalln(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(jaegertracer)
	grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(jaegertracer))

	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)

	// 後処理：postprocessing
	log.Printf(" Post Proc Message : %s", m)
	return m, err
}
