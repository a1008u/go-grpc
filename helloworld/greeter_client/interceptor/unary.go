package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

func UnaryClientInterceptor(
	ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Preprocessor phase
	log.Println("[client] Method : " + method)

	// Invoking the remote method
	err := invoker(ctx, method, req, reply, cc, opts...)

	// Postprocessor phase
	log.Println("[client] reply : " , reply)

	return err
}
