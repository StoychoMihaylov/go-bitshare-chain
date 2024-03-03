package commands

import (
	repositories "bitshare-chain/application/data-access/repositories"
	viewmodels "bitshare-chain/domain/viewmodels"
	services "bitshare-chain/infrastructure/services"
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

// Handle executes the CreateWalletAccountCommand and returns a JSON response.
func (handler *CreateWalletAccountCommandHandler) Handle(context context.Context, cmd CreateWalletAccountCommand) (interface{}, error) {
	keys := handler.keyGenerator.GeneratePublicAndPrivateKey()

	newWalletAccount := repositories.WalletAccount{
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
