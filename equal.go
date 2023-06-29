package main

import (
	"fmt"
	"reflect"
)

type option struct {
	skipDefaults bool
	onlyFields   []string
	dropFields   []string
}

func SkipDefaults() func(*option) {
	return func(o *option) {
		o.skipDefaults = true
	}
}

func OnlyFields(fields []string) func(*option) {
	return func(o *option) {
		o.onlyFields = fields
	}
}

func WithoutFields(fields []string) func(*option) {
	return func(o *option) {
		o.dropFields = fields
	}
}

// Equal dynamically compares if two objects are equal
func Equal(a, b interface{}, opts ...func(*option)) (bool, error) {
	// Load options
	options := &option{}
	for _, optFn := range opts {
		optFn(options)
	}

	// Get unwrapper
	if len(options.onlyFields) > 0 && len(options.dropFields) > 0 {
		return false, fmt.Errorf("cannot specify both OnlyFields and WithoutFields options")
	}
	unwrapper, fields := UnwrapWith, options.onlyFields
	if len(options.dropFields) > 0 {
		unwrapper, fields = UnwrapWithout, options.dropFields
	}

	// Fetch object maps
	aMap, errA := unwrapper(a, fields...)
	if errA != nil {
		return false, errA
	}
	bMap, errB := unwrapper(b, fields...)
	if errB != nil {
		return false, errB
	}

	// Get all keys from maps
	keys := make(map[string]bool)
	for k := range aMap {
		keys[k] = true
	}
	for k := range bMap {
		keys[k] = true
	}

	// Do comparison
	for key := range keys {
		aVal, aOk := aMap[key]
		bVal, bOk := bMap[key]
		if !aOk || !bOk { // key not found, not the same
			return false, nil
		}
		if options.skipDefaults {
			if isZeroOfUnderlyingType(aVal) != isZeroOfUnderlyingType(bVal) { // defaults not equal, not the same
				return false, nil
			}
		}
		if !reflect.DeepEqual(aVal, bVal) { // values not equal, not the same
			return false, nil
		}
	}
	return true, nil
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
