package field_access

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	"github.com/moisespsena-go/aorm"
)

var vfGetterType = reflect.TypeOf((*aorm.VirtualFieldsGetter)(nil)).Elem()

func get(rov reflect.Value, pth []string, index int) (r reflect.Value, ok bool) {
	for rov.Kind() == reflect.Ptr {
		rov = rov.Elem()
	}
	if rov.Kind() != reflect.Struct {
		panic(errors.New("Invalid Field path " + strconv.Quote(strings.Join(pth[0:index+1], "."))))
	}

	typ := rov.Type()

	if typ.Implements(vfGetterType) {
		if v, ok := rov.Interface().(aorm.VirtualFieldsGetter).GetVirtualField(pth[index]); ok {
			return reflect.ValueOf(v), true
		}
	}

	if f, ok := typ.FieldByName(pth[index]); ok {
		return rov.FieldByIndex(f.Index), true
	}

	return reflect.Value{}, false
}

func GetValue(o interface{}, pth string) (rov reflect.Value, ok bool) {
	parts := strings.Split(pth, ".")
	rov = reflect.ValueOf(o)
	for i := range parts {
		if rov, ok = get(rov, parts, i); !ok {
			return
		}
	}
	return
}

func Get(o interface{}, pth string) (v interface{}, ok bool) {
	if value, ok := GetValue(o, pth); ok {
		return value.Interface(), ok
	}
	return nil, false
}

func Set(o interface{}, pth string, value interface{}) (ok bool) {
	if v, ok := GetValue(o, pth); ok {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		v.Set(reflect.ValueOf(value))
		return true
	}
	return false
}
