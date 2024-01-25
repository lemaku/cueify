package cue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"regexp"
)

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

func Inspect(path []string, json string, raw string) InspectionResult {
	json = "#export & " + json

	context := cuecontext.New()
	schema := context.CompileString(raw)

	value := context.CompileString(json, cue.Scope(schema))
	value = value.LookupPath(toCuePath(path))

	var iter *cue.Iterator
	if value.IncompleteKind() == cue.StructKind {
		iter, _ = value.Fields(cue.Optional(true))
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
