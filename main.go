package main

import (
	"encoding/json"
	"fmt"
)

type Object struct {
	A string
	B []string
	C struct {
		D string
		F float64
		M struct {
			K map[string]int
			G int
			H struct {
				A string
			}
		}
	}
}

var obj = Object{
	A: "a",
	B: []string{"b", "c"},
	C: struct {
		D string
		F float64
		M struct {
			K map[string]int
			G int
			H struct {
				A string
			}
		}
	}{
		D: "D",
		F: 1.26,
		M: struct {
			K map[string]int
			G int
			H struct{ A string }
		}{K: map[string]int{"a": 1, "b": 2}, G: 1, H: struct{ A string }{A: "child"}},
	},
}

func main() {
	v, _ := json.Marshal(obj)
	var data map[string]interface{}
	_ = json.Unmarshal(v, &data)
	fmt.Println(unwrapMap(data))
}
