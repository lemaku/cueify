package main

import "cueify/http"

func main() {
	//success, errors := cue.Validate([]string{"universities", "tuwien", "students", "0", "matNr"}, cue.MissingPropVal)
	//fmt.Println(success)
	//fmt.Println(errors)

	http.RunServer("localhost:8080")
}
