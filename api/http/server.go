package http

import (
	"cueify/cue"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunServer(address string) {
	router := gin.Default()
	router.POST("/validate", validateValue)
	router.POST("/inspect", inspectValue)

	if err := router.Run(address); err != nil {
		panic(fmt.Sprintf("Could not start server on %s", address))
	}
}

type inspectBody struct {
	Path  []string    `json:"path"`
	Value interface{} `json:"value"`
}

func inspectValue(c *gin.Context) {
	var body inspectBody
	if err := c.BindJSON(&body); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid body"})
		return
	}

	jsonString, _ := json.Marshal(body.Value)
	properties := cue.Inspect(body.Path, string(jsonString))

	c.JSON(http.StatusOK, properties)
}

type validationResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors"`
}

type validationBody struct {
	Path  []string    `json:"path"`
	Value interface{} `json:"value"`
}

func validateValue(c *gin.Context) {
	var body validationBody
	if err := c.BindJSON(&body); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid body"})
		return
	}

	jsonString, _ := json.Marshal(body.Value)
	success, errors := cue.Validate(body.Path, string(jsonString))

	if errors != nil {
		c.JSON(http.StatusOK, validationResult{success, (*errors).Errors})
	} else {
		c.JSON(http.StatusOK, validationResult{success, nil})
	}
}
