package cue

import (
	"encoding/json"
	"testing"
)

const summarizeTestSchema = `
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

const summarizeTestStudentSchema = `
#student: {
	matNr:  string & =~"^[0-9]{8}$"
	name:   string
	active: *true | bool
    if active {
        semester: int
    }
}

#export: [...#student]
`

func TestSummarizeEmpty(t *testing.T) {
	result := Summarize("{}", summarizeTestSchema)

	if result.Valid {
		t.Fatalf("Result should not have been valid but was")
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic("WasmAPI.Summarize: couldn't serialize result")
	}

	t.Log(string(jsonResult))

}

func TestSummarizeTopLevelArray(t *testing.T) {
	// TODO snapshot tests
	if Summarize("[]", summarizeTestStudentSchema).Valid == false {
		t.Fatalf("Result should have been valid but was not")
	}

	if Summarize("[{}]", summarizeTestStudentSchema).Valid == true {
		t.Fatalf("Result should not have been valid but was")
	}

	if Summarize(`[{ name: "Max Mustermann"}]`, summarizeTestStudentSchema).Valid == true {
		t.Fatalf("Result should not have been valid but was")
	}

	result := Summarize(`[{ name: "Max Mustermann"}]`, summarizeTestStudentSchema)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic("WasmAPI.Summarize: couldn't serialize result")
	}

	t.Log(string(jsonResult))

}

func TestPartialExport(t *testing.T) {

	const allKindsSchema = `
#export: {
	a: null,
	b: string,
	c: float,
	d: int,
	e: bytes,
	f?: { a: string },
	g: [...string],
}`

	// TODO make cue accept JSON encoded data
	result := Summarize(`{ a: null, b: "abc", c: 1.3, d: 1, e: "ADA=" }`, allKindsSchema)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic("WasmAPI.Summarize: couldn't serialize result")
	}

	t.Log(string(jsonResult))

}
