package equl

import (
	"fmt"
	"testing"
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
	A: "initial",
	B: []string{"example", "example"},
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
		D: "HJK",
		F: 1.26,
		M: struct {
			K map[string]int
			G int
			H struct{ A string }
		}{K: map[string]int{"ABC": 1, "EFG": 2}, G: 1, H: struct{ A string }{A: "child"}},
	},
}

func TestEqual(t *testing.T) {
	fields := []string{"B.0", "C.D", "C.M.K"}
	fmt.Println()
	fmt.Println("=== Fields", fields)
	fmt.Println()
	fmt.Println("--- Only")
	fmt.Println(UnwrapWith(obj, fields...))
	fmt.Println()
	fmt.Println("--- Without")
	fmt.Println(UnwrapWithout(obj, fields...))
	fmt.Println()
	fmt.Println("--- Equal Only")
	fmt.Println(Equal(obj, obj, OnlyFields(fields...)))
	fmt.Println()
	fmt.Println("--- Equal Without")
	fmt.Println(Equal(obj, obj, WithoutFields(fields...)))
	fmt.Println()
	fmt.Println("--- NotEqual")
	fmt.Println(Equal(obj, Object{B: fields}, WithoutFields(fields...)))
}
