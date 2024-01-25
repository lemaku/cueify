package main

import (
	"cueify/cue"
	"encoding/json"
	"fmt"
	"syscall/js"
)

func main() {
	api := js.Global().Get("WasmAPI")
	api.Set("_validate", js.FuncOf(validateWasm))
	api.Set("_validateSchema", js.FuncOf(validateSchemaWasm))
	api.Set("_summarize", js.FuncOf(summarizeWasm))
	api.Set("_inspect", js.FuncOf(inspectWasm))

	select {}
}

func inspectWasm(this js.Value, args []js.Value) interface{} {
	const expArgs = 3

	if len(args) != expArgs {
		panic(fmt.Errorf("WasmAPI.Inspect: expected %v args, got %v", expArgs, len(args)))
	}

	// Parse path parameter
	arg0 := args[0]
	if arg0.Type() != js.TypeObject {
		panic(fmt.Errorf("WasmAPI.Inspect: expected arg %v to be of type syscall/js.TypeObject, got %v", 0, arg0.Type()))
	}
	path := make([]string, arg0.Length())
	for i := 0; i < len(path); i++ {
		item := arg0.Index(i)
		if item.Type() != js.TypeString {
			panic(fmt.Errorf("WasmAPI.Inspect: expected arg %v to contain items of type syscall/js.TypeString, got %v", 0, item.Type()))
			return nil
		}
		path[i] = item.String()
	}

	// Parse json parameter
	if args[1].Type() != js.TypeString {
		panic(fmt.Errorf("WasmAPI.Inspect: expected arg %v to be of type syscall/js.TypeString, got %v", 1, args[1].Type()))
	}
	jsonInput := args[1].String()

	if args[2].Type() != js.TypeString {
		panic(fmt.Errorf("WasmAPI.Inspect: expected arg %v to be of type syscall/js.TypeString, got %v", 2, args[2].Type()))
	}
	raw := args[2].String()

	result := cue.Inspect(path, jsonInput, raw)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic("WasmAPI.Inspect: couldn't serialize result")
	}
	return string(jsonResult)
}

func summarizeWasm(this js.Value, args []js.Value) interface{} {
	const expArgs = 2

	if len(args) != expArgs {
		panic(fmt.Errorf("WasmAPI.Summarize: expected %v args, got %v", expArgs, len(args)))
	}

	// Parse json parameter
	if args[0].Type() != js.TypeString {
		panic(fmt.Errorf("WasmAPI.Summarize: expected arg %v to be of type syscall/js.TypeString, got %v", 0, args[0].Type()))
	}
	jsonInput := args[0].String()

	if args[1].Type() != js.TypeString {
		panic(fmt.Errorf("WasmAPI.Inspect: expected arg %v to be of type syscall/js.TypeString, got %v", 1, args[1].Type()))
	}
	raw := args[1].String()

	result := cue.Summarize(jsonInput, raw)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic("WasmAPI.Summarize: couldn't serialize result")
	}

	return string(jsonResult)
}

func validateWasm(this js.Value, args []js.Value) interface{} {
	const expArgs = 3

	if len(args) != expArgs {
		panic(fmt.Errorf("WasmAPI.Validate: expected %v args, got %v", expArgs, len(args)))
	}

	// Parse path parameter
	arg0 := args[0]
	if arg0.Type() != js.TypeObject {
		panic(fmt.Errorf("WasmAPI.Validate: expected arg %v to be of type syscall/js.TypeObject, got %v", 0, arg0.Type()))
	}
	path := make([]string, arg0.Length())
	for i := 0; i < len(path); i++ {
		item := arg0.Index(i)
		if item.Type() != js.TypeString {
			panic(fmt.Errorf("WasmAPI.Validate: expected arg %v to contain items of type syscall/js.TypeString, got %v", 0, item.Type()))
			return nil
		}
		path[i] = item.String()
	}

	// Parse jsonInput parameter
	if args[1].Type() != js.TypeString {
		panic(fmt.Errorf("WasmAPI.Validate: expected arg %v to be of type syscall/js.TypeObject, got %v", 1, args[1].Type()))
	}
	jsonInput := args[1].String()

	if args[2].Type() != js.TypeString {
		panic(fmt.Errorf("WasmAPI.Inspect: expected arg %v to be of type syscall/js.TypeString, got %v", 1, args[1].Type()))
	}
	raw := args[2].String()

	result := cue.Validate(path, jsonInput, raw)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic("WasmAPI.Validate: couldn't serialize result")
	}

	return string(jsonResult)
}

func validateSchemaWasm(this js.Value, args []js.Value) interface{} {
	const expArgs = 1

	if len(args) != expArgs {
		panic(fmt.Errorf("WasmAPI.ValidateSchema: expected %v args, got %v", expArgs, len(args)))
	}

	// Parse schema parameter
	schema := args[0]
	if schema.Type() != js.TypeString {
		panic(fmt.Errorf("WasmAPI.ValidateSchema: expected arg %v to be of type syscall/js.TypeString, got %v", 0, schema.Type()))
	}

	result := cue.ValidateSchema(schema.String())
	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic("WasmAPI.ValidateSchema: couldn't serialize result")
	}

	return string(jsonResult)
}
