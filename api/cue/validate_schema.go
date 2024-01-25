package cue

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
)

type SchemaValidationResult struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

func ValidateSchema(raw string) SchemaValidationResult {
	context := cuecontext.New()

	// Check input for compilation errors
	schema := context.CompileString(raw)
	if schema.Err() != nil {
		return SchemaValidationResult{Valid: false, Error: errors.Details(schema.Err(), nil)}
	}

	// Check for existence of #export definition
	export := schema.LookupPath(cue.ParsePath("#export"))
	if export.Err() != nil {
		return SchemaValidationResult{Valid: false, Error: errors.Details(export.Err(), nil)}
	}

	// No validation error and #export exists
	return SchemaValidationResult{Valid: true, Error: ""}
}
