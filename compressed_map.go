package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// CompressedMap defines a map with keys such as: key.a.b, or key.[0]b.c
type extracted struct {
	req   []string
	orig  map[string]interface{}
	fresh map[string]interface{}
}

func UnwrapWithout(obj interface{}, fieldSkip ...string) (map[string]interface{}, error) {
	v, _ := json.Marshal(obj)
	var data map[string]interface{}
	_ = json.Unmarshal(v, &data)
	return unwrapMap("", data, []string{}, fieldSkip), nil
}

func UnwrapWith(obj interface{}, fieldReq ...string) (map[string]interface{}, error) {
	v, _ := json.Marshal(obj)
	var data map[string]interface{}
	_ = json.Unmarshal(v, &data)
	return unwrapMap("", data, fieldReq, []string{}), nil
}

func unwrapMap(parent string, data map[string]interface{}, req []string, skip []string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		newK := fmt.Sprintf("%s.%s", parent, k)
		if parent == "" {
			newK = k
		}
		// TODO: this needs to be simplified
		switch v.(type) {
		case map[string]interface{}:
			for ck, cv := range unwrapMap(newK, v.(map[string]interface{}), req, skip) {
				result[ck] = cv
			}
		case []interface{}:
			// TODO: selectors not working for slices
			for ci, cv := range v.([]interface{}) {
				for ck, cv := range unwrapMap(fmt.Sprintf("%s.%d", newK, ci), map[string]interface{}{"THISFLAG": cv},
					append(req, ".*THISFLAG.*"), skip) {
					result[strings.TrimRight(ck, "THISFLAG")] = cv
				}
			}
		default:
			if len(req) > 0 && matches(newK, req) {
				result[newK] = v
			} else if len(skip) > 0 && !matches(newK, skip) {
				result[newK] = v
			}
		}
	}
	return result
}

func newExtractedWith(orig map[string]interface{}, req []string) *extracted {
	r := &extracted{
		req:   req,
		orig:  orig,
		fresh: make(map[string]interface{}),
	}
	return r
}

func matches(key string, reqs []string) bool {
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
