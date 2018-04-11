package clone

import (
	"reflect"
)

func Clone(v interface{}) interface{} {
	return CloneValue(reflect.ValueOf(v)).Interface()
}

func CloneValue(v reflect.Value) reflect.Value {
	c := cloner{}
	r := c.cloneValue(v)
	return r

}

type cloner struct {
}

func (c cloner) cloneValue(v reflect.Value) reflect.Value {

	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		v = v.Elem()
		fallthrough
	case reflect.Ptr:
		if v.IsNil() {
			return v
		}

		nt := c.cloneValue(v.Elem())
		if nt.CanAddr() {
			nt = nt.Addr()
		} else if nt.CanInterface() {
			tt := nt.Interface()
			nt = reflect.ValueOf(&tt)
		}

		return nt
	case reflect.Struct:
		nf := v.NumField()
		nt := reflect.New(v.Type()).Elem()

		for i := 0; i != nf; i++ {
			mi := v.Field(i)
			if !mi.CanSet() {
				continue
			}
			mv := c.cloneValue(mi)
			nt.Field(i).Set(mv)
		}

		return nt
	case reflect.Map:
		nt := reflect.MakeMap(v.Type())

		for _, i := range v.MapKeys() {
			mi := v.MapIndex(i)
			mk := c.cloneValue(i)
			mv := c.cloneValue(mi)
			nt.SetMapIndex(mk, mv)
		}

		return nt
	case reflect.Array:
		nt := reflect.New(v.Type()).Elem()
		l := nt.Len()

		for i := 0; i != l; i++ {
			mi := v.Index(i)
			mv := c.cloneValue(mi)
			nt.Index(i).Set(mv)
		}

		return nt
	case reflect.Slice:
		l := v.Len()
		nt := reflect.MakeSlice(v.Type(), l, v.Cap())

		for i := 0; i != l; i++ {
			mi := v.Index(i)
			mv := c.cloneValue(mi)
			nt.Index(i).Set(mv)
		}

		return nt
	default:
		return v
	}

}
