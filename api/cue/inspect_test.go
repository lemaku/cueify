package cue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"encoding/json"
	"fmt"
	"github.com/gkampitakis/go-snaps/snaps"
	"testing"
)

func matchInspectSnapshot(t *testing.T, value interface{}) {
	jsonResult, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		panic("Error: couldn't serialize result to use for the snapshot")
	}
	snaps.MatchSnapshot(t, string(jsonResult))
}

func TestUniversitiesSchema(t *testing.T) {
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
		name: "Vienna University of Technology",
		students: [...#student]
	},
	countryCode?: string
}

#export: #universities
`
	matchInspectSnapshot(t, Inspect([]string{}, "{}", schema))
	matchInspectSnapshot(t, Inspect([]string{"tuwien"}, `{ "tuwien": { "students": [] } }`, schema))
	matchInspectSnapshot(t, Inspect([]string{"tuwien", "students"}, `{
		"tuwien": {
			"name": "Vienna University of Technology",
			"students": [
				{
					"name": "Test"
				}
			]
		}
	}`, schema))
	matchInspectSnapshot(t, Inspect([]string{"tuwien", "students", "0"}, `{
		"tuwien": {
			"name": "Vienna University of Technology",
			"students": [
				{
					"name": "Test"
				}
			]
		}
	}`, schema))
}

func TestInspectComplexSchema(t *testing.T) {
	raw := `{
	a: {
		k1: []
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
		k1: [...string]
		k2: [...string]
		l: [...string] | [...int]
		m: { a: string } | { b: string }
	}
}`
	matchInspectSnapshot(t, Inspect([]string{"a"}, raw, schema))
	matchInspectSnapshot(t, Inspect([]string{"a", "k1"}, raw, schema))
	matchInspectSnapshot(t, Inspect([]string{"a", "k2"}, raw, schema))
}

func TestInspectSuperComplexSchema(t *testing.T) {
	raw := `{
	"b": [],
	"c": {},
	"d": [],
	"e": {}
}`
	schema := `#export: {
	b?: [...string] | { b: string }
	b2: [...string] | { b: string }
	c: [...string] | { b: string }
	c2: [...string] | { b: string }
	d: [...string] | [...int|bool]
	d2: [...string] | [...int|bool]
	e: { a: string } | { b: string }
	e2: { a: string } | { b: string }
	f: { b?: string } | { b: int }
}`
	matchInspectSnapshot(t, Inspect([]string{}, raw, schema))
	matchInspectSnapshot(t, Inspect([]string{"b"}, raw, schema)) // TODO this is still a problem
	matchInspectSnapshot(t, Inspect([]string{"b2"}, raw, schema))
	matchInspectSnapshot(t, Inspect([]string{"d"}, raw, schema))  // TODO still issue because .Eval can't be used
	matchInspectSnapshot(t, Inspect([]string{"d2"}, raw, schema)) // TODO still issue because .Eval can't be used
	matchInspectSnapshot(t, Inspect([]string{"e"}, raw, schema))
	matchInspectSnapshot(t, Inspect([]string{"e2"}, raw, schema))
	matchInspectSnapshot(t, Inspect([]string{"f"}, raw, schema))
}

func TestLookupPath(t *testing.T) {
	context := cuecontext.New()

	// I would expect this to print "string"
	value := context.CompileString(`["abc"]`)
	fmt.Println(value.LookupPath(cue.MakePath(cue.AnyIndex)))
	fmt.Println(value.Eval().LookupPath(cue.MakePath(cue.AnyIndex)))

	// I would expect this to print "string" just like if "value" was simply [...string]
	value = context.CompileString(`([...string] | { a: string }) & []`)
	fmt.Println(value.Eval().LookupPath(cue.MakePath(cue.AnyIndex)))

	value = context.CompileString(`([...string] | [...int]) & []`)
	op, values := value.Eval().Expr()
	// op is OrOp
	fmt.Println(op)
	// => I would expect this to print "string"
	fmt.Println(values[0].Eval().LookupPath(cue.MakePath(cue.AnyIndex)))
	// => I would expect this to print "int"
	fmt.Println(values[1].Eval().LookupPath(cue.MakePath(cue.AnyIndex)))

	// I would expect this to print "string"
	value = context.CompileString(`([...string] | { a: string }) & ["abc"]`)
	fmt.Println(value.Eval().LookupPath(cue.MakePath(cue.AnyIndex)))

	// U would expect this to print "string"
	value = context.CompileString(`{ a: [...[...string]]} & {a: [[]]}`).LookupPath(cue.ParsePath("a[0]"))
	fmt.Println(value.Eval().LookupPath(cue.MakePath(cue.AnyIndex)))
}
