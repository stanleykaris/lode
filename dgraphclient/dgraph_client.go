package dgraphclient

import (
	"context"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	Endpoint string
	TLS      bool
	APIKey   string
}

func New(cfg Config) (*dgo.Dgraph, error) {
	// Build gRPC connection (TLS handling)
	var opts []grpc.DialOption
	if cfg.TLS {
		creds := credentials.NewTLS(nil)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// Use insecure credentials for non‑TLS connections
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Optional API key propagation
	if cfg.APIKey != "" {
		authInterceptor := func(
			ctx context.Context,
			method string,
			req, reply interface{},
			cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			callOpts ...grpc.CallOption,
		) error {
			md := metadata.Pairs("authorization", "Bearer "+cfg.APIKey)
			ctx = metadata.NewOutgoingContext(ctx, md)
			return invoker(ctx, method, req, reply, cc, callOpts...)
		}
		opts = append(opts, grpc.WithUnaryInterceptor(authInterceptor))
	}

	conn, err := grpc.NewClient(cfg.Endpoint, opts...)
	if err != nil {
		return nil, err
	}
	return dgo.NewDgraphClient(api.NewDgraphClient(conn)), nil
}
