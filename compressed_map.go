package main

import (
	"fmt"
	"regexp"
)

// CompressedMap defines a map with keys such as: key.a.b, or key.[0]b.c
type extracted struct {
	req   []string
	orig  map[string]interface{}
	fresh map[string]interface{}
}

func unwrapMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		switch v.(type) {
		case map[string]interface{}:
			for ck, cv := range unwrapMap(v.(map[string]interface{})) {
				result[fmt.Sprintf("%s.%s", k, ck)] = cv
			}
		default:
			result[k] = v
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
		if yes, _ := regexp.Match(req, []byte(key)); yes {
			return true
		}
	}
	return false
}
