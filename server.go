package main

import (
	commands "bitshare-chain/application/commands"
	mongo_context "bitshare-chain/application/data-access/context"
	repositories "bitshare-chain/application/data-access/repositories"
	services "bitshare-chain/application/services"
	validation "bitshare-chain/application/validation"
	settings "bitshare-chain/infrastructure/settings"
	utilities "bitshare-chain/infrastructure/utilities"
	controllers "bitshare-chain/web/controllers"

	gin "github.com/gin-gonic/gin"
)

func main() {
	ginRouter := gin.Default()

	// Initialize dependencies
	//DB
	dbOptions := settings.MongoDbOptions{
		DatabaseName:     "GoBitshareChain",
		ConnectionString: "mongodb://root:rootpassword@go-bitshare-mongodb:27000/?authMechanism=SCRAM-SHA-256",
	}

	mongoContext, err := mongo_context.NewMongoContext(&dbOptions)
	if err != nil {
		panic(err)
	}
	defer mongoContext.Close()

	//REPOS
	walletAccountRepository := repositories.NewWalletAccountRepository(mongoContext)
	nodeMetadataRepository := repositories.NewNodeMetadataRepository(mongoContext)

	//SERVICES
	keyGenerator := &utilities.KeyGenerator{}
	metadataService := services.NewMetadataService(*keyGenerator, *nodeMetadataRepository)

	//VALIDATOR
	validator := validation.NewValidator()

	//COMMANDS
	createWalletAccountCommandHandler := commands.NewCreateWalletAccountCommandHandler(walletAccountRepository, *keyGenerator, validator)
	testHandler := commands.NewTestCommandHandler(validator)

	//CONTROLLERS
	chainController := controllers.NewChainController(ginRouter, createWalletAccountCommandHandler, metadataService)
	chainController.SetupChainController()

	testController := controllers.NewTestController(ginRouter, testHandler)
	testController.SetupTestController()

	ginRouter.Run(":8000")
}
