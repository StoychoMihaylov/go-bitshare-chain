package commands

import (
	documents "bitshare-chain/application/data-access/documents"
	repositories "bitshare-chain/application/data-access/repositories"
	validation "bitshare-chain/application/validation"
	viewmodels "bitshare-chain/domain/view-models"
	services "bitshare-chain/infrastructure/utilities"
	context "context"
	json "encoding/json"
)

type CreateWalletAccountCommand struct{}

type CreateWalletAccountCommandHandler struct {
	walletAccountRepository repositories.WalletAccountRepository
	keyGenerator            services.KeyGenerator
	validator               *validation.Validator
}

func NewCreateWalletAccountCommandHandler(walletAccountRepository repositories.WalletAccountRepository, keyGenerator services.KeyGenerator, validator *validation.Validator) *CreateWalletAccountCommandHandler {
	return &CreateWalletAccountCommandHandler{
		walletAccountRepository: walletAccountRepository,
		keyGenerator:            keyGenerator,
		validator:               validator,
	}
}

func (handler *CreateWalletAccountCommandHandler) Handle(context context.Context, command CreateWalletAccountCommand) (interface{}, error) {
	publicKey, privateKey := handler.keyGenerator.GeneratePublicAndPrivateKey()

	newWalletAccount := documents.WalletAccountDocument{
		Address: publicKey,
	}

	err := handler.walletAccountRepository.CreateWalletAccount(&newWalletAccount)
	if err != nil {
		// Handle error accordingly
		return nil, err
	}

	walletKeysVM := viewmodels.NewWalletKeysVM(publicKey, privateKey, "")

	// Convert to JSON response
	response, err := json.Marshal(walletKeysVM)
	if err != nil {
		// Handle error accordingly
		return nil, err
	}

	return response, nil
}
