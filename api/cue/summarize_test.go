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

func TestSummarizeEmpty(t *testing.T) {
	result := Summarize("{}", summarizeTestSchema)

	if result.Valid != false {
		t.Fatalf("Result should not have been valid but was")
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic("WasmAPI.Summarize: couldn't serialize result")
	}

	t.Log(string(jsonResult))

}
