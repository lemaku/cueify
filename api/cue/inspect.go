package cue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"fmt"
	"regexp"
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

type Property struct {
	Path  []string `json:"path"`
	Type  Kind     `json:"type"`
	Index int      `json:"index"`
}

type InspectionResult struct {
	Type       Kind       `json:"type"`
	Properties []Property `json:"properties"`
}

type Kind string

const (
	String  Kind = "string"
	Number  Kind = "number"
	Bool    Kind = "bool"
	List    Kind = "list"
	Complex Kind = "complex"
)

func Inspect(path []string, json string) InspectionResult {
	json = "universities: #universities\n\n" + json

	context := cuecontext.New()
	schema := context.CompileString(schemaString)

	value := context.CompileString(json, cue.Scope(schema))
	value = value.LookupPath(toCuePath(path))

	var iter *cue.Iterator
	fmt.Println(value.Kind())
	if value.IncompleteKind() == cue.StructKind {
		iter, _ = value.Fields()
	} else if value.IncompleteKind() == cue.ListKind {
		tmp, _ := value.List()
		iter = &tmp
	} else {
		return InspectionResult{}
	}
	var properties []Property

	i := 0
	for iter.Next() {
		properties = append(properties, Property{Path: toStringPath(iter.Value().Path()), Type: getKind(iter.Value()), Index: i})
		i++
	}
	return InspectionResult{Type: getKind(value), Properties: properties}
}

func getKind(value cue.Value) Kind {
	switch value.IncompleteKind() {
	case cue.StringKind:
		return String
	case cue.BoolKind:
		return Bool
	case cue.IntKind:
		return Number
	case cue.ListKind:
		return List
	case cue.StructKind:
		return Complex
	default:
		panic("Not implemented")
	}
}

func toStringPath(path cue.Path) []string {
	selectors := path.Selectors()
	p := make([]string, len(selectors))
	for i, selector := range selectors {
		p[i] = selector.String()
	}
	return p
}

func toCuePath(path []string) cue.Path {
	s := ""
	for i, p := range path {
		if i == 0 {
			s += p
		} else {
			if match, _ := regexp.MatchString("[0-9]+", p); match {
				s += "[" + p + "]"
			} else {
				s += "." + p
			}
		}
	}
	return cue.ParsePath(s)
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

	fmt.Println(v)

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
		// not final or concrete
		cue.Concrete(true),
		// check minimally
		cue.Definitions(true),
		cue.Hidden(true),
		cue.Optional(true),
	)
}
