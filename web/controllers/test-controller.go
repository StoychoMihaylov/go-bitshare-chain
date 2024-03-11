package controllers

import (
	commands "bitshare-chain/application/commands"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

type TestController struct {
	ginRouter   *gin.Engine
	testCommand *commands.TestCommandHandler
}

type TestControllerer interface {
	SetupTestController()
	TestRequest(context *gin.Context)
}

func NewTestController(router *gin.Engine, command *commands.TestCommandHandler) TestControllerer {
	return &TestController{
		ginRouter:   router,
		testCommand: command,
	}
}

func (controller *TestController) SetupTestController() {
	controller.ginRouter.POST("/api/test", controller.TestRequest)
}

// "POST" "/test"
func (controller *TestController) TestRequest(context *gin.Context) {
	var testCmd commands.TestCommand
	if err := context.BindJSON(&testCmd); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	response, err := controller.testCommand.Handle(context.Request.Context(), testCmd)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, response)
}
