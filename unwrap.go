package equl

import (
	"container/list"
	"encoding/json"
	"fmt"
	"strings"
)

// Unwrap unwraps an object into a flat map. The object must implement json.Marshaler interface.
// The resulting object will contain all the keys that pass filter function check.
// TODO: add example usage
func Unwrap(obj interface{}, filter func(key string) bool) (map[string]interface{}, error) {
	data, err := objToMap(obj)
	if err != nil {
		return nil, err
	}
	return UnwrapMap(data, filter), nil
}

// UnwrapMap does the same as Unwrap, but on a map.
func UnwrapMap(data map[string]interface{}, filter func(key string) bool) map[string]interface{} {
	return flatten(data, filter)
}

// UnwrapWithout unwraps an object into a flat map. The object must implement json.Marshaler interface.
// The resulting object will be formed without all the keys and subkeys provided via fieldSkip argument.
// TODO: add example usage
func UnwrapWithout(obj interface{}, fieldSkip ...string) (map[string]interface{}, error) {
	data, err := objToMap(obj)
	if err != nil {
		return nil, err
	}
	return UnwrapMapWithout(data, fieldSkip...), nil
}

// UnwrapMapWithout does the same as UnwrapWithout, but on a map.
func UnwrapMapWithout(data map[string]interface{}, fieldSkip ...string) map[string]interface{} {
	skipMap := listToMap(fieldSkip) // memoize skip request data
	return flatten(data, func(key string) bool {
		return len(skipMap) > 0 && !matches(key, skipMap) || len(skipMap) == 0
	})
}

// UnwrapWith unwraps an object into a flat map. The object must implement json.Marshaler interface.
// The resulting object will be formed only with the keys that match provided fieldReq argument.
// TODO: add example usage
func UnwrapWith(obj interface{}, fieldReq ...string) (map[string]interface{}, error) {
	data, err := objToMap(obj)
	if err != nil {
		return nil, err
	}
	return UnwrapMapWith(data, fieldReq...), nil
}

// UnwrapMapWith does the same as UnwrapWith, but on a map.
func UnwrapMapWith(data map[string]interface{}, fieldReq ...string) map[string]interface{} {
	needMap := listToMap(fieldReq) // memoize need request data
	return flatten(data, func(key string) bool {
		return len(needMap) > 0 && matches(key, needMap) || len(needMap) == 0
	})
}

// flatten flattens a given map. Extracted key will be added to resulting map if filterFn returns true.
//   - maps/objects will be flattened as: map{"parent": map{"child": true}} => map{"parent.child": true}
//   - arrays will be flattened as: map{"array": []{5,6,7,8}} => map{"array.0": 5, ..., "array.3": 8}
func flatten(input map[string]interface{}, filterFn func(key string) bool) map[string]interface{} {
	resultMap := make(map[string]interface{})

	// TODO: check with concurrency
	type mapItem struct {
		isRoot bool
		key    string
		value  interface{}
	}
	items := list.New()
	items.PushBack(mapItem{isRoot: true, value: input})

	for {
		// get latest extracted map item
		if items.Len() == 0 {
			break
		}
		item := items.Remove(items.Back()).(mapItem)

		switch itemValue := item.value.(type) {
		// handle maps
		case map[string]interface{}:
			for key, val := range itemValue {
				if !item.isRoot {
					key = fmt.Sprintf("%s.%s", item.key, key)
				}
				items.PushBack(mapItem{key: key, value: val})
			}
		// handle lists
		case []interface{}:
			for index, val := range itemValue {
				items.PushBack(mapItem{key: fmt.Sprintf("%s.%d", item.key, index), value: val})
			}
		// handle primitives
		default:
			if filterFn(item.key) {
				// TODO: if the object has not been unmarshalled properly,
				//  this will not work properly (e.g. objects)
				resultMap[item.key] = itemValue
			}
		}
	}

	return resultMap
}

// matches checks if a given flattened key should be processed. It will compare
// against sub keys as well.
//
// process = matches("key.to.check", map[string]bool{"key.to.check": true}) // true
// process = matches("key.to.check", map[string]bool{"key.to": true}) // true
// process = matches("key.to.check", map[string]bool{"key": true}) // true
//
// process = matches("key.to.check", map[string]bool{"to.check": true}) // false
// process = matches("key.to.check", map[string]bool{"check": true}) // false
func matches(key string, reqs map[string]bool) bool {
	subKey := ""
	for i, keyPart := range strings.Split(key, ".") {
		if i == 0 {
			subKey = keyPart
		} else {
			subKey = subKey + "." + keyPart
		}
		if _, ok := reqs[subKey]; ok {
			return true
		}
	}
	return false
}

func listToMap(lst []string) map[string]bool {
	mp := make(map[string]bool)
	for _, item := range lst {
		mp[item] = true
	}
	return mp
}

func objToMap(obj interface{}) (map[string]interface{}, error) {
	v, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(v, &data); err != nil {
		return nil, err
	}
	return data, nil
}
