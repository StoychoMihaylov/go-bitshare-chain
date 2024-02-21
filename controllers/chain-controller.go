package controllers

import (
	"github.com/gin-gonic/gin"
)

type controller struct{}

type ChainController interface {
	SetRewardAndBlockSigningKeys(context *gin.Context)
	RequestTransaction(context *gin.Context)
	GetPendingTransaction(context *gin.Context)
	MineTransactions(context *gin.Context)
	GetBalanceOfAddress()
	CreateNewWalletAccount()
	IsTheChainValid()
}
