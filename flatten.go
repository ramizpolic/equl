package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func UnwrapWithout(obj interface{}, fieldSkip ...string) (map[string]interface{}, error) {
	v, _ := json.Marshal(obj)
	var data map[string]interface{}
	_ = json.Unmarshal(v, &data)
	return flatten(data, []string{}, fieldSkip), nil
}

func UnwrapWith(obj interface{}, fieldReq ...string) (map[string]interface{}, error) {
	v, _ := json.Marshal(obj)
	var data map[string]interface{}
	_ = json.Unmarshal(v, &data)
	return flatten(data, fieldReq, []string{}), nil
}

type mapItem struct {
	isRoot bool
	key    string
	value  interface{}
}

// flatten flattens a given map
func flatten(input map[string]interface{}, need, skip []string) map[string]interface{} {
	resultMap := make(map[string]interface{})

	stack := list.New()
	stack.PushBack(mapItem{isRoot: true, value: input})

	// keys := []string{}
	addToMap := func(key string, val interface{}) {
		if len(need) > 0 && matches(key, need) {
			resultMap[key] = val
			// keys = append(keys, key)
		} else if len(skip) > 0 && !matches(key, skip) {
			resultMap[key] = val
			// keys = append(keys, key)
		}
	}

	for {
		element := stack.Back()
		if element == nil {
			break
		}
		stack.Remove(element)
		item := element.Value.(mapItem)

		switch itemValue := item.value.(type) {
		// Handle maps
		case map[string]interface{}:
			for key, val := range itemValue {
				if !item.isRoot {
					key = fmt.Sprintf("%s.%s", item.key, key)
				}
				stack.PushBack(mapItem{key: key, value: val})
			}
		// Handle lists
		case []interface{}:
			for index, val := range itemValue {
				key := fmt.Sprintf("%s[%d]", item.key, index)
				stack.PushBack(mapItem{key: key, value: val})
			}
		// Handle primitives
		default:
			addToMap(item.key, itemValue)
		}
	}

	return resultMap
}

func matches(key string, reqs []string) bool {
	// TODO: this needs optimization, probably
	for _, req := range reqs {
		// fmt.Printf("Comparing %s with %s\n", key, req)
		if strings.Contains(req, "*") {
			if yes, _ := regexp.Match(req, []byte(key)); yes {
				return true
			}
		} else if strings.EqualFold(req, key) {
			return true
		}
	}
	return false
}
