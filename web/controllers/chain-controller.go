package controllers

import (
	commands "bitshare-chain/application/commands"
	services "bitshare-chain/application/services"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

type ChainController struct {
	ginRouter                         *gin.Engine
	createWalletAccountCommandHandler *commands.CreateWalletAccountCommandHandler
	metadataService                   *services.MetadataService
}

type ChainControllerer interface {
	SetupChainController()
	// SetRewardAndBlockSigningKeys(context *gin.Context)
	// RequestTransaction(context *gin.Context)
	// GetPendingTransaction(context *gin.Context)
	// MineTransactions(context *gin.Context)
	// GetBalanceOfAddress(context *gin.Context)
	CreateNewWalletAccount(context *gin.Context)
	// IsTheChainValid()
}

func NewChainController(ginRouter *gin.Engine, createWalletAccountCommandHandler *commands.CreateWalletAccountCommandHandler, metadataService *services.MetadataService) ChainControllerer {
	return &ChainController{
		ginRouter:                         ginRouter,
		createWalletAccountCommandHandler: createWalletAccountCommandHandler,
		metadataService:                   metadataService,
	}
}

func (controller *ChainController) SetupChainController() {
	controller.ginRouter.POST("/api/create-wallet", controller.CreateNewWalletAccount)
	controller.ginRouter.PUT("/api/set-block-signing-keys", controller.SetBlockSigningKeys)
}

// "POST" "api/create-wallet"
func (controller *ChainController) CreateNewWalletAccount(context *gin.Context) {
	var createWalletAccountCommand commands.CreateWalletAccountCommand

	response, err := controller.createWalletAccountCommandHandler.Handle(context.Request.Context(), createWalletAccountCommand)
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

// "POST" "api/set-block-signing-keys"
func (controller *ChainController) SetBlockSigningKeys(context *gin.Context) {
	privateKey := context.Query("privateKey")

	keys := controller.metadataService.CreateOrUpdateKeys(privateKey)

	context.JSON(http.StatusOK, keys)
}
