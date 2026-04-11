package docproxy

import (
	"math"
	"reflect"
	"strings"
	"time"
)

// As converts a (Value, bool) pair — as returned by Get and similar methods —
// into a typed Go value T. It is the primary way to extract values without
// explicit type assertions on view types.
//
// Returns (zero, false) if ok is false (key absent), or if the conversion from
// v's type to T is not supported. For NullView, returns (zero, false) for
// concrete targets and (nil, true) for pointer targets.
//
// Supported conversions:
//
//   - string:       StringView, TextView
//   - bool:         BoolView
//   - int…int64:    Int64View, Uint64View, Float64View, CounterView — range-checked
//   - uint…uint64:  Int64View (≥0), Uint64View, Float64View — range-checked
//   - float32/64:   Float64View, Int64View, Uint64View
//   - []byte:       BytesView
//   - time.Time:    TimestampView
//   - struct:       MapView, via exported fields and `automerge:"name"` tags
//   - map[string]T: MapView, values converted recursively
//   - []T:          ListView, elements converted recursively
//   - any:          any Value, via [Value.Native]
//   - *T:           same as T but wrapped in a pointer; NullView → nil (ok=true)
//   - MapView, ListView, TextView: direct pass-through
//
// Example:
//
//	name, ok := docproxy.As[string](doc.Get("name"))
//
//	type Config struct {
//	    Debug   bool   `automerge:"debug"`
//	    Version int64  `automerge:"version"`
//	}
//	cfg, ok := docproxy.As[Config](doc.Get("config"))
func As[T any](v Value, ok bool) (T, bool) {
	var zero T
	if !ok {
		return zero, false
	}
	t := reflect.TypeOf(&zero).Elem()
	result, converted := asValue(v, t)
	if !converted {
		return zero, false
	}
	typed, cast := result.(T)
	if !cast {
		return zero, false
	}
	return typed, true
}

// asValue is the recursive core of As. It converts Value v to reflect.Type t.
func asValue(v Value, t reflect.Type) (any, bool) {
	// Empty interface (any): return the native Go representation.
	if t.Kind() == reflect.Interface && t.NumMethod() == 0 {
		if v == nil {
			return nil, false
		}
		return v.Native(), true
	}

	// Pointer target: unwrap, recurse, re-wrap. NullView → typed nil pointer.
	if t.Kind() == reflect.Ptr {
		if _, isNull := v.(NullView); isNull {
			return reflect.Zero(t).Interface(), true
		}
		inner, ok := asValue(v, t.Elem())
		if !ok {
			return nil, false
		}
		ptr := reflect.New(t.Elem())
		ptr.Elem().Set(reflect.ValueOf(inner))
		return ptr.Interface(), true
	}

	// Null for non-pointer targets always fails.
	if _, isNull := v.(NullView); isNull {
		return nil, false
	}
	if v == nil {
		return nil, false
	}

	switch t.Kind() {
	case reflect.String:
		switch sv := v.(type) {
		case StringView:
			return sv.Value(), true
		case TextView:
			return sv.Value(), true
		}

	case reflect.Bool:
		if bv, ok := v.(BoolView); ok {
			return bv.Value(), true
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convertToInt(v, t)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convertToUint(v, t)

	case reflect.Float32, reflect.Float64:
		return convertToFloat(v, t)

	case reflect.Slice:
		if t == reflect.TypeOf([]byte(nil)) {
			if bv, ok := v.(BytesView); ok {
				return bv.Value(), true
			}
			return nil, false
		}
		if lv, ok := v.(ListView); ok {
			return asSlice(lv, t)
		}

	case reflect.Map:
		if mv, ok := v.(MapView); ok {
			return asMap(mv, t)
		}

	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			if tv, ok := v.(TimestampView); ok {
				return tv.Value(), true
			}
			return nil, false
		}
		if mv, ok := v.(MapView); ok {
			return asStruct(mv, t)
		}
	}

	// Direct pass-through for view types.
	rv := reflect.ValueOf(v)
	if rv.Type().AssignableTo(t) {
		return v, true
	}
	return nil, false
}

func convertToInt(v Value, t reflect.Type) (any, bool) {
	var i64 int64
	switch sv := v.(type) {
	case Int64View:
		i64 = sv.Value()
	case Uint64View:
		if sv.Value() > math.MaxInt64 {
			return nil, false
		}
		i64 = int64(sv.Value())
	case Float64View:
		i64 = int64(sv.Value())
	case CounterView:
		i64 = sv.Value()
	default:
		return nil, false
	}
	rv := reflect.New(t).Elem()
	rv.SetInt(i64)
	if rv.Int() != i64 {
		return nil, false // overflow on narrowing
	}
	return rv.Interface(), true
}

func convertToUint(v Value, t reflect.Type) (any, bool) {
	var u64 uint64
	switch sv := v.(type) {
	case Uint64View:
		u64 = sv.Value()
	case Int64View:
		if sv.Value() < 0 {
			return nil, false
		}
		u64 = uint64(sv.Value())
	case Float64View:
		if sv.Value() < 0 {
			return nil, false
		}
		u64 = uint64(sv.Value())
	default:
		return nil, false
	}
	rv := reflect.New(t).Elem()
	rv.SetUint(u64)
	if rv.Uint() != u64 {
		return nil, false // overflow on narrowing
	}
	return rv.Interface(), true
}

func convertToFloat(v Value, t reflect.Type) (any, bool) {
	var f64 float64
	switch sv := v.(type) {
	case Float64View:
		f64 = sv.Value()
	case Int64View:
		f64 = float64(sv.Value())
	case Uint64View:
		f64 = float64(sv.Value())
	default:
		return nil, false
	}
	rv := reflect.New(t).Elem()
	rv.SetFloat(f64)
	return rv.Interface(), true
}

// asSlice converts a ListView to a slice type t (elem type is t.Elem()).
func asSlice(lv ListView, t reflect.Type) (any, bool) {
	elemType := t.Elem()
	isIface := elemType.Kind() == reflect.Interface
	result := reflect.MakeSlice(t, 0, lv.Len())
	for _, elem := range lv.Values() {
		converted, ok := asValue(elem, elemType)
		if !ok {
			return nil, false
		}
		result = reflect.Append(result, reflectElem(elemType, isIface, converted))
	}
	return result.Interface(), true
}

// asMap converts a MapView to a map type t (must be map[string]V).
func asMap(mv MapView, t reflect.Type) (any, bool) {
	if t.Key().Kind() != reflect.String {
		return nil, false
	}
	valType := t.Elem()
	isIface := valType.Kind() == reflect.Interface
	result := reflect.MakeMap(t)
	for key, elem := range mv.Values() {
		converted, ok := asValue(elem, valType)
		if !ok {
			return nil, false
		}
		result.SetMapIndex(reflect.ValueOf(key), reflectElem(valType, isIface, converted))
	}
	return result.Interface(), true
}

// asStruct converts a MapView to a struct type t using automerge struct tags.
// Missing keys are left as zero values; conversion failures are silently skipped.
func asStruct(mv MapView, t reflect.Type) (any, bool) {
	result := reflect.New(t).Elem()
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		key := structFieldKey(field)
		if key == "-" {
			continue
		}
		fv, ok := mv.Get(key)
		if !ok {
			continue
		}
		converted, ok := asValue(fv, field.Type)
		if !ok {
			continue
		}
		result.Field(i).Set(reflectElem(field.Type, field.Type.Kind() == reflect.Interface, converted))
	}
	return result.Interface(), true
}

// structFieldKey returns the map key for a struct field from its automerge tag
// or field name.
func structFieldKey(f reflect.StructField) string {
	tag := f.Tag.Get("automerge")
	if tag == "" {
		return f.Name
	}
	name, _, _ := strings.Cut(tag, ",")
	return name
}

// reflectElem wraps converted into a reflect.Value suitable for assignment to
// a target of type t. When t is an interface, the value must be boxed into an
// interface reflect.Value rather than left as its concrete type.
func reflectElem(t reflect.Type, isIface bool, converted any) reflect.Value {
	if isIface {
		ev := reflect.New(t).Elem()
		if converted != nil {
			ev.Set(reflect.ValueOf(converted))
		}
		return ev
	}
	return reflect.ValueOf(converted)
}
