package equl

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestExtract(t *testing.T) {
	obj := object{
		Base:        1,
		Slice:       []string{"a", "b"},
		Parent:      parent{Name: "parent-name", Child: "parent-child"},
		Grandparent: grandparent{Parent: parent{Name: "grand-parent-name", Child: "grand-parent-child"}},
	}
	objMap := map[string]interface{}{
		"Base":  1,
		"Slice": []string{"a", "b"},
		"Parent": map[string]interface{}{
			"Name":  "parent-name",
			"Child": "parent-child",
		},
		"Grandparent": map[string]interface{}{
			"Parent": map[string]interface{}{
				"Name":  "grand-parent-name",
				"Child": "grand-parent-child",
			},
		},
	}
	filterFunc := func(key string) bool {
		return key != ".Base" && !strings.Contains(key, ".Grandparent") // skip base and grandparent
	}
	expected := map[string]interface{}{
		".Parent.Child": "parent-child",
		".Parent.Name":  "parent-name",
		".Slice.0":      "a",
		".Slice.1":      "b",
	}

	actualObj, err := extract(obj, filterFunc)
	assert.Nil(t, err)
	assert.Equal(t, expected, actualObj)

	actualMap, err := extract(objMap, filterFunc)
	assert.Nil(t, err)
	assert.Equal(t, expected, actualMap)
}
