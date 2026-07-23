package dbconn

import (
	"context"
	"fmt"

	"db-querry/backend/internal/api"
	"db-querry/backend/internal/mysqlconn"
	"db-querry/backend/internal/pgconn"
)

type Connector struct {
	postgres pgconn.Connector
	mysql    mysqlconn.Connector
}

func NewConnector() Connector {
	return Connector{postgres: pgconn.NewConnector(), mysql: mysqlconn.NewConnector()}
}

func (c Connector) Test(ctx context.Context, databaseType api.DatabaseType, rawURL string) error {
	switch api.NormalizeDatabaseType(databaseType) {
	case api.DatabaseTypeMySQL:
		return c.mysql.Test(ctx, rawURL)
	case api.DatabaseTypePostgres:
		return c.postgres.Test(ctx, rawURL)
	default:
		return fmt.Errorf("unsupported database type: %s", databaseType)
	}
}
