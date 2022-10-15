package oracle

import (
	"database/sql"
	"fmt"
	"github.com/idasilva/gorm-oracle"
"reflect"
	"strconv"
	"strings"
)



var dialectsMap = map[string]gorm_oracle.Dialect{}

func NewDialect(name string) gorm_oracle.Dialect {
	if value, ok := dialectsMap[name]; ok {
		dialect := reflect.New(reflect.TypeOf(value).Elem()).Interface().(gorm_oracle.Dialect)
		return dialect
	}

	fmt.Printf("`%v` is not officially supported, running under compatibility mode.\n", name)
	commontDialect := &commonDialect{}
	return commontDialect
}

// RegisterDialect register new dialect
func RegisterDialect(name string, dialect gorm_oracle.Dialect) {
	dialectsMap[name] = dialect
}

// ParseFieldStructForDialect get field's sql data type
var ParseFieldStructForDialect = func(field *gorm_oracle.StructField, dialect gorm_oracle.Dialect) (fieldValue reflect.Value, sqlType string, size int, additionalType string) {
	// Get redirected field type
	var (
		reflectType = field.Struct.Type
		dataType    = field.TagSettings["TYPE"]
	)

	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	// Get redirected field value
	fieldValue = reflect.Indirect(reflect.New(reflectType))

	if gormDataType, ok := fieldValue.Interface().(interface {
		GormDataType(gorm_oracle.Dialect) string
	}); ok {
		dataType = gormDataType.GormDataType(dialect)
	}

	// Get scanner's real value
	var getScannerValue func(reflect.Value)
	getScannerValue = func(value reflect.Value) {
		fieldValue = value
		if _, isScanner := reflect.New(fieldValue.Type()).Interface().(sql.Scanner); isScanner && fieldValue.Kind() == reflect.Struct {
			getScannerValue(fieldValue.Field(0))
		}
	}
	getScannerValue(fieldValue)

	// Default Size
	if num, ok := field.TagSettings["SIZE"]; ok {
		size, _ = strconv.Atoi(num)
	} else {
		size = 255
	}

	// Default type from tag setting
	additionalType = field.TagSettings["NOT NULL"] + " " + field.TagSettings["UNIQUE"]
	if value, ok := field.TagSettings["DEFAULT"]; ok {
		additionalType = additionalType + " DEFAULT " + value
	}

	return fieldValue, dataType, size, strings.TrimSpace(additionalType)
}
