package callbacks

import (
	"fmt"
	gorm "github.com/idasilva/gorm-oracle"
)

// Define callbacks for row query
func init() {
	gorm.DefaultCallback.RowQuery().Register("gorm:row_query", rowQueryCallback)
	fmt.Println("row query")
}

// queryCallback used to query data from database
func rowQueryCallback(scope *gorm.Scope) {
	if result, ok := scope.InstanceGet("row_query_result"); ok {
		scope.PrepareQuerySQL()

		if rowResult, ok := result.(*gorm.RowQueryResult); ok {
			rowResult.Row = scope.SQLDB().QueryRowContext(scope.DB().Context,scope.SQL, scope.SQLVars...)
		} else if rowsResult, ok := result.(*gorm.RowsQueryResult); ok {
			rowsResult.Rows, rowsResult.Error = scope.SQLDB().QueryContext(scope.DB().Context,scope.SQL, scope.SQLVars...)
		}
	}
}
