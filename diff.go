package equl

import (
	"reflect"
	"strings"
)

// Result defines the result of the Diff operation
type Result struct {
	equal bool
	diffs map[string][2]interface{}
}

func (r *Result) Equal() bool {
	return r.equal
}

func (r *Result) Diffs() map[string][2]interface{} {
	return r.diffs
}

// Diff dynamically compares two objects and returns their difference. Diff is
// slower thatn Equal as it calculate the whole difference map between objects.
// Use Diff when you need to find out which fields are different.
// Objects must implement json.Marshaler interface.
func Diff(a, b interface{}, opts ...func(*options)) (*Result, error) {
	// Load options
	req := &options{}
	for _, optFn := range opts {
		optFn(req)
	}

	// Get filter function from request or fallback to default
	filter := req.filter
	if filter == nil {
		includes := sliceToMap(req.includeFields)
		ignores := sliceToMap(req.ignoreFields)
		filter = func(key string) bool {
			return defaultFilter(key, includes, ignores)
		}
	}

	// Extract compressed maps from objects based on filter function
	aMap, err := extract(a, filter)
	if err != nil {
		return nil, err
	}
	bMap, err := extract(b, filter)
	if err != nil {
		return nil, err
	}

	// Get all keys from both object maps to use for comparisons
	keys := make(map[string]bool)
	for k := range aMap {
		keys[k] = true
	}
	for k := range bMap {
		keys[k] = true
	}

	// Do comparison, TODO: this can be optimized
	diffs := make(map[string][2]interface{})
	equal := true
	for key := range keys {
		aVal, aOk := aMap[key]
		bVal, bOk := bMap[key]

		// Values are only equal if the key is present in both maps,
		// and if the values matching the same key are equal.
		if !aOk || !bOk || !reflect.DeepEqual(aVal, bVal) {
			equal = false
			diffs[key] = [2]interface{}{aVal, bVal}
		}

		// Return instantly if only the equality check was requested
		if !equal && req.onlyEqual {
			return &Result{equal: false}, nil
		}
	}

	// Return difference between maps
	return &Result{
		equal: equal,
		diffs: diffs,
	}, nil
}

// Equal dynamically compares if two objects are equal. This is faster that Diff
// since it does not need to calculate difference map. Use Equal when you need to
// find out if two objects are equal or not, without knowing the differences.
// Objects must implement json.Marshaler interface.
func Equal(a, b interface{}, opts ...func(*options)) (bool, error) {
	result, err := Diff(a, b, append(opts, withOnlyEqual())...)
	if err != nil {
		return false, err
	}
	return result.equal, err
}

// defaultFilter checks if a flattened key should be processed.
// It will compare against sub keys as well.
func defaultFilter(key string, includes, ignores map[string]bool) bool {
	toAdd := false

	subKey := ""
	for _, keyPart := range strings.Split(key, ".") {
		if keyPart == "" {
			continue
		}
		subKey = subKey + "." + keyPart

		// If the key should be ignored, ignore it instantly
		if _, ok := ignores[subKey]; ok {
			return false
		}
		// Otherwise, only mark key for addition if it is not ignored
		if _, ok := includes[subKey]; ok {
			toAdd = true
		}
	}

	return toAdd || len(includes) == 0
}

// sliceToMap converts a slice to a lookup map
func sliceToMap(slice []string) map[string]bool {
	mp := make(map[string]bool)
	for _, item := range slice {
		mp[item] = true
	}
	return mp
}
