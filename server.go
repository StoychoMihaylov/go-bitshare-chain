package main

import (
	commands "bitshare-chain/application/commands"
	mongo_context "bitshare-chain/application/data-access/context"
	repositories "bitshare-chain/application/data-access/repositories"
	settings "bitshare-chain/infrastructure/settings"
	services "bitshare-chain/infrastructure/utilities"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

func main() {
	ginRouter := gin.Default()

	ginRouter.GET("/gin", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello from GIN!")
	})

	ginRouter.GET("/test", func(context *gin.Context) {
		context.String(http.StatusOK, "TEST!")
	})

	//TO DO: FIND A WAY TO MAKE THAT SIMPLE!!!!
	//---------------------------------------------------------------------------------------------------------------------------------

	dbOptions := settings.MongoDbOptions{
		DatabaseName:     "GoBitshareChain",
		ConnectionString: "mongodb://root:rootpassword@go-bitshare-mongodb:27017/?authMechanism=SCRAM-SHA-256",
	}

	mongoContext, err := mongo_context.NewMongoContext(&dbOptions)
	if err != nil {
		// Handle error
		panic(err)
	}
	defer mongoContext.Close()

	// Initialize dependencies
	walletAccountRepository := repositories.NewWalletAccountRepository(mongoContext)
	keyGenerator := &services.KeyGenerator{} // Assuming you have a KeyGenerator implementation

	// TO DO: Initialize validation middleware
	// TO DO: Add validation middlewere package
	//validator := validation.NewValidator()

	// Inject dependencies into the handler
	createWalletAccountHandler := commands.NewCreateWalletAccountCommandHandler(walletAccountRepository, *keyGenerator) // walletAccountRepository still not implemented

	ginRouter.POST("/create-wallet", func(context *gin.Context) {
		response, err := createWalletAccountHandler.Handle(context.Request.Context(), commands.CreateWalletAccountCommand{})
		if err != nil {
			context.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		// Type assert the response to []byte
		responseBytes, ok := response.([]byte)
		if !ok {
			context.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		// Assuming response is a JSON-encoded byte slice
		context.Data(200, "application/json; charset=utf-8", responseBytes)
	})

	//---------------------------------------------------------------------------------------------------------------------------------
	ginRouter.Run(":8000")
}
