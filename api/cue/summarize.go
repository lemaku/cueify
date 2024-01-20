package cue

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
)

type Error struct {
	Path   []string `json:"path"`
	Errors []string `json:"errors"`
}

type SummarizeResult struct {
	Value  interface{} `json:"value"`
	Valid  bool        `json:"valid"`
	Errors []Error     `json:"errors"`
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
		return SummarizeResult{Value: partialExport(value), Valid: true, Errors: []Error{}}
	}

	m := make(map[string][]string)
	for _, e := range errors.Errors(err) {
		path := strings.Join(e.Path(), ".")
		msgFormat, args := e.Msg()
		m[path] = append(m[path], fmt.Sprintf(msgFormat, args...))
	}

	var errs []Error
	for key, value := range m {
		errs = append(errs, Error{Path: strings.Split(key, "."), Errors: value})
	}
	return SummarizeResult{Value: partialExport(value), Valid: false, Errors: errs}
}

func partialExport(value cue.Value) map[string]interface{} {
	export := make(map[string]interface{})

	switch value.IncompleteKind() {
	case cue.StructKind:
		s, _ := value.Fields(defaultOptions...)

		for s.Next() {
			property := s.Value()
			pathSelectors := property.Path().Selectors()
			propertyName := pathSelectors[len(pathSelectors)-1].String()

			if property.IncompleteKind() == cue.IntKind || property.IncompleteKind() == cue.StringKind || property.IncompleteKind() == cue.BoolKind {
				if property.IsConcrete() {
					export[propertyName] = s.Value()
				}
				defaultVal, hasDefault := property.Default()
				if hasDefault {
					export[propertyName] = defaultVal
				}

			} else if property.IncompleteKind() == cue.StructKind {
				export[propertyName] = partialExport(property)
			} else if property.IncompleteKind() == cue.ListKind {
				list, _ := property.List()

				var arr []interface{}

				for list.Next() {
					listElement := list.Value()
					if listElement.IncompleteKind() == cue.IntKind || listElement.IncompleteKind() == cue.StringKind || listElement.IncompleteKind() == cue.BoolKind {
						arr = append(arr, listElement.Value())
					} else if listElement.IncompleteKind() == cue.StructKind {
						arr = append(arr, partialExport(listElement))
					}
				}
				export[propertyName] = arr
			}
		}
	}

	return export
}
