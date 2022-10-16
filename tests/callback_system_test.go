package tests

import (
	gorm	"github.com/idasilva/gorm-oracle"
"reflect"
	"runtime"
	"strings"
	"testing"
)

func equalFuncs(funcs []*func(s *gorm.Scope), fnames []string) bool {
	var names []string
	for _, f := range funcs {
		fnames := strings.Split(runtime.FuncForPC(reflect.ValueOf(*f).Pointer()).Name(), ".")
		names = append(names, fnames[len(fnames)-1])
	}
	return reflect.DeepEqual(names, fnames)
}

func create(s *gorm.Scope)        {}
func beforeCreate1(s *gorm.Scope) {}
func beforeCreate2(s *gorm.Scope) {}
func afterCreate1(s *gorm.Scope)  {}
func afterCreate2(s *gorm.Scope)  {}

func TestRegisterCallback(t *testing.T) {
	var callback = &gorm.Callback{}

	callback.Create().Register("before_create1", beforeCreate1)
	callback.Create().Register("before_create2", beforeCreate2)
	callback.Create().Register("create", create)
	callback.Create().Register("after_create1", afterCreate1)
	callback.Create().Register("after_create2", afterCreate2)

	if !equalFuncs(callback.Creates, []string{"beforeCreate1", "beforeCreate2", "create", "afterCreate1", "afterCreate2"}) {
		t.Errorf("register callback")
	}
}

func TestRegisterCallbackWithOrder(t *testing.T) {
	var callback1 = &gorm.Callback{}
	callback1.Create().Register("before_create1", beforeCreate1)
	callback1.Create().Register("create", create)
	callback1.Create().Register("after_create1", afterCreate1)
	callback1.Create().Before("after_create1").Register("after_create2", afterCreate2)
	if !equalFuncs(callback1.Creates, []string{"beforeCreate1", "create", "afterCreate2", "afterCreate1"}) {
		t.Errorf("register callback with order")
	}

	var callback2 = &gorm.Callback{}

	callback2.Update().Register("create", create)
	callback2.Update().Before("create").Register("before_create1", beforeCreate1)
	callback2.Update().After("after_create2").Register("after_create1", afterCreate1)
	callback2.Update().Before("before_create1").Register("before_create2", beforeCreate2)
	callback2.Update().Register("after_create2", afterCreate2)

	if !equalFuncs(callback2.Updates, []string{"beforeCreate2", "beforeCreate1", "create", "afterCreate2", "afterCreate1"}) {
		t.Errorf("register callback with order")
	}
}

func TestRegisterCallbackWithComplexOrder(t *testing.T) {
	var callback1 = &gorm.Callback{}

	callback1.Query().Before("after_create1").After("before_create1").Register("create", create)
	callback1.Query().Register("before_create1", beforeCreate1)
	callback1.Query().Register("after_create1", afterCreate1)

	if !equalFuncs(callback1.Queries, []string{"beforeCreate1", "create", "afterCreate1"}) {
		t.Errorf("register callback with order")
	}

	var callback2 = &gorm.Callback{}

	callback2.Delete().Before("after_create1").After("before_create1").Register("create", create)
	callback2.Delete().Before("create").Register("before_create1", beforeCreate1)
	callback2.Delete().After("before_create1").Register("before_create2", beforeCreate2)
	callback2.Delete().Register("after_create1", afterCreate1)
	callback2.Delete().After("after_create1").Register("after_create2", afterCreate2)

	if !equalFuncs(callback2.Deletes, []string{"beforeCreate1", "beforeCreate2", "create", "afterCreate1", "afterCreate2"}) {
		t.Errorf("register callback with order")
	}
}

func replaceCreate(s *gorm.Scope) {}

func TestReplaceCallback(t *testing.T) {
	var callback = &gorm.Callback{}

	callback.Create().Before("after_create1").After("before_create1").Register("create", create)
	callback.Create().Register("before_create1", beforeCreate1)
	callback.Create().Register("after_create1", afterCreate1)
	callback.Create().Replace("create", replaceCreate)

	if !equalFuncs(callback.Creates, []string{"beforeCreate1", "replaceCreate", "afterCreate1"}) {
		t.Errorf("replace callback")
	}
}

func TestRemoveCallback(t *testing.T) {
	var callback = &gorm.Callback{}

	callback.Create().Before("after_create1").After("before_create1").Register("create", create)
	callback.Create().Register("before_create1", beforeCreate1)
	callback.Create().Register("after_create1", afterCreate1)
	callback.Create().Remove("create")

	if !equalFuncs(callback.Creates, []string{"beforeCreate1", "afterCreate1"}) {
		t.Errorf("remove callback")
	}
}
