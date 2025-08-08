package knowledgegraph

import (
	"context"
	"log"
	"encore.dev/types/uuid"
	"encore.app/neo4jclient"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CreateNodeRequest struct {
	Name string `json:"name"`
	EntityType string `json:"entity_type"`
}

type CreateNodeResponse struct {
	Id uuid.UUID `json:"id"`
}

func CreateNodeAPI(ctx context.Context, req *CreateNodeRequest) (*CreateNodeResponse, error) {
	driver, err := neo4jclient.Init()
	if err != nil {
		log.Printf("Failed to initialize neo4j driver: %v", err)
		return nil, err
	}
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, "CREATE (n {name: $name, type: $type}) RETURN id(n)", map[string]any{
			"name": req.Name,
			"type": req.EntityType,
		})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		log.Printf("Failed to create node: %v", err)
		return nil, err
	}

	return &CreateNodeResponse{Id: uuid.NewV5(uuid.Nil, "")}, nil
}