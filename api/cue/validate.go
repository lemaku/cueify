package cue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"fmt"
	"reflect"
	"strings"
)

type Errors struct {
	Self   []string            `json:"self"`
	Others map[string][]string `json:"others"`
}

type OtherError struct {
	Path   []string `json:"path"`
	Errors []string `json:"errors"`
}

type ValidationResult struct {
	Valid  bool   `json:"valid"`
	Errors Errors `json:"errors"`
}

func Validate(path []string, json string, raw string) ValidationResult {
	json = "#export & " + json

	context := cuecontext.New()
	schema := context.CompileString(raw)

	value := context.CompileString(json, cue.Scope(schema))

	err := value.Validate(
		cue.Concrete(false),
		cue.Optional(true))

	if err != nil {
		errs := Errors{
			Self:   make([]string, 0),
			Others: make(map[string][]string),
		}

		for _, e := range errors.Errors(err) {
			format, args := e.Msg()
			msg := fmt.Sprintf(format, args...)
			errPath := strings.Join(e.Path(), ".")

			if reflect.DeepEqual(e.Path(), path) {
				errs.Self = append(errs.Self, msg)
			} else if len(errs.Others[errPath]) > 0 {
				errs.Others[errPath] = append(errs.Others[errPath], msg)
			} else {
				errs.Others[errPath] = []string{msg}
			}
		}

		return ValidationResult{Valid: false, Errors: errs}
	}

	return ValidationResult{Valid: true}
}
