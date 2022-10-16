package callbacks

import (

"context"
"errors"
"fmt"
	gorm "github.com/idasilva/gorm-oracle"
"github.com/idasilva/gorm-oracle/utils"
"reflect"
)

// Define callbacks for querying
func init() {
	gorm.DefaultCallback.Query().Register("gorm:query", queryCallback)
	gorm.DefaultCallback.Query().Register("gorm:preload", gorm.PreloadCallback)
	gorm.DefaultCallback.Query().Register("gorm:after_query", afterQueryCallback)
	fmt.Println("query")
}

// queryCallback used to query data from database
func queryCallback(scope *gorm.Scope) {
	if scope.Context == nil{
		scope.Context = context.Background()
	}
	defer scope.Trace(utils.NowFunc())

	var (
		isSlice, isPtr bool
		resultType     reflect.Type
		results        = scope.IndirectValue()
	)

	if orderBy, ok := scope.Get("gorm:order_by_primary_key"); ok {
		if primaryField := scope.PrimaryField(); primaryField != nil {
			scope.Search.Order(fmt.Sprintf("%v.%v %v", scope.QuotedTableName(), scope.Quote(primaryField.DBName), orderBy))
		}
	}

	if value, ok := scope.Get("gorm:query_destination"); ok {
		results = utils.Indirect(reflect.ValueOf(value))
	}

	if kind := results.Kind(); kind == reflect.Slice {
		isSlice = true
		resultType = results.Type().Elem()
		results.Set(reflect.MakeSlice(results.Type(), 0, 0))

		if resultType.Kind() == reflect.Ptr {
			isPtr = true
			resultType = resultType.Elem()
		}
	} else if kind != reflect.Struct {
		scope.Err(errors.New("unsupported destination, should be slice or struct"))
		return
	}

	scope.PrepareQuerySQL()

	if !scope.HasError() {
		scope.DB().RowsAffected = 0
		if str, ok := scope.Get("gorm:query_option"); ok {
			scope.SQL += utils.AddExtraSpaceIfExist(fmt.Sprint(str))
		}

		if rows, err := scope.SQLDB().QueryContext(scope.Context,scope.SQL, scope.SQLVars...); scope.Err(err) == nil {
			defer rows.Close()

			columns, _ := rows.Columns()
			for rows.Next() {
				scope.DB().RowsAffected++

				elem := results
				if isSlice {
					elem = reflect.New(resultType).Elem()
				}

				scope.Scan(rows, columns, scope.New(elem.Addr().Interface()).Fields())

				if isSlice {
					if isPtr {
						results.Set(reflect.Append(results, elem.Addr()))
					} else {
						results.Set(reflect.Append(results, elem))
					}
				}
			}

			if err := rows.Err(); err != nil {
				scope.Err(err)
			} else if scope.DB().RowsAffected == 0 && !isSlice {
				scope.Err(gorm.ErrRecordNotFound)
			}
		}
	}
}

// afterQueryCallback will invoke `AfterFind` method after querying
func afterQueryCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		scope.CallMethod("AfterFind")
	}
}
