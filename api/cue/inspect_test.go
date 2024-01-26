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
	expectedType := []Kind{Complex}
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
	expectedType := []Kind{Complex}
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
	expectedType := []Kind{Complex}
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
		b: number | bool
		c: { d: string } | string
		e: float | _
		f: null
		g: >2
		h: string
		i: "Hello, \(h)!"
		j: bytes
	}
}`

	test := func(path []string, expectedType []Kind) {
		result := Inspect(path, json, schema)
		if !reflect.DeepEqual(result.Type, expectedType) {
			t.Fatalf("Type for %v should have been %v but was %v", strings.Join(path, "."), expectedType, result.Type)
		}
	}
	test([]string{"a"}, []Kind{Complex})
	test([]string{"a", "b"}, []Kind{Bool, Number})
	test([]string{"a", "c"}, []Kind{String, Complex})
	test([]string{"a", "e"}, []Kind{Null, Bool, Number, String, Bytes, List, Complex})
	test([]string{"a", "f"}, []Kind{Null})
	test([]string{"a", "g"}, []Kind{Number})
	test([]string{"a", "h"}, []Kind{String})
	test([]string{"a", "i"}, []Kind{Bottom})
	test([]string{"a", "j"}, []Kind{Bytes})
}
