package dgraphclient

import (
	"context"

	"encore.dev/config"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var secrets = config.Load[struct {
	DgraphEndpoint string `config:"DGRAPH_ENDPOINT,secret"` // e.g. "localhost:9081"
	DgraphTLS      bool   `config:"DGRAPH_TLS,secret"`      // optional
	DgraphAPIKey   string `config:"DGRAPH_API_KEY,secret"`  // optional
}]()

func Init() (*dgo.Dgraph, error) {
	// Build gRPC connection (TLS handling omitted for brevity)
	var opts []grpc.DialOption
	if secrets.DgraphTLS {
		creds := credentials.NewTLS(nil)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// Use insecure credentials for non‑TLS connections
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	authInterceptor := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md := metadata.Pairs("authorization", "Bearer "+secrets.DgraphAPIKey)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
	opts = append(opts, grpc.WithUnaryInterceptor(authInterceptor))

	conn, err := grpc.NewClient(secrets.DgraphEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	return dgo.NewDgraphClient(api.NewDgraphClient(conn)), nil
}
