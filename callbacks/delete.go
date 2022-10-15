package callbacks

import (

"errors"
"fmt"
"github.com/idasilva/gorm-oracle"
"github.com/idasilva/gorm-oracle/utils"
)

// Define callbacks for deleting
func init() {
	gorm_oracle.DefaultCallback.Delete().Register("gorm:begin_transaction", beginTransactionCallback)
	gorm_oracle.DefaultCallback.Delete().Register("gorm:before_delete", beforeDeleteCallback)
	gorm_oracle.DefaultCallback.Delete().Register("gorm:delete", deleteCallback)
	gorm_oracle.DefaultCallback.Delete().Register("gorm:after_delete", afterDeleteCallback)
	gorm_oracle.DefaultCallback.Delete().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
}

// beforeDeleteCallback will invoke `BeforeDelete` method before deleting
func beforeDeleteCallback(scope *gorm_oracle.Scope) {
	if scope.DB().HasBlockGlobalUpdate() && !scope.HasConditions() {
		scope.Err(errors.New("Missing WHERE clause while deleting"))
		return
	}
	if !scope.HasError() {
		scope.CallMethod("BeforeDelete")
	}
}

// deleteCallback used to delete data from database or set deleted_at to current time (when using with soft delete)
func deleteCallback(scope *gorm_oracle.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedAtField, hasDeletedAtField := scope.FieldByName("DeletedAt")

		if !scope.Search.Unscoped && hasDeletedAtField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedAtField.DBName),
				scope.AddToVars(utils.NowFunc()),
				utils.AddExtraSpaceIfExist(scope.CombinedConditionSql()),
				utils.AddExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				utils.AddExtraSpaceIfExist(scope.CombinedConditionSql()),
				utils.AddExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// afterDeleteCallback will invoke `AfterDelete` method after deleting
func afterDeleteCallback(scope *gorm_oracle.Scope) {
	if !scope.HasError() {
		scope.CallMethod("AfterDelete")
	}
}
