package cue

import (
	"fmt"
	"sort"
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

func Summarize(json string, raw string) SummarizeResult {
	json = "#export & " + json

	context := cuecontext.New()
	schema := context.CompileString(raw)

	value := context.CompileString(json, cue.Scope(schema))

	// Don't allow any errors other than errors due to values not being concrete
	if len(errors.Errors(value.Validate(cue.Concrete(false)))) > 0 {
		panic("WasmAPI.Summarize: Unexpected errors - current value has other errors than IncompleteError")
	}

	err := value.Validate(
		cue.Concrete(true),
		cue.Optional(true),
	)

	if err == nil {
		return SummarizeResult{Value: partialExport(value), Valid: true, Errors: []Error{}}
	}

	// Map path of errors to the messages
	m := make(map[string][]string)
	for _, e := range errors.Errors(err) {
		path := strings.Join(e.Path(), ".")
		msgFormat, args := e.Msg()
		m[path] = append(m[path], fmt.Sprintf(msgFormat, args...))
	}

	// Analyze dependencies of the values and suggest what values to use next by using a heuristic
	paths := make([]string, len(m))
	i := 0
	for k := range m {
		paths[i] = k
		i++
	}
	depMap, err := analyzeDeps(value, paths)
	if err != nil {
		fmt.Println(err)
	} else {
		sort.Slice(paths, func(i, j int) bool {
			return depMap[paths[i]] < depMap[paths[j]]
		})
	}

	var errs []Error
	for _, key := range paths {
		errs = append(errs, Error{Path: strings.Split(key, "."), Errors: m[key]})
	}
	return SummarizeResult{Value: partialExport(value), Valid: false, Errors: errs}
}

func analyzeDeps(value cue.Value, paths []string) (map[string]int, error) {
	dependenciesCount := make(map[string]int)
	visited := make(map[string]bool)
	inProgress := make(map[string]bool)

	var count func(path string) (int, error)
	count = func(path string) (int, error) {
		curVal := value.LookupPath(toCuePath(strings.Split(path, ".")))
		if inProgress[path] {
			return 0, fmt.Errorf("WasmAPI.Summarize: Dependency analysis failed - cycle detected")
		}
		if visited[path] {
			return dependenciesCount[path], nil
		}
		inProgress[path] = true
		cnt := 0
		if !curVal.IsConcrete() {
			for _, dep := range curVal.Deps() {
				if !dep.IsConcrete() {
					depCnt, err := count(strings.Join(toStringPath(dep.Path()), "."))
					if err != nil {
						return 0, err
					}
					cnt += 1 + depCnt
				}
			}
		}
		dependenciesCount[path] = cnt
		visited[path] = true
		delete(inProgress, path)
		return cnt, nil
	}

	retVal := make(map[string]int)
	for _, p := range paths {
		cnt, err := count(p)
		if err != nil {
			return nil, err
		}
		retVal[p] = cnt
	}

	return retVal, nil
}

func partialExport(value cue.Value) interface{} {
	switch value.IncompleteKind() {
	case cue.StructKind:
		return partialStructExport(value)
	case cue.ListKind:
		return partialListExport(value)
	}

	return make(map[string]interface{})
}

func partialListExport(value cue.Value) interface{} {
	export := make([]interface{}, 0)

	if value.IncompleteKind() != cue.ListKind {
		panic("WasmAPI.Summarize: Unexpected errors - could not partially export value")
	}

	list, _ := value.List()
	for list.Next() {
		listElement := list.Value()
		if listElement.IncompleteKind() == cue.IntKind || listElement.IncompleteKind() == cue.StringKind || listElement.IncompleteKind() == cue.BoolKind {
			export = append(export, listElement.Value())
		} else if listElement.IncompleteKind() == cue.StructKind {
			export = append(export, partialStructExport(listElement))
		} else if listElement.IncompleteKind() == cue.ListKind {
			export = append(export, partialListExport(listElement))
		}
	}
	return export
}

func partialStructExport(value cue.Value) interface{} {
	export := make(map[string]interface{})

	if value.IncompleteKind() != cue.StructKind {
		panic("WasmAPI.Summarize: Unexpected errors - could not partially export value")
	}

	s, _ := value.Fields()

	for s.Next() {
		property := s.Value()
		pathSelectors := property.Path().Selectors()
		propertyName := pathSelectors[len(pathSelectors)-1].String()

		if property.IncompleteKind().IsAnyOf(cue.NumberKind | cue.StringKind | cue.BoolKind | cue.NullKind | cue.BytesKind) {
			if property.IsConcrete() {
				export[propertyName] = s.Value()
			}
			defaultVal, hasDefault := property.Default()
			if hasDefault {
				export[propertyName] = defaultVal
			}
		} else if property.IncompleteKind() == cue.StructKind {
			export[propertyName] = partialStructExport(property)
		} else if property.IncompleteKind() == cue.ListKind {
			export[propertyName] = partialListExport(property)
		}
	}
	return export
}
