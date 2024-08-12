# cueify

This project was developed in the context of an academic project, creating a prototype for automatically generating web forms for CUE configurations. Please note that it is not a feature-complete and robust application. The respective thesis was submitted at the Technical University of Vienna and can be requested by contacting @lemaku.

## Interface to CUE

The `api` folder contains the GO-based API for dealing with CUE values. For the usage of the CUE Go-API, please clone the `cue` repository, apply the following `git diff`, and point to the repository in `api/go.mod` at `replace cuelang.org/go => /path/to/cue`:
```
--- a/cue/types.go
+++ b/cue/types.go
@@ -33,6 +33,7 @@ import (
        "cuelang.org/go/internal/core/adt"
        "cuelang.org/go/internal/core/compile"
        "cuelang.org/go/internal/core/convert"
+       "cuelang.org/go/internal/core/dep"
        "cuelang.org/go/internal/core/eval"
        "cuelang.org/go/internal/core/export"
        "cuelang.org/go/internal/core/runtime"
@@ -313,6 +314,21 @@ func marshalList(l *Iterator) (b []byte, err errors.Error) {
        return b, nil
 }

+func (v Value) Deps() (a []Value) {
+
+       cfg := &dep.Config{
+               Descend: true,
+       }
+       ctx := v.ctx()
+       dep.Visit(cfg, ctx, v.v, func(d dep.Dependency) error {
+               node := d.Node
+
+               a = append(a, makeValue(v.idx, node, nil))
+               return nil
+       })
+       return a
+}
+
 func (v Value) getNum(k adt.Kind) (*adt.Num, errors.Error) {
        v, _ = v.Default()
        ctx := v.ctx()
```
Use `build.sh` to compile the API into WASM; the compiled API (`main.wasm`) will be picked up and loaded when the user first visits the application. Please note that the test files where mostly used for regression tests, experimenting, and documentation, rather than for hardening the API.

## Webapp

The heart of the project is a web app developed in VueJS. Run `npm run dev` to run the application locally.

## Known issues
* Nested arrays may not work (`Inspect` function at `api/inspect.go` is affected) due to an issue where the usage of `Eval()` leads to `LookupPath(cue.MakePath(cue.AnyIndex))` not returning the expected type of the array (at the time of the development of this project). The function `TestLookupPath` in `api/cue/inspect_test.go` describes the issue in more detail.
* Certain features (like dynamic labels) of CUE are not supported.