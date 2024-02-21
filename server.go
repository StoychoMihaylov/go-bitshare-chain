package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
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

	http.ListenAndServe(":8000", nil)
}
