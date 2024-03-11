package cue

import (
	"cuelang.org/go/cue/cuecontext"
	"encoding/json"
	"github.com/gkampitakis/go-snaps/snaps"
	"testing"
)

func matchSummarizeSnapshot(t *testing.T, value interface{}) {
	jsonResult, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		panic("Error: couldn't serialize result to use for the snapshot")
	}
	snaps.MatchSnapshot(t, string(jsonResult))
}

func TestSummarizeEmpty(t *testing.T) {
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
	matchSummarizeSnapshot(t, Summarize("{}", schema))
}

func TestSummarizeTopLevelArray(t *testing.T) {
	const schema = `
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
	matchSummarizeSnapshot(t, Summarize("[]", schema))
	matchSummarizeSnapshot(t, Summarize("[{}]", schema))
	matchSummarizeSnapshot(t, Summarize(`[{ name: "Max Mustermann"}]`, schema))
	matchSummarizeSnapshot(t, Summarize(`[{ name: "Max Mustermann", matNr: "11111111"}]`, schema))
	matchSummarizeSnapshot(t, Summarize(`[{ name: "Max Mustermann", matNr: "11111111", semester: 5}]`, schema))
}

func TestDependencyHeuristic(t *testing.T) {
	raw := `
{
	z: y
	y: x
	x: int
} & {}`
	context := cuecontext.New()
	value := context.CompileString(raw)

	depMap, _ := analyzeDeps(value, []string{"x", "y", "z"})
	matchSummarizeSnapshot(t, depMap)
}

func TestDependencyHeuristicWithConcreteValues(t *testing.T) {
	raw := `
{
	z: y
	y: x
	x: t
	t: int
} & { x: 4 }`
	context := cuecontext.New()
	value := context.CompileString(raw)

	depMap, _ := analyzeDeps(value, []string{"t", "x", "y", "z"})
	matchSummarizeSnapshot(t, depMap)
}

func TestDependencyHeuristicWithCycles(t *testing.T) {
	raw := `
{
	x: y
	y: x
} & {}`
	context := cuecontext.New()
	value := context.CompileString(raw)

	_, err := analyzeDeps(value, []string{"x", "y", "z"})
	snaps.MatchSnapshot(t, err.Error())
}

func TestPartialExportShouldNotIncludeOptionalProperties(t *testing.T) {
	// "a" should not be included since it is optional
	matchSummarizeSnapshot(t, Summarize(`{ c: 1.3 }`, `
#export: {
	a?: null,
	b: string,
	c: float
}`).Value)
	// now "a" is not optional anymore and has concrete value, so should be included
	matchSummarizeSnapshot(t, Summarize(`{ c: 1.3 }`, `
#export: {
	a: null,
	b: string,
	c: float
}`).Value)
}

func TestPartialExportShouldIncludeDefaults(t *testing.T) {
	matchSummarizeSnapshot(t, Summarize(`{}`, `
#export: {
	a: string | *"this is the default",
}`).Value)
}

func TestPartialExportShouldHandleStructsAndListsCorrectly(t *testing.T) {
	const schema = `
#export: {
	a: { b: string },
	c?: { d: string },
	e: [...string],
	f?: [...string],
}`

	matchSummarizeSnapshot(t, Summarize(`{}`, schema).Value)
	matchSummarizeSnapshot(t, Summarize(`{c: {d: "hello"}, f: [",", "world"]}`, schema).Value)
}
