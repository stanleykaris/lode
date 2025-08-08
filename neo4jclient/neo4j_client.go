package neo4jclient

import (
	"encore.dev/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Secrets struct {
	Neo4jUri      string `config:"NEO4J_URI,secret"`
	Neo4jUser     string `config:"NEO4J_USER,secret"`
	Neo4jPassword string `config:"NEO4J_PASSWORD,secret"`
}

var secrets = config.Load[Secrets]()

func Init() (neo4j.DriverWithContext, error) {
	// secrets are loaded at package level
	// Load secrets
	uri := secrets.Neo4jUri
	user := secrets.Neo4jUser
	password := secrets.Neo4jPassword

	auth := neo4j.BasicAuth(user, password, "")
	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
