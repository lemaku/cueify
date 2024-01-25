package cue

import (
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
	expectedType := Complex
	expectedProperties := 2

	if result.Type != expectedType {
		t.Fatalf("Type should have been %v but was %v", expectedType, result.Type)
	}

	if len(result.Properties) != expectedProperties {
		t.Fatalf("There should have been %d properties but there were only %d", expectedProperties, len(result.Properties))
	}
}

func TestInspectMiddle(t *testing.T) {
	result := Inspect([]string{"tuwien"}, "{ tuwien: { students: null } }", schema)
	expectedType := Complex
	expectedProperties := 2

	if result.Type != expectedType {
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
	expectedType := Complex
	expectedProperties := 4

	if result.Type != expectedType {
		t.Fatalf("Type should have been %v but was %v", expectedType, result.Type)
	}

	if len(result.Properties) != expectedProperties {
		t.Fatalf("There should have been %d properties but there were only %d", expectedProperties, len(result.Properties))
	}
}
