package main

import "sync"

// MapComparator checks if a CompressedMap item satisfies a rule
type MapComparator interface {
	Check(key string, value interface{}) bool
}

// Filter returns keys of the map that satisfy provided StringComparator
func (m *CompressedMap) Filter(filters ...MapComparator) []string {
	resultCh := make(chan string, 1)
	wg := sync.WaitGroup{}

	for key, value := range m.v {
		wg.Add(1)
		go func(key string, value interface{}) {
			defer wg.Done()
			for _, filter := range filters {
				if !filter.Check(key, value) {
					return
				}
			}
			resultCh <- key
		}(key, value)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	validKeys := make([]string, 0, len(m.v))
	for result := range resultCh {
		validKeys = append(validKeys, result)
	}
	return validKeys
}
