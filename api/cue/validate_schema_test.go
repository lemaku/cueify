package cue

import (
	"testing"
)

func TestValidSchemaWithMultipleDefinitions(t *testing.T) {
	schema := `
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
		name: "Vienna University of Technology" | "University of Vienna",
		students: [...#student]
	},
	countryCode: string
}

#export: #universities
`
	result := ValidateSchema(schema)
	if !result.Valid {
		t.Fatalf("Schema should be valid but was not")
	}
}

func TestInValidSchemaWithMultipleDefinitionsWithoutExport(t *testing.T) {
	schema := `
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
		name: "Vienna University of Technology" | "University of Vienna",
		students: [...#student]
	},
	countryCode: string
}
`
	result := ValidateSchema(schema)
	if result.Valid {
		t.Fatalf("Schema should not be valid but was")
	}
}

func TestValidSchema(t *testing.T) {
	schema := `
#export: {
	abc: string
	def: int
}
`
	result := ValidateSchema(schema)
	if !result.Valid {
		t.Fatalf("Schema should be valid but was not")
	}
}

func TestInvalidSchemaWithSyntaxError(t *testing.T) {
	schema := `
#export: {
	abc: stringi
	def: int
}
`
	result := ValidateSchema(schema)
	if result.Valid {
		t.Fatalf("Schema should not valid but was")
	}
}
