export type ApiError = {
  code: string
  message: string
  details: Record<string, unknown>
}

export type ApiResponse<T> =
  | { success: true; data: T; error: null }
  | { success: false; data: null; error: ApiError }

export type DatabaseType = 'postgres' | 'mysql'

export type DbSummary = {
  name: string
  databaseType: DatabaseType
  displayDsn: string
  metadataStatus: 'pending' | 'ready' | 'failed'
  connectionStatus: 'online' | 'offline' | 'unknown'
  metadataUpdatedAt: string | null
}

export type MetadataColumn = {
  name: string
  dataType: string
  nullable: boolean
  primaryKey: boolean
  ordinal: number
  comment: string
}

export type MetadataObject = {
  schema: string
  name: string
  type: 'table' | 'view'
  comment: string
  columns: MetadataColumn[]
}

export type MetadataSchema = {
  name: string
  objects: MetadataObject[]
}

export type MetadataDocument = {
  databaseType: DatabaseType
  schemas: MetadataSchema[]
}

export type SqlValidationResult = {
  valid: boolean
  executable: boolean
  statementType: 'select' | 'unknown'
  normalizedSql: string
  limitApplied: boolean
  limit: number | null
  errors: Array<{ code: string; message: string }>
}

export type QueryResult = {
  columns: Array<{ name: string; dataType: string }>
  rows: Array<Record<string, unknown>>
  rowCount: number
  durationMs: number
  limitApplied: boolean
  limit: number | null
  empty: boolean
  validation?: SqlValidationResult
}

export type ExportFormat = 'csv' | 'json'

export type ExportFile = {
  filename: string
  mimeType: string
  content: string
}

export type ExportContext = {
  dbName: string | null
  result: QueryResult | null
  loading: boolean
  now: Date
}

export type GeneratedSqlDraft = {
  prompt: string
  sql: string
  explanation: string
  referencedObjects: string[]
  validation: SqlValidationResult
}

export type DbMetadataResponse = {
  name: string
  metadataStatus: DbSummary['metadataStatus']
  connectionStatus: DbSummary['connectionStatus']
  metadataUpdatedAt: string | null
  metadata: MetadataDocument
}
