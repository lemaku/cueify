package cue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"fmt"
	"reflect"
)

type ValidationError struct {
	Errors []string `json:"errors"`
}

type ValidationResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

func Validate(path []string, json string, raw string) ValidationResult {
	json = "#export & " + json

	context := cuecontext.New()
	schema := context.CompileString(raw)

	value := context.CompileString(json, cue.Scope(schema))

	err := value.Validate(
		cue.Attributes(true),
		cue.Concrete(true),
		cue.Definitions(false),
		cue.DisallowCycles(false),
		cue.Docs(false),
		cue.Hidden(false),
		cue.Optional(false))

	if err != nil {
		var validationErrors []string

		errs := errors.Errors(err)

		for _, e := range errs {
			if reflect.DeepEqual(e.Path(), path) {
				format, args := e.Msg()
				validationErrors = append(validationErrors, fmt.Sprintf(format, args...))
			}
		}

		if len(validationErrors) > 0 {
			return ValidationResult{Valid: false, Errors: validationErrors}
		}
	}

	return ValidationResult{Valid: true, Errors: make([]string, 0)}
}
