package commands

import (
	"context"
	"encoding/json"
	"myapp/domain"
	"myapp/domain/viewmodels"
	"myapp/infrastructure/services"
)

type CreateWalletAccountCommand struct{}

type CreateWalletAccountCommandHandler struct {
	walletAccountRepository domain.WalletAccountRepository
	keyGenerator            services.KeyGenerator
}

func NewCreateWalletAccountCommandHandler(walletAccountRepository domain.WalletAccountRepository, keyGenerator services.KeyGenerator) *CreateWalletAccountCommandHandler {
	return &CreateWalletAccountCommandHandler{
		walletAccountRepository: walletAccountRepository,
		keyGenerator:            keyGenerator,
	}
}

// Handle executes the CreateWalletAccountCommand and returns a JSON response.
func (handler *CreateWalletAccountCommandHandler) Handle(context context.Context, cmd CreateWalletAccountCommand) (interface{}, error) {
	keys := handler.keyGenerator.GeneratePublicAndPrivateKey()

	newWalletAccount := domain.WalletAccount{
		Address: keys.PublicKey,
	}

	err := handler.walletAccountRepository.CreateWalletAccount(context, newWalletAccount)
	if err != nil {
		// Handle error accordingly
		return nil, err
	}

	walletKeysVM := viewmodels.NewWalletKeysVM(keys.PublicKey, keys.PrivateKey, nil)

	// Convert to JSON response
	response, err := json.Marshal(walletKeysVM)
	if err != nil {
		// Handle error accordingly
		return nil, err
	}

	return response, nil
}
