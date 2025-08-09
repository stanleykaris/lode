package knowledgegraph

import (
	"context"
	"log"

	"encore.app/dgraphclient"
	"encore.dev/types/uuid"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

type CreateNodeRequest struct {
	Name       string `json:"name"`
	EntityType string `json:"entity_type"`
}

type CreateNodeResponse struct {
	Id uuid.UUID `json:"id"`
}

func CreateNodeAPI(ctx context.Context, req *CreateNodeRequest) (*CreateNodeResponse, error) {
	driver, err := dgraphclient.Init()
	if err != nil {
		log.Printf("Failed to initialize dgraph driver: %v", err)
		return nil, err
	}
	session := driver.NewTxn()
	defer func(session *dgo.Txn, ctx context.Context) {
		err := session.Discard(ctx)
		if err != nil {
			log.Printf("Failed to discard session: %v", err)
		}
	}(session, ctx)

	_, err = session.Mutate(ctx, &api.Mutation{
		SetJson: []byte(`{"name": "` + req.Name + `", "type": "` + req.EntityType + `"}`),
	})
	if err != nil {
		log.Printf("Failed to create node: %v", err)
		return nil, err
	}

	return &CreateNodeResponse{Id: uuid.NewV5(uuid.Nil, "")}, nil
}
