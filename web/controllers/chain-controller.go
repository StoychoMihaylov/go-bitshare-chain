package controllers

import (
	commands "bitshare-chain/application/commands"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

type ChainController struct{}

var (
	createWalletAccountCommandHandler *commands.CreateWalletAccountCommandHandler
)

type ChainControllerer interface {
	// SetRewardAndBlockSigningKeys(context *gin.Context)
	// RequestTransaction(context *gin.Context)
	// GetPendingTransaction(context *gin.Context)
	// MineTransactions(context *gin.Context)
	// GetBalanceOfAddress(context *gin.Context)
	CreateNewWalletAccount(context *gin.Context)
	// IsTheChainValid()
}

func NewChainController(commnad *commands.CreateWalletAccountCommandHandler) ChainControllerer {
	createWalletAccountCommandHandler = commnad
	return &ChainController{}
}

func (controller *ChainController) CreateNewWalletAccount(context *gin.Context) {
	var createWalletAccountCommand commands.CreateWalletAccountCommand

	response, err := createWalletAccountCommandHandler.Handle(context.Request.Context(), createWalletAccountCommand)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseBytes, ok := response.([]byte)
	if !ok {
		context.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	context.Data(200, "application/json; charset=utf-8", responseBytes)
}
