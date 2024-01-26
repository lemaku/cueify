package cue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"fmt"
	"regexp"
	"slices"
)

type Property struct {
	Path  []string `json:"path"`
	Type  []Kind   `json:"type"`
	Index int      `json:"index"`
}

type InspectionResult struct {
	Type       []Kind     `json:"type"`
	Properties []Property `json:"properties"`
}

type Kind string

const (
	Bottom  Kind = "bottom"
	Null    Kind = "null"
	String  Kind = "string"
	Bytes   Kind = "bytes"
	Number  Kind = "number"
	Bool    Kind = "bool"
	List    Kind = "list"
	Complex Kind = "complex"
	Any     Kind = "any"
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
		return InspectionResult{Type: getKind(value), Properties: make([]Property, 0)}
	}
	var properties []Property

	i := 0
	for iter.Next() {
		properties = append(properties, Property{Path: toStringPath(iter.Value().Path()), Type: getKind(iter.Value()), Index: i})
		i++
	}
	return InspectionResult{Type: getKind(value), Properties: properties}
}

func checkKind(value cue.Value, cueKind cue.Kind, cueifyKind Kind, collection []Kind) []Kind {
	if value.IncompleteKind().IsAnyOf(cueKind) || value.IncompleteKind() == cueKind {
		return append(collection, cueifyKind)
	}
	return collection
}

func getKind(value cue.Value) []Kind {
	var kind []Kind
	kind = checkKind(value, cue.BottomKind, Bottom, kind)
	kind = checkKind(value, cue.NullKind, Null, kind)
	kind = checkKind(value, cue.BoolKind, Bool, kind)
	kind = checkKind(value, cue.IntKind, Number, kind)
	kind = checkKind(value, cue.FloatKind, Number, kind)
	kind = checkKind(value, cue.StringKind, String, kind)
	kind = checkKind(value, cue.BytesKind, Bytes, kind)
	kind = checkKind(value, cue.ListKind, List, kind)
	kind = checkKind(value, cue.StructKind, Complex, kind)

	if len(kind) == 0 {
		panic(fmt.Errorf("Kind %v not implemented", value.IncompleteKind().TypeString()))

	}

	// This only removes consecutive equal copies => number, number need to be right after each other!
	return slices.Compact(kind)
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
