package equl

import (
	"container/list"
	"encoding/json"
	"fmt"
)

// extract compresses an object into a flat map. The object must implement json.Marshaler interface.
// The resulting object will contain all the keys that pass filter function check.
//   - maps/objects will be flattened as: map{"parent": map{"child": true}} => map{"parent.child": true}
//   - arrays will be flattened as: map{"array": []{5,6,7,8}} => map{"array.0": 5, ..., "array.3": 8}
func extract(obj interface{}, filter func(key string) bool) (map[string]interface{}, error) {
	// Convert object to map
	objMap, err := objectToMap(obj)
	if err != nil {
		return nil, err
	}

	// Add the whole map to queue
	type mapItem struct {
		isRoot bool
		key    string
		value  interface{}
	}
	queue := list.New()
	queue.PushBack(mapItem{isRoot: true, value: objMap})

	// Iterate over the queue to extract key-by-key and add to
	// the compressed resulting map
	compressedMap := make(map[string]interface{})
	for {
		// Get latest extracted map item
		if queue.Len() == 0 {
			break
		}
		item := queue.Remove(queue.Back()).(mapItem)

		// Switch based on the item
		switch itemValue := item.value.(type) {
		case map[string]interface{}: // Handle maps
			for key, val := range itemValue {
				if !item.isRoot {
					key = fmt.Sprintf("%s.%s", item.key, key)
				}
				queue.PushBack(mapItem{key: key, value: val})
			}

		case []interface{}: // Handle lists
			for index, val := range itemValue {
				queue.PushBack(mapItem{key: fmt.Sprintf("%s.%d", item.key, index), value: val})
			}

		default: // Handle primitives
			key := "." + item.key
			if filter(key) {
				// TODO: if the object has not been unmarshalled properly,
				//  this will not work properly (e.g. objects)
				compressedMap[key] = itemValue
			}
		}
	}
	return compressedMap, nil
}

// objectToMap converts an object to map using json.Marshaler converter
func objectToMap(obj interface{}) (map[string]interface{}, error) {
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
