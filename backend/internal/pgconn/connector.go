package pgconn

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Connector struct{}

func NewConnector() Connector { return Connector{} }

func (Connector) Test(ctx context.Context, url string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	return conn.Ping(ctx)
}

