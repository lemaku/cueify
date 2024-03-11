package cue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Property struct {
	Path     []string `json:"path"`
	Type     []Kind   `json:"type"`
	Optional bool     `json:"optional"`
}

type InspectionResult struct {
	Type       []Kind     `json:"type"`
	Of         []Kind     `json:"of"`
	Properties []Property `json:"properties"`
}

type Kind string

const (
	Bottom Kind = "bottom"
	Null   Kind = "null"
	String Kind = "string"
	Bytes  Kind = "bytes"
	Int    Kind = "int"
	Float  Kind = "float"
	Bool   Kind = "bool"
	List   Kind = "list"
	Struct Kind = "struct"
)

func Inspect(path []string, json string, raw string) InspectionResult {
	cueVal := "#export & " + json

	context := cuecontext.New()
	schema := context.CompileString(raw)

	value := context.CompileString(cueVal, cue.Scope(schema))
	value = value.LookupPath(toCuePath(path))

	if value.IncompleteKind() == cue.StructKind {
		var properties []Property
		op, values := value.Eval().Expr()
		switch op {
		case cue.OrOp:
			iter, _ := values[0].Fields(cue.Optional(true))
			properties = getProperties(iter)
			iter, _ = values[1].Fields(cue.Optional(true))
			properties = mergeProperties(append(properties, getProperties(iter)...))
			break
		case cue.NoOp:
			iter, _ := values[0].Fields(cue.Optional(true))
			properties = getProperties(iter)
			break
		default:
			panic("WasmAPI.Inspect: Given configuration was too complex - at most one disjunction within the same type is allowed")
		}
		return InspectionResult{Type: getKind(value), Properties: properties}
	} else if value.IncompleteKind() == cue.ListKind {
		var of []Kind
		op, values := value.Expr()
		if isPropertySet(path, json) {
			// Because #export & {} is used, if an input is given, the value is a conjunction between the definition and the input
			if op != cue.AndOp {
				panic("WasmAPI.Inspect: Unexpected error")
			}

			// If input is empty, use the definition to find out the types it could be
			if inputIter, _ := values[1].List(); !inputIter.Next() {
				op, values = values[0].Expr()
			} else {
				// If input is not empty, take the first value and return its type
				of = getKind(inputIter.Value())
			}
		}
		switch op {
		case cue.OrOp:
			of = append(
				getKind(values[0].LookupPath(cue.MakePath(cue.AnyIndex))),
				getKind(values[1].LookupPath(cue.MakePath(cue.AnyIndex)))...)
			break
		case cue.NoOp:
			of = getKind(values[0].LookupPath(cue.MakePath(cue.AnyIndex)))
			break
		default:
			if of == nil {
				panic("WasmAPI.Inspect: Given configuration was too complex - at most one disjunction within the same type is allowed")
			}
		}

		iter, _ := value.List()
		return InspectionResult{Type: getKind(value), Of: removeDuplicate(of), Properties: getProperties(&iter)}
	}

	return InspectionResult{Type: getKind(value), Properties: make([]Property, 0)}
}

func getProperties(iter *cue.Iterator) []Property {
	var properties []Property
	for iter.Next() {
		properties = append(properties, Property{
			Path:     toStringPath(iter.Value().Path()),
			Type:     getKind(iter.Value()),
			Optional: iter.Selector().ConstraintType() == cue.OptionalConstraint})
	}
	return properties
}

func checkKind(value cue.Value, cueKind cue.Kind, cueifyKind Kind, collection []Kind) []Kind {
	if value.IncompleteKind().IsAnyOf(cueKind) || value.IncompleteKind() == cueKind {
		return append(collection, cueifyKind)
	}
	return collection
}

func getKind(value cue.Value) []Kind {
	var kind []Kind
	kind = checkKind(value, cue.NullKind, Null, kind)
	kind = checkKind(value, cue.BoolKind, Bool, kind)
	kind = checkKind(value, cue.IntKind, Int, kind)
	kind = checkKind(value, cue.FloatKind, Float, kind)
	kind = checkKind(value, cue.StringKind, String, kind)
	kind = checkKind(value, cue.ListKind, List, kind)
	kind = checkKind(value, cue.StructKind, Struct, kind)

	if len(kind) == 0 {
		panic(fmt.Errorf("WasmAPI.Inspect: Couldn't handle %v", value.IncompleteKind().TypeString()))
	}

	return removeDuplicate(kind)
}

/*
Transform CUE path into string path
*/
func toStringPath(path cue.Path) []string {
	selectors := path.Selectors()
	p := make([]string, len(selectors))
	for i, selector := range selectors {
		p[i] = selector.String()
	}
	return p
}

/*
Transform string path into CUE path
*/
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

/*
Function to check whether the path is already set or not.
*/
func isPropertySet(path []string, raw string) bool {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(raw), &data)
	if err != nil {
		panic("WasmAPI.Inspect: couldn't unmarshall JSON value")
	}

	currentData := data
	for _, key := range path {
		value, ok := currentData[key]
		if !ok {
			return false
		}

		switch v := value.(type) {
		case map[string]interface{}:
			currentData = v
		case []interface{}:
			currentData = make(map[string]interface{})
			for i, item := range v {
				currentData[fmt.Sprintf("%d", i)] = item
			}
		}
	}

	return currentData != nil
}

/*
Remove duplicates from a slice
*/
func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

/*
Merge properties with the same path. When disjunctions between to structs are used and both structs have
a property with the same path, this function merges them into one, using the types of both.
*/
func mergeProperties(properties []Property) []Property {
	pathMap := make(map[string][]Property)

	for _, prop := range properties {
		pathStr := strings.Join(prop.Path, ".")
		pathMap[pathStr] = append(pathMap[pathStr], prop)
	}

	var mergedProperties []Property

	for _, props := range pathMap {
		if len(props) == 1 {
			mergedProperties = append(mergedProperties, props[0])
		} else {
			mergedProp := Property{
				Path:     props[0].Path,
				Type:     make([]Kind, 0),
				Optional: true,
			}
			for _, prop := range props {
				mergedProp.Type = append(mergedProp.Type, prop.Type...)
				mergedProp.Optional = mergedProp.Optional && prop.Optional
			}

			mergedProp.Type = removeDuplicate(mergedProp.Type)
			mergedProperties = append(mergedProperties, mergedProp)
		}
	}

	return mergedProperties
}
