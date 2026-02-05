package grpcutil

import (
	"context"
	"log"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
)

type Options struct {
	options    []grpc.ClientOption
	middleware []middleware.Middleware
}

type Option func(o *Options)

func WithDiscovery(discovery registry.Discovery) Option {
	return func(o *Options) {
		o.options = append(o.options, grpc.WithDiscovery(discovery))
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.options = append(o.options, grpc.WithTimeout(timeout))
	}
}

func WithMiddleware(middleware ...middleware.Middleware) Option {
	return func(o *Options) {
		o.middleware = append(o.middleware, middleware...)
	}
}

func WithClientOptions(options ...grpc.ClientOption) Option {
	return func(o *Options) {
		o.options = append(o.options, options...)
	}
}

func Auth(token string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				tr.RequestHeader().Set("token", token)
			}

			return handler(ctx, req)
		}
	}
}

// NewClient 创建一个grpc client，服务发现需要从外部实现
func NewClient(ctx context.Context, endpoint string, options ...Option) *ggrpc.ClientConn {
	log.Println("[gRPC] client:", endpoint)

	var opts Options
	for _, o := range options {
		o(&opts)
	}

	opts.middleware = append(opts.middleware, tracing.Client())

	var clientOpts []grpc.ClientOption
	clientOpts = append(clientOpts, grpc.WithTimeout(15*time.Second)) // 请求超时时间
	clientOpts = append(clientOpts, grpc.WithEndpoint(endpoint))
	clientOpts = append(clientOpts, grpc.WithMiddleware(opts.middleware...))
	if len(opts.options) > 0 {
		clientOpts = append(clientOpts, opts.options...)
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
