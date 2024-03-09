package controllers

import (
	commands "bitshare-chain/application/commands"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

type TestController struct{}

var (
	testCommand *commands.TestCommandHandler
)

type TestControllerer interface {
	TestRequest(context *gin.Context)
}

func NewTestController(command *commands.TestCommandHandler) TestControllerer {
	testCommand = command
	return &TestController{}
}

func (controller *TestController) TestRequest(context *gin.Context) {
	var testCmd commands.TestCommand
	if err := context.BindJSON(&testCmd); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	response, err := testCommand.Handle(context.Request.Context(), testCmd)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, response)
}
