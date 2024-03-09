package main

import (
	commands "bitshare-chain/application/commands"
	mongo_context "bitshare-chain/application/data-access/context"
	repositories "bitshare-chain/application/data-access/repositories"
	validation "bitshare-chain/application/validation"
	settings "bitshare-chain/infrastructure/settings"
	services "bitshare-chain/infrastructure/utilities"
	controllers "bitshare-chain/web/controllers"

	gin "github.com/gin-gonic/gin"
)

func main() {
	ginRouter := gin.Default()

	dbOptions := settings.MongoDbOptions{
		DatabaseName:     "GoBitshareChain",
		ConnectionString: "mongodb://root:rootpassword@go-bitshare-mongodb:27017/?authMechanism=SCRAM-SHA-256",
	}

	mongoContext, err := mongo_context.NewMongoContext(&dbOptions)
	if err != nil {
		panic(err)
	}
	defer mongoContext.Close()

	// Initialize dependencies
	walletAccountRepository := repositories.NewWalletAccountRepository(mongoContext)
	keyGenerator := &services.KeyGenerator{}

	validator := validation.NewValidator()
	createWalletAccountCommandHandler := commands.NewCreateWalletAccountCommandHandler(walletAccountRepository, *keyGenerator, validator)

	// Correct the variable name here to avoid conflicts
	chainController := controllers.NewChainController(createWalletAccountCommandHandler)

	testHandler := commands.NewTestCommandHandler(validator)
	testController := controllers.NewTestController(testHandler)

	//Routs
	ginRouter.POST("/test", testController.TestRequest)

	ginRouter.POST("/create-wallet", chainController.CreateNewWalletAccount)

	ginRouter.Run(":8000")
}
