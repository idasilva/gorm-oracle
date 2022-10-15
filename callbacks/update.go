package callbacks

import (

"errors"
"fmt"
"github.com/idasilva/gorm-oracle"
"github.com/idasilva/gorm-oracle/utils"
"strings"
)

// Define callbacks for updating
func init() {
	gorm_oracle.DefaultCallback.Update().Register("gorm:assign_updating_attributes", assignUpdatingAttributesCallback)
	gorm_oracle.DefaultCallback.Update().Register("gorm:begin_transaction", beginTransactionCallback)
	gorm_oracle.DefaultCallback.Update().Register("gorm:before_update", beforeUpdateCallback)
	gorm_oracle.DefaultCallback.Update().Register("gorm:save_before_associations", saveBeforeAssociationsCallback)
	gorm_oracle.DefaultCallback.Update().Register("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	gorm_oracle.DefaultCallback.Update().Register("gorm:update", updateCallback)
	gorm_oracle.DefaultCallback.Update().Register("gorm:save_after_associations", saveAfterAssociationsCallback)
	gorm_oracle.DefaultCallback.Update().Register("gorm:after_update", afterUpdateCallback)
	gorm_oracle.DefaultCallback.Update().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
}

// assignUpdatingAttributesCallback assign updating attributes to model
func assignUpdatingAttributesCallback(scope *gorm_oracle.Scope) {
	if attrs, ok := scope.InstanceGet("gorm:update_interface"); ok {
		if updateMaps, hasUpdate := scope.UpdatedAttrsWithValues(attrs); hasUpdate {
			scope.InstanceSet("gorm:update_attrs", updateMaps)
		} else {
			scope.SkipLeft()
		}
	}
}

// beforeUpdateCallback will invoke `BeforeSave`, `BeforeUpdate` method before updating
func beforeUpdateCallback(scope *gorm_oracle.Scope) {
	if scope.DB().HasBlockGlobalUpdate() && !scope.HasConditions() {
		scope.Err(errors.New("Missing WHERE clause while updating"))
		return
	}
	if _, ok := scope.Get("gorm:update_column"); !ok {
		if !scope.HasError() {
			scope.CallMethod("BeforeSave")
		}
		if !scope.HasError() {
			scope.CallMethod("BeforeUpdate")
		}
	}
}

// updateTimeStampForUpdateCallback will set `UpdatedAt` when updating
func updateTimeStampForUpdateCallback(scope *gorm_oracle.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdatedAt", utils.NowFunc())
	}
}

// updateCallback the callback used to update data to database
func updateCallback(scope *gorm_oracle.Scope) {
	if !scope.HasError() {
		var sqls []string

		if updateAttrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
			for column, value := range updateAttrs.(map[string]interface{}) {
				sqls = append(sqls, fmt.Sprintf("%v = %v", scope.Quote(column), scope.AddToVars(value)))
			}
		} else {
			for _, field := range scope.Fields() {
				if scope.ChangeableField(field) {
					if !field.IsPrimaryKey && field.IsNormal {
						sqls = append(sqls, fmt.Sprintf("%v = %v", scope.Quote(field.DBName), scope.AddToVars(field.Field.Interface())))
					} else if relationship := field.Relationship; relationship != nil && relationship.Kind == "belongs_to" {
						for _, foreignKey := range relationship.ForeignDBNames {
							if foreignField, ok := scope.FieldByName(foreignKey); ok && !scope.ChangeableField(foreignField) {
								sqls = append(sqls,
									fmt.Sprintf("%v = %v", scope.Quote(foreignField.DBName), scope.AddToVars(foreignField.Field.Interface())))
							}
						}
					}
				}
			}
		}

		var extraOption string
		if str, ok := scope.Get("gorm:update_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		if len(sqls) > 0 {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v%v%v",
				scope.QuotedTableName(),
				strings.Join(sqls, ", "),
				utils.AddExtraSpaceIfExist(scope.CombinedConditionSql()),
				utils.AddExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// afterUpdateCallback will invoke `AfterUpdate`, `AfterSave` method after updating
func afterUpdateCallback(scope *gorm_oracle.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		if !scope.HasError() {
			scope.CallMethod("AfterUpdate")
		}
		if !scope.HasError() {
			scope.CallMethod("AfterSave")
		}
	}
}
