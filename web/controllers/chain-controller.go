package controllers

import (
	commands "bitshare-chain/application/commands"
	services "bitshare-chain/application/services"
	bindingmodels "bitshare-chain/domain/binding-models"
	enums "bitshare-chain/domain/enums"
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	hex "encoding/hex"
	errors "errors"
	big "math/big"
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
	controller.ginRouter.POST("/api/get-pending-transaction", controller.GetPendingTransaction)
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
func (controller *ChainController) RequestTransaction(context *gin.Context) {

	// !!! Temporary !!!
	// Could also be signed in the front-end end sent to the back-end (probably more secured way since the private key is not traveling in the network)
	// !!! Temporary !!!

	var transactionBM bindingmodels.TransactionBindingModel
	if err := context.BindJSON(&transactionBM); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	privateKey := context.Query("privateKey")
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid private key"})
		return
	}

	privateKeyECDSA, err := parsePrivateKey(privateKeyBytes)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate keys"})
		return
	}

	if err := transactionBM.SignTransaction(privateKeyECDSA); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controller.cache.Store(enums.PendingTransactions, transactionBM)

	context.JSON(http.StatusOK, gin.H{"message": "Transaction signed successfully"})
}

func parsePrivateKey(privateKeyBytes []byte) (*ecdsa.PrivateKey, error) {
	curve := elliptic.P256()
	privateKey := new(ecdsa.PrivateKey)
	privateKey.Curve = curve
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.Curve = curve
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(privateKeyBytes)
	if privateKey.PublicKey.X == nil || privateKey.PublicKey.Y == nil {
		return nil, errors.New("invalid private key")
	}
	return privateKey, nil
}

// "POST" "/api/get-pending-transaction"
func (controller *ChainController) GetPendingTransaction(context *gin.Context) {
	var pendingTransactions []bindingmodels.TransactionBindingModel

	values, ok := controller.cache.Load(enums.PendingTransactions)
	if !ok {
		context.JSON(http.StatusOK, values)
	} else {
		context.JSON(http.StatusNotFound, "")
	}

	// Iterate over all entries in the cache
	controller.cache.Range(func(key, value interface{}) bool {
		// Check if the value is of type TransactionBindingModel
		if transaction, ok := value.(bindingmodels.TransactionBindingModel); ok {
			// Append the transaction to the pendingTransactions slice
			pendingTransactions = append(pendingTransactions, transaction)
		}
		return true // Continue iterating
	})

	context.JSON(http.StatusOK, pendingTransactions)
}
