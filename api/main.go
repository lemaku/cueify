package main

import (
	"cueify/http"
)

func main() {
	//fmt.Println(cue.Inspect([]string{"universities", "tuwien", "students", "0"}, cue.MissingPropVal))
	//fmt.Println(cue.Inspect([]string{"universities", "tuwien", "students"}, cue.MissingPropVal))
	//fmt.Println(cue.Inspect([]string{"universities"}, cue.MissingPropVal))
	//fmt.Println(success)
	//fmt.Println(errors)
	//fmt.Println(cue.Summarize(cue.CorrectVal).Value)

	http.RunServer("localhost:8080")
}
