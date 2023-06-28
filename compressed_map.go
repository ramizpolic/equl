package main

import (
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

func unwrapMap(parent string, data map[string]interface{}, req []string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		newK := fmt.Sprintf("%s.%s", parent, k)
		if parent == "" {
			newK = k
		}
		// TODO: what about the array of maps/serializable objects?
		switch v.(type) {
		case map[string]interface{}:
			for ck, cv := range unwrapMap(newK, v.(map[string]interface{}), req) {
				result[ck] = cv
			}
		default:
			if matches(newK, req) {
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
