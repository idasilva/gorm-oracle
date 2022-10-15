package callbacks

import (
	"github.com/idasilva/gorm-oracle"
)

// Define callbacks for row query
func init() {
	gorm_oracle.DefaultCallback.RowQuery().Register("gorm:row_query", rowQueryCallback)
}

// queryCallback used to query data from database
func rowQueryCallback(scope *gorm_oracle.Scope) {
	if result, ok := scope.InstanceGet("row_query_result"); ok {
		scope.PrepareQuerySQL()

		if rowResult, ok := result.(*gorm_oracle.RowQueryResult); ok {
			rowResult.Row = scope.SQLDB().QueryRowContext(scope.DB().Context,scope.SQL, scope.SQLVars...)
		} else if rowsResult, ok := result.(*gorm_oracle.RowsQueryResult); ok {
			rowsResult.Rows, rowsResult.Error = scope.SQLDB().QueryContext(scope.DB().Context,scope.SQL, scope.SQLVars...)
		}
	}
}
