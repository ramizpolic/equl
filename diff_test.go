package equl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type parent struct {
	Name  string
	Child string
}

type grandparent struct {
	Name   string
	Parent parent
}

type object struct {
	Base        int
	Slice       []string
	Parent      parent
	Grandparent grandparent
}

func TestDiff(t *testing.T) {
	objA := object{
		Base:        1,
		Slice:       []string{"same", "diff-A"},
		Parent:      parent{Name: "diff-A", Child: "same"},
		Grandparent: grandparent{Parent: parent{Name: "same", Child: "diff-A"}},
	}

	objB := object{
		Base:        1,
		Slice:       []string{"same", "diff-B"},
		Parent:      parent{Name: "diff-B", Child: "same"},
		Grandparent: grandparent{Parent: parent{Name: "same", Child: "diff-B"}},
	}

	allFields := []string{".Base", ".Slice", ".Parent", ".Grandparent"}
	sameFields := []string{".Base", ".Slice.0", ".Parent.Child", ".Grandparent.Parent.Name"}
	diffFields := []string{".Slice.1", ".Parent.Name", ".Grandparent.Parent.Child"}
	diffValues := map[string][2]interface{}{
		".Slice.1":                  {"diff-A", "diff-B"},
		".Parent.Name":              {"diff-A", "diff-B"},
		".Grandparent.Parent.Child": {"diff-A", "diff-B"},
	}

	result, err := Diff(objA, objB)
	assert.Nil(t, err)
	assert.False(t, result.Equal(), "default should be false")
	assert.Equal(t, result.Diffs(), diffValues)

	result, err = Diff(objA, objB, WithFields(allFields...))
	assert.Nil(t, err)
	assert.False(t, result.Equal(), "allFields included should be false")
	assert.Equal(t, result.Diffs(), diffValues)

	result, err = Diff(objA, objB, WithoutFields(allFields...))
	assert.Nil(t, err)
	assert.True(t, result.Equal(), "allFields ignored should be true")

	result, err = Diff(objA, objB, WithFields(sameFields...))
	assert.Nil(t, err)
	assert.True(t, result.Equal(), "sameFields included should be true")

	result, err = Diff(objA, objB, WithoutFields(diffFields...))
	assert.Nil(t, err)
	assert.True(t, result.Equal(), "diffFields ignored should be true")

	result, err = Diff(objA, objB, WithFields(allFields...), WithoutFields(diffFields...))
	assert.Nil(t, err)
	assert.True(t, result.Equal(), "allFields included and diffFields ignored should be true")

	result, err = Diff(objA, objB, WithFields(allFields...), WithoutFields(allFields...))
	assert.Nil(t, err)
	assert.True(t, result.Equal(), "allFields included and allFields ignored should be true")

	for _, diffField := range diffFields {
		result, err = Diff(objA, objB, WithFields(diffField))
		assert.Nil(t, err)
		assert.False(t, result.Equal(), "diffField \"%s\" included should be false", diffField)
		assert.Equal(t, result.Diffs(), map[string][2]interface{}{
			diffField: diffValues[diffField],
		}, "diffField \"%s\" included should be equal", diffField)
	}
}
