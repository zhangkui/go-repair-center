package repository

import (
 "context"
 "database/sql"
)

type Record map[string]any

type Page struct {
 Items []Record
 Page int
 PageSize int
 Total int64
 TotalPages int
}

type CRUDRepository interface {
 Create(context.Context, Record) (Record, error)
 Get(context.Context, int64) (Record, error)
 List(context.Context, int, int, string, string) (Page, error)
 Update(context.Context, int64, Record) (Record, error)
 Delete(context.Context, int64) error
}

type TransactionManager interface {
 WithinTransaction(context.Context, func(*sql.Tx) error) error
}

type UserRepository interface {
 CRUDRepository
 FindByUsername(context.Context, string) (Record, error)
 Permissions(context.Context, int64) ([]string, error)
 AssignRole(context.Context, int64, int64) error
 RevokeSessions(context.Context, int64) error
}

type PartRepository interface {
 CRUDRepository
 UseStock(context.Context, *sql.Tx, int64, int, int64, int64, string) (Record, error)
 AdjustStock(context.Context, int64, int, string, int64) (Record, error)
}

type QuotationRepository interface {
 CRUDRepository
 CreateVersion(context.Context, *sql.Tx, Record, []Record) (Record, error)
 VersionHistory(context.Context, int64) ([]Record, error)
}

type AuditRepository interface {
 CRUDRepository
 Write(context.Context, *sql.Tx, Record) error
}