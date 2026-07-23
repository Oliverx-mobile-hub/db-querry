package api

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	UpsertConnection(ctx context.Context, record DBConnectionRecord) error
	DeleteConnection(ctx context.Context, name string) error
	ListConnections(ctx context.Context) ([]DBConnectionRecord, error)
	GetConnection(ctx context.Context, name string) (DBConnectionRecord, error)
	UpdateMetadataStatus(ctx context.Context, name string, status DBMetadataStatus, metadataError string) error
	InsertMetadataSnapshot(ctx context.Context, dbName string, metadata MetadataDocument, objectCount int, warnings []string) error
	GetLatestMetadataSnapshot(ctx context.Context, dbName string) (MetadataDocument, *string, error)
	InsertGeneratedSQLDraft(ctx context.Context, dbName string, draft GeneratedSQLDraft) error
	Close() error
}

type Connector interface {
	Test(ctx context.Context, databaseType DatabaseType, url string) error
}

type MetadataCollector interface {
	Collect(ctx context.Context, databaseType DatabaseType, url string) (MetadataDocument, []string, error)
}

type SQLValidator interface {
	Validate(databaseType DatabaseType, sql string) SQLValidationResult
}

type QueryExecutor interface {
	Execute(ctx context.Context, dbName string, sql string) (QueryResult, error)
}

type SQLGenerator interface {
	GenerateSQL(ctx context.Context, prompt string, metadata MetadataDocument) (GeneratedSQLDraft, error)
}
