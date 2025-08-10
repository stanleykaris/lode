package knowledgegraph

import (
	"context"
	"encoding/json"
	"log"

	"encore.app/dgraphclient"
	"encore.dev/config"
	"encore.dev/types/uuid"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

var secrets = config.Load[struct {
	DgraphEndpoint string `config:"DGRAPH_ENDPOINT,secret"`
	DgraphTLS      bool   `config:"DGRAPH_TLS,secret"`
	DgraphAPIKey   string `config:"DGRAPH_API_KEY,secret"`
}]()

type CreateNodeRequest struct {
	Name       string `json:"name"`
	EntityType string `json:"entity_type"`
}

type CreateNodeResponse struct {
	Id uuid.UUID `json:"id"`
}

func CreateNodeAPI(ctx context.Context, req *CreateNodeRequest) (*CreateNodeResponse, error) {
	driver, err := dgraphclient.New(dgraphclient.Config{
		Endpoint: secrets.DgraphEndpoint,
		TLS:      secrets.DgraphTLS,
		APIKey:   secrets.DgraphAPIKey,
	})
	if err != nil {
		log.Printf("Failed to initialize dgraph driver: %v", err)
		return nil, err
	}
	session := driver.NewTxn()
	defer func(session *dgo.Txn, ctx context.Context) {
		if err := session.Discard(ctx); err != nil {
			log.Printf("Failed to discard session: %v", err)
		}
	}(session, ctx)

	payload := map[string]string{
		"name":       req.Name,
		"entityType": req.EntityType,
	}
	b, _ := json.Marshal(payload)
	muResp, err := session.Mutate(ctx, &api.Mutation{
		SetJson:   b,
		CommitNow: true,
	})
	if err != nil {
		log.Printf("Failed to create node: %v", err)
		return nil, err
	}

	uid := muResp.GetUids()["node"]
	return &CreateNodeResponse{Id: uuid.NewV5(uuid.Nil, uid)}, nil
}
