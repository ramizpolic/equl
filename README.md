# equl
equl is a dynamic and performant rule-based object equality comparator best suited for checks of mutating objects often
found in Kubernetes environments.

## Usage
```go
package main

import (
	"github.com/ramizpolic/equl"
	"strings"
)

type Parent struct {
	Name  string
	Child string
}

type Object struct {
	Base   int
	Parent Parent
}

func main() {
	objA := Object{
		Base:   1,
		Parent: Parent{Name: "same", Child: "diff-A"},
	}
	objB := Object{
		Base:   1,
		Parent: Parent{Name: "same", Child: "diff-B"},
	}

	// By default, all fields between objects will be compared
	//
	result, _ := equl.Diff(objA, objB)
	result.Equal() // False
	result.Diffs() // map[.Parent.Child:[diff-A diff-B]]

	// Fields specified using WithFields will be compared, including all their children
	//
	result, _ = equl.Diff(objA, objB, equl.WithFields(".Base", ".Parent"))
	result.Equal() // False
	result.Diffs() // map[.Parent.Child:[diff-A diff-B]]

	// Fields specified using WithoutFields will be ignored, including all their children
	//
	result, _ = equl.Diff(objA, objB, equl.WithoutFields(".Parent.Child"))
	result.Equal() // True

	// Specifying both WithFields and WithoutFields allows to create dynamic rule-based comparisons,
	// for example: compare the whole .Base and .Parent structs, but ignore everything in .Parent.Child
	//
	result, _ = equl.Diff(objA, objB, equl.WithFields(".Base", ".Parent"), equl.WithoutFields(".Parent.Child"))
	result.Equal() // True

	// It is also possible to specify a custom field fitler function to decide which fields should be
	// compared and which ones ignored.
	//
	result, _ = equl.Diff(objA, objB, equl.WithFieldFilter(func(key string) bool {
		return strings.HasPrefix(key, ".Parent") // Anything that starts with .Parent
	}))
	result.Equal() // False
	result.Diffs() // map[.Parent.Child:[diff-A diff-B]]

	// Equal dynamically compares if two objects are equal. This is faster that Diff
	// since it does not need to calculate difference map. Same rules apply.
	//
	equal, _ := equl.Equal(objA, objB)                                      // False
	equal, _ := equl.Equal(objA, objB, equl.WithoutFields(".Parent.Child")) // True
}
```
