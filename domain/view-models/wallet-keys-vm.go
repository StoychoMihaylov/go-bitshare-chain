package viewmodels

// WalletKeysVM represents the view model for wallet keys.
type WalletKeysVM struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
	NodeId     string `json:"nodeId"`
}

// NewWalletKeysVM creates a new instance of WalletKeysVM.
func NewWalletKeysVM(publicKey, privateKey, nodeId string) WalletKeysVM {
	return WalletKeysVM{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		NodeId:     nodeId,
	}
}
