package sqlite

import (
	"context"
	"github.com/idasilva/gorm-oracle"
	_ "github.com/mattn/go-sqlite3"
)


type Sqlite struct {
	context context.Context
	db gorm.SQLCommon
}

func (s Sqlite) GetName() string {
	return "sqlite3"
}

func (s *Sqlite) SetDB(db gorm.SQLCommon) {
	s.db = db
}

func (s *Sqlite) SetContext(context context.Context) {
	s.context = context
}

func (s Sqlite) BindVar(i int) string {
	panic("implement me")
}

func (s Sqlite) Quote(key string) string {
	panic("implement me")
}

func (s Sqlite) DataTypeOf(field *gorm.StructField) string {
	panic("implement me")
}

func (s Sqlite) HasIndex(tableName string, indexName string) bool {
	panic("implement me")
}

func (s Sqlite) HasForeignKey(tableName string, foreignKeyName string) bool {
	panic("implement me")
}

func (s Sqlite) RemoveIndex(tableName string, indexName string) error {
	panic("implement me")
}

func (s Sqlite) HasTable(tableName string) bool {
	panic("implement me")
}

func (s Sqlite) HasColumn(tableName string, columnName string) bool {
	panic("implement me")
}

func (s Sqlite) LimitWhereSQL(limit interface{}) string {
	panic("implement me")
}

func (s Sqlite) LimitAndOffsetSQL(limit, offset interface{}) string {
	panic("implement me")
}

func (s Sqlite) SelectFromDummyTable() string {
	panic("implement me")
}

func (s Sqlite) LastInsertIDReturningSuffix(tableName, columnName string) string {
	panic("implement me")
}

func (s Sqlite) BuildForeignKeyName(tableName, field, dest string) string {
	panic("implement me")
}

func (s Sqlite) CurrentDatabase() string {
	panic("implement me")
}

func NewDialect() gorm.Dialect {
	commontDialect := &Sqlite{}
	return commontDialect
}
