package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

// Server - Unary Interceptor ---------------------------------------
func UnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	// 前処理：preprocessing
	log.Printf("======= [Server Interceptor]  FullMethod is %s ===== Serer is %s", info.FullMethod, info.Server)

	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)

	// 後処理：postprocessing
	log.Printf(" Post Proc Message : %s", m)
	return m, err
}
