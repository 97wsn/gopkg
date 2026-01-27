package grpcutil

import (
	"context"
	"log"
	"time"

	registry "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
)

type ClientOption = grpc.ClientOption

func authMiddleware(token string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				tr.RequestHeader().Set("token", token)
			}

			return handler(ctx, req)
		}
	}
}

func NewGrpcClient(ctx context.Context, endpoint string, reg *registry.Registry, options ...ClientOption) *ggrpc.ClientConn {
	log.Println("[gRPC] client:", endpoint)
	var clientOpts []grpc.ClientOption
	clientOpts = append(clientOpts, grpc.WithTimeout(15*time.Second)) // 请求超时时间
	clientOpts = append(clientOpts, grpc.WithEndpoint(endpoint))
	clientOpts = append(clientOpts, grpc.WithDiscovery(reg))
	clientOpts = append(clientOpts, grpc.WithMiddleware(
		//authMiddleware(etcdCfg.Token), // 简单鉴权
		tracing.Client(),
	))
	if len(options) > 0 {
		clientOpts = append(clientOpts, options...)
	}
	conn, err := grpc.DialInsecure(
		ctx,
		clientOpts...,
	)

	if err != nil {
		log.Panicf("[gRPC] Failed to client connection error:%v %s\n", endpoint, err)
	}

	return conn
}
