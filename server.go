package main

import (
	commands "bitshare-chain/application/commands"
	repositories "bitshare-chain/application/data-access/repositories"
	settings_options "bitshare-chain/application/options"
	services "bitshare-chain/infrastructure/services"
	http "net/http"

	gin "github.com/gin-gonic/gin"
	mux "github.com/gorilla/mux"
)

func main() {
	ginRouter := gin.Default()
	ginRouter.GET("/gin", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello from GIN!")
	})

	gorillaRouter := mux.NewRouter()
	gorillaRouter.HandleFunc("/gorilla", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("Hello from GORILLA Mux!"))
	}).Methods("GET")

	// Create a single server for GIN and GORILLA
	http.Handle("/gin", ginRouter)
	http.Handle("/gorilla", gorillaRouter)

	//TO DO: FIND A WAY TO MAKE THAT SIMPLE!!!!
	//---------------------------------------------------------------------------------------------------------------------------------

	dbOptions := settings_options.MongoDbOptions{
		DatabaseName:     "GoBitshareChain",
		ConnectionString: "mongodb://root:rootpassword@go-bitshare-mongodb:27017/?authMechanism=SCRAM-SHA-256",
	}

	// Initialize dependencies
	walletAccountRepository := repositories.NewWalletAccountRepository(dbOptions)
	keyGenerator := &services.KeyGenerator{} // Assuming you have a KeyGenerator implementation

	// Inject dependencies into the handler
	createWalletAccountHandler := commands.NewCreateWalletAccountCommandHandler(walletAccountRepository, keyGenerator) // walletAccountRepository still not implemented

	ginRouter.GET("/api/test", func(context *gin.Context) {
		context.String(http.StatusOK, "TEST!")
	})

	ginRouter.POST("/api/create-wallet", func(context *gin.Context) {
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

	http.ListenAndServe(":8000", nil)
}
