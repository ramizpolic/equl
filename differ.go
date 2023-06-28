package main

func Diff(a, b CompressedMap, keys []string) bool {
	for _, key := range keys {
		aVal, _ := a.Get(key)
		bVal, _ := b.Get(key)
		if aVal != bVal { // this can be slow
			return false
		}
	}
	return true
}
