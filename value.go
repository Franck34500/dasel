package dasel

import (
	"fmt"
	"log"
	"reflect"
)

var deletePlaceholder = reflect.ValueOf("__dasel.delete_placeholder__")

// Value is a wrapper around reflect.Value that adds some handy helper funcs.
type Value struct {
	reflect.Value
	setFn    func(value Value)
	deleteFn func()
	metadata map[string]interface{}
}

// Metadata returns the metadata with a key of key for v.
func (v Value) Metadata(key string) interface{} {
	if m, ok := v.metadata[key]; ok {
		return m
	}
	return nil
}

// WithMetadata sets the given value into the values metadata.
func (v Value) WithMetadata(key string, value interface{}) Value {
	if v.metadata == nil {
		v.metadata = map[string]interface{}{}
	}
	v.metadata[key] = value
	return v
}

// Interface returns the interface{} value of v.
func (v Value) Interface() interface{} {
	return v.Unpack().Interface()
}

// Len returns v's length.
func (v Value) Len() int {
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Unpack().Len()
	case reflect.Bool:
		if v.Interface() == true {
			return 1
		} else {
			return 0
		}
	default:
		return len(fmt.Sprint(v.Interface()))
	}
}

// IsEmpty returns true is v represents an empty reflect.Value.
func (v Value) IsEmpty() bool {
	return unpackReflectValue(v.Value) == reflect.Value{}
}

// IsDeletePlaceholder returns true is v represents a delete placeholder.
func (v Value) IsDeletePlaceholder() bool {
	return unpackReflectValue(v.Value) == deletePlaceholder
}

// Kind returns the underlying type of v.
func (v Value) Kind() reflect.Kind {
	return v.Unpack().Kind()
}

func unpackReflectValue(value reflect.Value) reflect.Value {
	res := value
	for res.Kind() == reflect.Ptr || res.Kind() == reflect.Interface {
		res = res.Elem()
	}
	return res
}

// Unpack returns the underlying reflect.Value after resolving any pointers or interface types.
func (v Value) Unpack() reflect.Value {
	return unpackReflectValue(v.Value)
}

func (v Value) Type() reflect.Type {
	return v.Unpack().Type()
}

// Set sets underlying value of v.
// Depends on setFn since the implementation can differ depending on how the Value was initialised.
// Will panic if no setFn is present.
func (v Value) Set(value Value) {
	if v.setFn != nil {
		v.setFn(value)
		return
	}
	log.Println("unable to set value with missing setFn")
}

func (v Value) Delete() {
	if v.deleteFn != nil {
		v.deleteFn()
		return
	}
	log.Println("unable to delete value with missing deleteFn")
}

// MapIndex returns the value associated with key in the map v.
// It returns the zero Value if no field was found.
func (v Value) MapIndex(key Value) Value {
	return Value{
		Value: v.Unpack().MapIndex(key.Value),
		setFn: func(value Value) {
			v.Unpack().SetMapIndex(key.Value, value.Value)
		},
		deleteFn: func() {
			v.Unpack().SetMapIndex(key.Value, reflect.Value{})
		},
		metadata: map[string]interface{}{
			"type":   unpackReflectValue(v.Unpack().MapIndex(key.Value)).Kind().String(),
			"key":    key.Interface(),
			"parent": v,
		},
	}
}

func (v Value) MapKeys() []Value {
	res := make([]Value, 0)
	for _, k := range v.Unpack().MapKeys() {
		res = append(res, Value{Value: k})
	}
	return res
}

// FieldByName returns the struct field with the given name.
// It returns the zero Value if no field was found.
func (v Value) FieldByName(name string) Value {
	return Value{
		Value: v.Unpack().FieldByName(name),
		setFn: func(value Value) {
			v.Unpack().FieldByName(name).Set(value.Value)
		},
		deleteFn: func() {
			field := v.Unpack().FieldByName(name)
			field.Set(reflect.New(field.Type()))
		},
		metadata: map[string]interface{}{
			"type":   unpackReflectValue(v.Unpack().FieldByName(name)).Kind().String(),
			"key":    name,
			"parent": v,
		},
	}
}

// NumField returns the number of fields in the struct v.
func (v Value) NumField() int {
	return v.Unpack().NumField()
}

// Index returns v's i'th element.
// It panics if v's Kind is not Array, Slice, or String or i is out of range.
func (v Value) Index(i int) Value {
	return Value{
		Value: v.Unpack().Index(i),
		setFn: func(value Value) {
			v.Unpack().Index(i).Set(value.Value)
		},
		deleteFn: func() {
			// todo : find a way to remove this slice element
			// The slice index, v, and v.Metadata("parent") are all returning false for CanSet.
			v.Unpack().Index(i).Set(deletePlaceholder)
			return
		},
		metadata: map[string]interface{}{
			"type":   unpackReflectValue(v.Unpack().Index(i)).Kind().String(),
			"key":    i,
			"parent": v,
		},
	}
}

// Values represents a list of Value's.
type Values []Value

// Interfaces returns the interface values for the underlying values stored in v.
func (v Values) Interfaces() []interface{} {
	res := make([]interface{}, 0)
	for _, val := range v {
		res = append(res, val.Interface())
	}
	return res
}

// ValueOf returns a Value wrapped around value.
func ValueOf(value interface{}) Value {
	switch v := value.(type) {
	case Value:
		return v
	case reflect.Value:
		return Value{
			Value: v,
		}
	default:
		return Value{
			Value: reflect.ValueOf(value),
		}
	}
}