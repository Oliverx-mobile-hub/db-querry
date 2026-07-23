package mysqlconn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Connector struct{}

func NewConnector() Connector { return Connector{} }

func (Connector) Test(ctx context.Context, rawURL string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := Open(rawURL)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.PingContext(ctx)
}

func Open(rawURL string) (*sql.DB, error) {
	dsn, err := DriverDSN(rawURL)
	if err != nil {
		return nil, err
	}
	return sql.Open("mysql", dsn)
}

func DriverDSN(rawURL string) (string, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", errors.New("empty mysql url")
	}
	if !strings.Contains(trimmed, "://") {
		return normalizeRawDSN(trimmed), nil
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "mysql" && parsed.Scheme != "mysql2" {
		return "", fmt.Errorf("unsupported mysql url scheme: %s", parsed.Scheme)
	}
	if parsed.User == nil || parsed.User.Username() == "" {
		return "", errors.New("mysql url must include user")
	}
	dbName := strings.TrimPrefix(parsed.EscapedPath(), "/")
	if dbName == "" {
		return "", errors.New("mysql url must include database name")
	}
	dbName, err = url.PathUnescape(dbName)
	if err != nil {
		return "", err
	}

	host := parsed.Hostname()
	if host == "" {
		host = "localhost"
	}
	port := parsed.Port()
	if port == "" {
		port = "3306"
	}

	cfg := mysql.NewConfig()
	cfg.User = parsed.User.Username()
	cfg.Passwd, _ = parsed.User.Password()
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(host, port)
	cfg.DBName = dbName
	cfg.ParseTime = true
	cfg.MultiStatements = false
	cfg.AllowNativePasswords = true
	cfg.Params = map[string]string{}

	for key, values := range parsed.Query() {
		if len(values) == 0 {
			continue
		}
		switch strings.ToLower(key) {
		case "parsetime", "multistatements":
			continue
		default:
			cfg.Params[key] = values[len(values)-1]
		}
	}

	return cfg.FormatDSN(), nil
}

func normalizeRawDSN(dsn string) string {
	if strings.Contains(strings.ToLower(dsn), "multistatements=true") {
		dsn = strings.ReplaceAll(dsn, "multiStatements=true", "multiStatements=false")
		dsn = strings.ReplaceAll(dsn, "multistatements=true", "multiStatements=false")
	}
	separator := "?"
	if strings.Contains(dsn, "?") {
		separator = "&"
	}
	if !strings.Contains(strings.ToLower(dsn), "parsetime=") {
		dsn += separator + "parseTime=true"
		separator = "&"
	}
	if !strings.Contains(strings.ToLower(dsn), "multistatements=") {
		dsn += separator + "multiStatements=false"
	}
	return dsn
}
