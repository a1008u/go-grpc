package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

// Server - Unary Interceptor ---------------------------------------
func UnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	// 前処理：preprocessing
	beforeMessage := fmt.Sprintf("======= [Server Interceptor]  FullMethod is %s ===== Serer is %s", info.FullMethod, info.Server)
	before(beforeMessage)

	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)

	// 後処理：postprocessing
	afterMessage := fmt.Sprintf("Post Proc Message : %s", m)
	after(afterMessage)

	return m, err
}

func Ux(opts ...Option) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// 前処理：preprocessing
		beforeMessage := fmt.Sprintf("======= [Server Interceptor][grpc.UnaryServerInterceptor]  FullMethod is %s ===== Serer is %s", info.FullMethod, info.Server)
		before(beforeMessage)

		// Invoking the handler to complete the normal execution of a unary RPC.
		m, err := handler(ctx, req)

		// 後処理：postprocessing
		afterMessage := fmt.Sprintf("[Server Interceptor][grpc.UnaryServerInterceptor] Post Proc Message : %s", m)
		after(afterMessage)

		return m, err
	}
}

func before(message string){
	log.Printf(message)
}

func after(message string){
	log.Printf(message)
}
