package gorm

import (
	"context"
	"database/sql"
)

// SQLCommon is the minimal database connection functionality gorm requires.  Implemented by *sql.DB.
type SQLCommon interface {
	ExecContext(context context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context context.Context, query string) (*sql.Stmt, error)
	QueryContext(context context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(context context.Context, query string, args ...interface{}) *sql.Row
}

type sqlDb interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type sqlTx interface {
	Commit() error
	Rollback() error
}

// Dialect interface contains behaviors that differ across SQL database
type Dialect interface {
	// GetName get dialect's name
	GetName() string

	// SetDB set db for dialect
	SetDB(db SQLCommon)

	// SetDB set db for dialect
	SetContext(context context.Context)

	// BindVar return the placeholder for actual values in SQL statements, in many dbs it is "?", Postgres using $1
	BindVar(i int) string
	// Quote quotes field name to avoid SQL parsing exceptions by using a reserved word as a field name
	Quote(key string) string
	// DataTypeOf return data's sql type
	DataTypeOf(field *StructField) string

	// HasIndex check has index or not
	HasIndex(tableName string, indexName string) bool
	// HasForeignKey check has foreign key or not
	HasForeignKey(tableName string, foreignKeyName string) bool
	// RemoveIndex remove index
	RemoveIndex(tableName string, indexName string) error
	// HasTable check has table or not
	HasTable(tableName string) bool
	// HasColumn check has column or not
	HasColumn(tableName string, columnName string) bool

	// LimitWhereSQL return generated SQL with rownum clause, as oracle has special case
	LimitWhereSQL(limit interface{}) string
	// LimitAndOffsetSQL return generated SQL with Limit and Offset, as mssql has special case
	LimitAndOffsetSQL(limit, offset interface{}) string
	// SelectFromDummyTable return select values, for most dbs, `SELECT values` just works, mysql needs `SELECT value FROM DUAL`
	SelectFromDummyTable() string
	// LastInsertIdReturningSuffix most dbs support LastInsertId, but postgres needs to use `RETURNING`
	LastInsertIDReturningSuffix(tableName, columnName string) string

	// BuildForeignKeyName returns a foreign key name for the given table, field and reference
	BuildForeignKeyName(tableName, field, dest string) string

	// CurrentDatabase return current database name
	CurrentDatabase() string
}
