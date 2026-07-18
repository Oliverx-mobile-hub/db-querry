package api

type DBMetadataStatus string

const (
	MetadataPending DBMetadataStatus = "pending"
	MetadataReady   DBMetadataStatus = "ready"
	MetadataFailed  DBMetadataStatus = "failed"
)

type DBConnectionRecord struct {
	Name              string
	DatabaseType      string
	URL               string
	DisplayDSN        string
	MetadataStatus    DBMetadataStatus
	MetadataError     string
	MetadataUpdatedAt *string
	CreatedAt         string
	UpdatedAt         string
}

type DBSummary struct {
	Name              string           `json:"name"`
	DatabaseType      string           `json:"databaseType"`
	DisplayDSN        string           `json:"displayDsn"`
	MetadataStatus    DBMetadataStatus `json:"metadataStatus"`
	ConnectionStatus  string           `json:"connectionStatus"`
	MetadataUpdatedAt *string          `json:"metadataUpdatedAt"`
}

type MetadataDocument struct {
	DatabaseType string           `json:"databaseType"`
	Schemas      []MetadataSchema `json:"schemas"`
}

type MetadataSchema struct {
	Name    string           `json:"name"`
	Objects []MetadataObject `json:"objects"`
}

type MetadataObject struct {
	Schema  string           `json:"schema"`
	Name    string           `json:"name"`
	Type    string           `json:"type"`
	Comment string           `json:"comment"`
	Columns []MetadataColumn `json:"columns"`
}

type MetadataColumn struct {
	Name       string `json:"name"`
	DataType   string `json:"dataType"`
	Nullable   bool   `json:"nullable"`
	PrimaryKey bool   `json:"primaryKey"`
	Ordinal    int    `json:"ordinal"`
	Comment    string `json:"comment"`
}

type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SQLValidationResult struct {
	Valid         bool              `json:"valid"`
	Executable    bool              `json:"executable"`
	StatementType string            `json:"statementType"`
	NormalizedSQL string            `json:"normalizedSql"`
	LimitApplied  bool              `json:"limitApplied"`
	Limit         *int              `json:"limit"`
	Errors        []ValidationError `json:"errors"`
}

type QueryColumn struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type QueryResult struct {
	Columns      []QueryColumn        `json:"columns"`
	Rows         []map[string]any     `json:"rows"`
	RowCount     int                  `json:"rowCount"`
	DurationMs   int64                `json:"durationMs"`
	LimitApplied bool                 `json:"limitApplied"`
	Limit        *int                 `json:"limit"`
	Empty        bool                 `json:"empty"`
	Validation   *SQLValidationResult `json:"validation,omitempty"`
}

type GeneratedSQLDraft struct {
	Prompt            string              `json:"prompt"`
	SQL               string              `json:"sql"`
	Explanation       string              `json:"explanation"`
	ReferencedObjects []string            `json:"referencedObjects"`
	Validation        SQLValidationResult `json:"validation"`
}

func summaryFromRecord(record DBConnectionRecord) DBSummary {
	return DBSummary{
		Name:              record.Name,
		DatabaseType:      record.DatabaseType,
		DisplayDSN:        record.DisplayDSN,
		MetadataStatus:    record.MetadataStatus,
		ConnectionStatus:  "unknown",
		MetadataUpdatedAt: record.MetadataUpdatedAt,
	}
}
