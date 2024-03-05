package commands

import (
	documents "bitshare-chain/application/data-access/documents"
	repositories "bitshare-chain/application/data-access/repositories"
	viewmodels "bitshare-chain/domain/view-models"
	services "bitshare-chain/infrastructure/utilities"
	context "context"
	json "encoding/json"
)

type CreateWalletAccountCommand struct{}

type CreateWalletAccountCommandHandler struct {
	walletAccountRepository repositories.WalletAccountRepository
	keyGenerator            services.KeyGenerator
}

func NewCreateWalletAccountCommandHandler(walletAccountRepository repositories.WalletAccountRepository, keyGenerator services.KeyGenerator) *CreateWalletAccountCommandHandler {
	return &CreateWalletAccountCommandHandler{
		walletAccountRepository: walletAccountRepository,
		keyGenerator:            keyGenerator,
	}
}

func (handler *CreateWalletAccountCommandHandler) Handle(context context.Context, cmd CreateWalletAccountCommand) (interface{}, error) {
	// Generate public and private keys
	publicKey, privateKey := handler.keyGenerator.GeneratePublicAndPrivateKey()

	newWalletAccount := documents.WalletAccountDocument{
		Address: publicKey,
	}

	// Create a new wallet account
	err := handler.walletAccountRepository.CreateWalletAccount(&newWalletAccount)
	if err != nil {
		// Handle error accordingly
		return nil, err
	}

	// Create a view model for wallet keys
	walletKeysVM := viewmodels.NewWalletKeysVM(publicKey, privateKey, "")

	// Convert to JSON response
	response, err := json.Marshal(walletKeysVM)
	if err != nil {
		// Handle error accordingly
		return nil, err
	}

	return response, nil
}
