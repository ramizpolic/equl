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
	A: "afasf",
	B: []string{"asfasf", "asfasf"},
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
		D: "asfasf",
		F: 1.26,
		M: struct {
			K map[string]int
			G int
			H struct{ A string }
		}{K: map[string]int{"asfasf": 1, "asfasfs": 2}, G: 1, H: struct{ A string }{A: "child"}},
	},
}

func main() {
	v, _ := json.Marshal(obj)
	var data map[string]interface{}
	_ = json.Unmarshal(v, &data)
	fmt.Println(unwrapMap("", data, []string{"A", "C.D", "C.M.K.*"}))
}
