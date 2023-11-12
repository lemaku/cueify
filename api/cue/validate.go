package cue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"fmt"
	"reflect"
)

type ValidationError struct {
	Errors []string
}

func Validate(path []string, json string) (bool, *ValidationError) {
	json = "universities: #universities\n\n" + json

	context := cuecontext.New()
	schema := context.CompileString(schemaString)

	value := context.CompileString(json, cue.Scope(schema))

	err := value.Validate(cue.Attributes(true),
		cue.Concrete(false),
		cue.Definitions(false),
		cue.DisallowCycles(false),
		cue.Docs(false),
		cue.Hidden(false),
		cue.Optional(false))

	if err != nil {
		var validationErrors []string

		tmp := errors.Errors(err)

		for _, e := range tmp {
			if reflect.DeepEqual(e.Path(), path) {
				format, args := e.Msg()
				validationErrors = append(validationErrors, fmt.Sprintf(format, args...))
			}
		}

		if len(validationErrors) > 0 {
			return false, &ValidationError{validationErrors}
		}

	}

	return true, nil
}
