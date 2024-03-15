package controllers

import (
	commands "bitshare-chain/application/commands"
	services "bitshare-chain/application/services"
	bindingmodels "bitshare-chain/domain/binding-models"
	enums "bitshare-chain/domain/enums"
	bytes "bytes"
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	hex "encoding/hex"
	http "net/http"
	sync "sync"

	gin "github.com/gin-gonic/gin"
)

type ChainController struct {
	cache                             sync.Map
	ginRouter                         *gin.Engine
	createWalletAccountCommandHandler *commands.CreateWalletAccountCommandHandler
	metadataService                   *services.MetadataService
}

type ChainControllerer interface {
	SetupChainController()
	SetBlockSigningKeys(context *gin.Context)
	RequestTransaction(context *gin.Context)
	// GetPendingTransaction(context *gin.Context)
	// MineTransactions(context *gin.Context)
	// GetBalanceOfAddress(context *gin.Context)
	CreateNewWalletAccount(context *gin.Context)
	// IsTheChainValid()
}

func NewChainController(
	ginRouter *gin.Engine,
	createWalletAccountCommandHandler *commands.CreateWalletAccountCommandHandler,
	metadataService *services.MetadataService) ChainControllerer {
	return &ChainController{
		ginRouter:                         ginRouter,
		createWalletAccountCommandHandler: createWalletAccountCommandHandler,
		metadataService:                   metadataService,
	}
}

func (controller *ChainController) SetupChainController() {
	controller.ginRouter.POST("/api/create-wallet", controller.CreateNewWalletAccount)
	controller.ginRouter.POST("/api/set-block-signing-keys", controller.SetBlockSigningKeys)
	controller.ginRouter.POST("/api/request-transaction", controller.RequestTransaction)
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
	if privateKey == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	keys, err := controller.metadataService.CreateOrUpdateKeys(privateKey)
	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
	}

	context.JSON(http.StatusOK, keys)
}

// "POST" "/api/request-transaction"
func (controller *ChainController) RequestTransaction(c *gin.Context) {

	// !!! Temporary !!!
	// Could also be signed in the front-end end sent to the back-end (probably more secured way since the private key is not traveling in the network)
	// !!! Temporary !!!

	var transactionBM bindingmodels.TransactionBindingModel
	if err := c.BindJSON(&transactionBM); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	privateKey := c.Query("privateKey")
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key"})
		return
	}

	privateKeyECDSA, err := ecdsa.GenerateKey(elliptic.P256(), bytes.NewReader(privateKeyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate keys"})
		return
	}

	if err := transactionBM.SignTransaction(privateKeyECDSA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controller.cache.Store(enums.PendingTransactions, transactionBM)

	c.JSON(http.StatusOK, gin.H{"message": "Transaction signed successfully"})
}
