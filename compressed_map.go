package main

// CompressedMap defines a map with keys such as: key.a.b, or key.[0]b.c
type CompressedMap struct {
	v map[string]interface{}
}

func (m *CompressedMap) Get(key string) (interface{}, error) {
	return m.v[key], nil
}
