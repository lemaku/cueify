package cue

import (
	"reflect"
	"strings"
	"testing"
)

const schema = `
#student: {
	matNr:  string & =~"^[0-9]{8}$"
	name:   string
	active: *true | bool
    if active {
        semester: int
    }
}

#universities: {
	tuwien: {
		name: "Vienna University of Technology",
		students: [...#student]
	},
	countryCode: string
}

#export: #universities
`

func TestInspectRoot(t *testing.T) {
	result := Inspect([]string{}, "{}", schema)
	expectedType := []Kind{Struct}
	expectedProperties := 2

	if !reflect.DeepEqual(result.Type, expectedType) {
		t.Fatalf("Type should have been %v but was %v", expectedType, result.Type)
	}

	if len(result.Properties) != expectedProperties {
		t.Fatalf("There should have been %d properties but there were only %d", expectedProperties, len(result.Properties))
	}
}

func TestInspectMiddle(t *testing.T) {
	result := Inspect([]string{"tuwien"}, "{ tuwien: { students: [] } }", schema)
	expectedType := []Kind{Struct}
	expectedProperties := 2

	if !reflect.DeepEqual(result.Type, expectedType) {
		t.Fatalf("Type should have been %v but was %v", expectedType, result.Type)
	}

	if len(result.Properties) != expectedProperties {
		t.Fatalf("There should have been %d properties but there were only %d", expectedProperties, len(result.Properties))
	}
}

func TestInspectDeep(t *testing.T) {
	result := Inspect([]string{"tuwien", "students", "0"}, `{
	tuwien: {
		name: "Vienna University of Technology",
		students: [
			{
				name: "Test"
			}
		]
	}
}`, schema)
	expectedType := []Kind{Struct}
	expectedProperties := 4

	if !reflect.DeepEqual(result.Type, expectedType) {
		t.Fatalf("Type should have been %v but was %v", expectedType, result.Type)
	}

	if len(result.Properties) != expectedProperties {
		t.Fatalf("There should have been %d properties but there were only %d", expectedProperties, len(result.Properties))
	}
}

func TestInspectComplexSchema(t *testing.T) {
	json := `{
	a: {
	}
}`
	schema := `#export: {
	a: {
		a?: string
		b: number | bool
		c: { d: string } | string
		e: float | _
		f: null
		g: >2
		h: string
		i: "Hello, \(h)!"
		j: bytes,
		k: [...string] | [...int]
		l: { a: string } | { b: string }
	}
}`

	test := func(path []string, expectedType []Kind) {
		result := Inspect(path, json, schema)
		if !reflect.DeepEqual(result.Type, expectedType) {
			t.Fatalf("Type for %v should have been %v but was %v", strings.Join(path, "."), expectedType, result.Type)
		}
	}
	// TODO: use snapshot tests
	test([]string{"a"}, []Kind{Struct})
	test([]string{"a", "b"}, []Kind{Bool, Int, Float})
	test([]string{"a", "c"}, []Kind{String, Struct})
	test([]string{"a", "e"}, []Kind{Null, Bool, Int, Float, String, Bytes, List, Struct})
	test([]string{"a", "f"}, []Kind{Null})
	test([]string{"a", "g"}, []Kind{Int, Float})
	test([]string{"a", "h"}, []Kind{String})
	test([]string{"a", "i"}, []Kind{Bottom})
	test([]string{"a", "j"}, []Kind{Bytes})
	test([]string{"a", "k"}, []Kind{List})
	test([]string{"a", "l"}, []Kind{Struct})
}

func TestInspectSuperComplexSchema(t *testing.T) {
	json := `{
	a: {
		b: [],
		c: {},
		d: [],
		e: {}
	}
}`
	schema := `#export: {
	a: {
		b: [...string] | { b: string }
		c: [...string] | { b: string }
		d: [...string] | [...int]
		e: { a: string } | { b: string }
	}
}`

	test := func(path []string, expectedType []Kind) {
		result := Inspect(path, json, schema)
		if !reflect.DeepEqual(result.Type, expectedType) {
			t.Fatalf("Type for %v should have been %v but was %v", strings.Join(path, "."), expectedType, result.Type)
		}
	}
	test([]string{"a", "b"}, []Kind{List}) // TODO: Should say type of list
	test([]string{"a", "c"}, []Kind{Struct})
	test([]string{"a", "d"}, []Kind{List})   // TODO: Should give option of type
	test([]string{"a", "e"}, []Kind{Struct}) // TODO: Should give both props as options
}
