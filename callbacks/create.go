package callbacks

import (
	"fmt"
	gorm "github.com/idasilva/gorm-oracle"
	"github.com/idasilva/gorm-oracle/utils"
	"strings"
)

// Define callbacks for creating
func init() {
	gorm.DefaultCallback.Create().Register("gorm:begin_transaction", beginTransactionCallback)
	gorm.DefaultCallback.Create().Register("gorm:before_create", beforeCreateCallback)
	gorm.DefaultCallback.Create().Register("gorm:save_before_associations", saveBeforeAssociationsCallback)
	gorm.DefaultCallback.Create().Register("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	gorm.DefaultCallback.Create().Register("gorm:create", createCallback)
	gorm.DefaultCallback.Create().Register("gorm:force_reload_after_create", forceReloadAfterCreateCallback)
	gorm.DefaultCallback.Create().Register("gorm:save_after_associations", saveAfterAssociationsCallback)
	gorm.DefaultCallback.Create().Register("gorm:after_create", afterCreateCallback)
	gorm.DefaultCallback.Create().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
}

// beforeCreateCallback will invoke `BeforeSave`, `BeforeCreate` method before creating
func beforeCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		scope.CallMethod("BeforeSave")
	}
	if !scope.HasError() {
		scope.CallMethod("BeforeCreate")
	}
}

// updateTimeStampForCreateCallback will set `CreatedAt`, `UpdatedAt` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		now := utils.NowFunc()

		if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
			if createdAtField.IsBlank {
				createdAtField.Set(now)
			}
		}

		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(now)
			}
		}
	}
}

// createCallback the callback used to insert data into database
func createCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		defer scope.Trace(utils.NowFunc())

		var (
			columns, placeholders        []string
			blankColumnsWithDefaultValue []string
		)

		for _, field := range scope.Fields() {
			if scope.ChangeableField(field) {
				if field.IsNormal {
					if field.IsBlank && field.HasDefaultValue {
						blankColumnsWithDefaultValue = append(blankColumnsWithDefaultValue, scope.Quote(field.DBName))
						scope.InstanceSet("gorm:blank_columns_with_default_value", blankColumnsWithDefaultValue)
					} else if !field.IsPrimaryKey || !field.IsBlank {
						columns = append(columns, scope.Quote(field.DBName))
						placeholders = append(placeholders, scope.AddToVars(field.Field.Interface()))
					}
				} else if field.Relationship != nil && field.Relationship.Kind == "belongs_to" {
					for _, foreignKey := range field.Relationship.ForeignDBNames {
						if foreignField, ok := scope.FieldByName(foreignKey); ok && !scope.ChangeableField(foreignField) {
							columns = append(columns, scope.Quote(foreignField.DBName))
							placeholders = append(placeholders, scope.AddToVars(foreignField.Field.Interface()))
						}
					}
				}
			}
		}

		var (
			returningColumn = "*"
			quotedTableName = scope.QuotedTableName()
			primaryField    = scope.PrimaryField()
			extraOption     string
		)

		if str, ok := scope.Get("gorm:insert_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		if primaryField != nil {
			returningColumn = scope.Quote(primaryField.DBName)
		}

		lastInsertIDReturningSuffix := scope.Dialect().LastInsertIDReturningSuffix(quotedTableName, returningColumn)

		if len(columns) == 0 {
			scope.Raw(fmt.Sprintf(
				"INSERT INTO %v DEFAULT VALUES%v%v",
				quotedTableName,
				utils.AddExtraSpaceIfExist(extraOption),
				utils.AddExtraSpaceIfExist(lastInsertIDReturningSuffix),
			))
		} else {
			scope.Raw(fmt.Sprintf(
				"INSERT INTO %v (%v) VALUES (%v)%v%v",
				scope.QuotedTableName(),
				strings.Join(columns, ","),
				strings.Join(placeholders, ","),
				utils.AddExtraSpaceIfExist(extraOption),
				utils.AddExtraSpaceIfExist(lastInsertIDReturningSuffix),
			))
		}

		// execute create sql
		if lastInsertIDReturningSuffix == "" || primaryField == nil {
			if result, err := scope.SQLDB().ExecContext(scope.DB().Context, scope.SQL, scope.SQLVars...); scope.Err(err) == nil {
				// set rows affected count
				scope.DB().RowsAffected, _ = result.RowsAffected()

				// set primary value to primary field
				if primaryField != nil && primaryField.IsBlank {
					if primaryValue, err := result.LastInsertId(); scope.Err(err) == nil {
						scope.Err(primaryField.Set(primaryValue))
					}
				}
			}
		} else {
			if primaryField.Field.CanAddr() {
				if err := scope.SQLDB().QueryRowContext(scope.DB().Context, scope.SQL, scope.SQLVars...).Scan(primaryField.Field.Addr().Interface()); scope.Err(err) == nil {
					primaryField.IsBlank = false
					scope.DB().RowsAffected = 1
				}
			} else {
				scope.Err(gorm.ErrUnaddressable)
			}
		}
	}
}

// forceReloadAfterCreateCallback will reload columns that having default value, and set it back to current object
func forceReloadAfterCreateCallback(scope *gorm.Scope) {
	if blankColumnsWithDefaultValue, ok := scope.InstanceGet("gorm:blank_columns_with_default_value"); ok {
		db := scope.DB().New().Table(scope.TableName()).Select(blankColumnsWithDefaultValue.([]string))
		for _, field := range scope.Fields() {
			if field.IsPrimaryKey && !field.IsBlank {
				db = db.Where(fmt.Sprintf("%v = ?", field.DBName), field.Field.Interface())
			}
		}
		db.Scan(scope.Value)
	}
}

// afterCreateCallback will invoke `AfterCreate`, `AfterSave` method after creating
func afterCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		scope.CallMethod("AfterCreate")
	}
	if !scope.HasError() {
		scope.CallMethod("AfterSave")
	}
}
