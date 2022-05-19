package EqualExclude

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func EqualExclude(t *testing.T, v1, v2 interface{}, s ...string) {
	for _, d := range s {
		q1 := reflect.ValueOf(v1)
		if q1.Kind() != reflect.Slice && q1.Kind() != reflect.Map {
			q1 = q1.Elem()
		}

		q2 := reflect.ValueOf(v2)
		if q2.Kind() != reflect.Slice && q2.Kind() != reflect.Map {
			q2 = q2.Elem()
		}

		setZeroValue(q1, d)
		setZeroValue(q2, d)
	}

	require.Equal(t, v1, v2)
}

func setZeroValue(v reflect.Value, s string) {
	if i := strings.Index(s, "."); i == -1 {
		//if s == "*" {
		//	for j := 0; j < v.Len(); j++ {
		//		v1 := v.Index(j)
		//		t2(v1)
		//	}
		//} else {
		t1 := t1(v, s)
		if v.Kind() == reflect.Map {
			v.SetMapIndex(reflect.ValueOf(s), reflect.Zero(t1.Type()))
		} else {
			t2(t1)
		}
		//}
	} else {
		if s[:i] == "*" {
			for j := 0; j < v.Len(); j++ {
				v1 := v.Index(j)
				if v1.Kind() == reflect.Interface {
					if v1.Elem().Kind() == reflect.Slice {
						setZeroValue(v1.Elem(), s[i+1:])
					}
				} else {
					setZeroValue(v1, s[i+1:])
				}
			}
		} else {
			setZeroValue(t1(v, s[:i]), s[i+1:])
		}
	}
}

func t1(v reflect.Value, s string) reflect.Value {
	if v.Kind() == reflect.Slice {
		i, _ := strconv.Atoi(s)
		f := v.Index(i)
		if f.Kind() == reflect.Interface && f.Elem().Kind() == reflect.Slice {
			return f.Elem()
		}
		return f
	} else if v.Kind() == reflect.Map {
		return v.MapIndex(reflect.ValueOf(s))
	}
	return v.FieldByName(s)
}

func t2(field reflect.Value) {
	val := reflect.Zero(field.Type())
	if field.CanSet() {
		field.Set(val)
	} else {
		reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(val)
	}
}
