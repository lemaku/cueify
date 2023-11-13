package cue

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
)

const schemaString = `
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
		name: "Vienna University of Technology" | "WU Wien",
		students: [...#student]
	},
	countryCode: string
}
`

const correctVal = `
universities: #universities

universities: {
    tuwien: {
		name: "Vienna University of Technology",
        students: [
            {
                matNr: "12119877"
                name: "Leon K"
                semester: 5
            }
        ]
    }
}
`

const MissingPropVal = `
universities: {
    tuwien: {
        students: [
            {
				matNr: "12119877"
            }
        ]
    }
}
`

const emptyVal = `
universities: #universities

universities: {}
`

type Error struct {
	Path   []string `json:"path"`
	Errors []string `json:"errors"`
}

type SummarizeResult struct {
	Valid  bool    `json:"valid"`
	Errors []Error `json:"errors"`
}

func Summarize(json string) SummarizeResult {
	json = "universities: #universities\n\n" + json

	context := cuecontext.New()
	schema := context.CompileString(schemaString)

	value := context.CompileString(json, cue.Scope(schema))

	err := value.Validate(
		cue.Concrete(true),
		cue.Definitions(true),
		cue.Hidden(false),
		cue.Optional(true),
	)

	if err == nil {
		return SummarizeResult{Valid: true, Errors: []Error{}}
	}

	m := make(map[string][]string)
	for _, e := range errors.Errors(err) {
		path := strings.Join(e.Path(), ".")
		format, args := e.Msg()
		m[path] = append(m[path], fmt.Sprintf(format, args...))
	}

	var errs []Error
	for key, value := range m {
		errs = append(errs, Error{Path: strings.Split(key, "."), Errors: value})
	}
	return SummarizeResult{Valid: false, Errors: errs}
}

func Do() {
	var (
		c      *cue.Context
		schema cue.Value
		v      cue.Value
	)

	c = cuecontext.New()
	schema = c.CompileString(schemaString)
	v = c.CompileString(MissingPropVal, cue.Scope(schema))

	printErr("test", schema.Subsume(v))

	printErr("validate", loose(v))

	// tmp := v.LookupPath(cue.ParsePath("universities.tuwien.students[0]"))

	//fmt.Println(tmp)

	// Walk(v, customOptions...)
	//
	//// print the value
	//fmt.Println(v)
	//printFields(v.Fields())
	//
	//printErr("test", schema.Subsume(v))
	//
	// printErr("loose error", loose(v))
}

// Walk through object and the children of all non-trivial fields
func Walk(v cue.Value, options ...cue.Option) {

	if v.IncompleteKind() == cue.StructKind || v.IncompleteKind() == cue.ListKind {
		fmt.Println(">", v.Path())
		iter, _ := v.Fields(options...)
		for iter.Next() {
			fmt.Println(iter.Selector(), ":", iter.Value().IncompleteKind())
		}

		fmt.Println()
	}

	// possibly recurse
	switch v.IncompleteKind() {
	case cue.StructKind:
		if options == nil {
			options = defaultOptions
		}
		s, _ := v.Fields(options...)

		for s.Next() {
			tmp := s.Value()
			Walk(tmp, options...)
		}

	case cue.ListKind:
		l, _ := v.List()
		for l.Next() {
			tmp := l.Value()
			Walk(tmp.Value(), options...)
		}
	}
}

// Cue's default
var defaultOptions = []cue.Option{
	cue.Attributes(true),
	cue.Concrete(false),
	cue.Definitions(false),
	cue.DisallowCycles(false),
	cue.Docs(false),
	cue.Hidden(false),
	cue.Optional(false),
}

// Our custom options
var customOptions = []cue.Option{
	cue.Concrete(false),
	cue.Definitions(true),
	cue.Hidden(true),
	cue.Optional(true),
}

func printFields(iter *cue.Iterator, err error) {
	for iter.Next() {
		fmt.Printf("%v: %v\n", iter.Selector(), iter.Value())
	}
	fmt.Println()
}

func printErr(prefix string, err error) {
	if err != nil {
		msg := errors.Details(err, nil)
		fmt.Printf("%s:\n%s\n", prefix, msg)
	}
}

func loose(v cue.Value) error {
	return v.Validate(
		cue.Concrete(true),
		cue.Definitions(true),
		cue.Hidden(false),
		cue.Optional(true),
	)
}
