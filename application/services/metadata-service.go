package services

import (
	repositories "bitshare-chain/application/data-access/repositories"
	enums "bitshare-chain/domain/enums"
	viewmodels "bitshare-chain/domain/view-models"
	services "bitshare-chain/infrastructure/utilities"
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	hex "encoding/hex"
	big "math/big"
	sync "sync"
)

type MetadataService struct {
	cache                  sync.Map
	keyGenerator           services.KeyGenerator
	nodeMetadataRepository repositories.NodeMetadataRepository
}

func NewMetadataService(keyGenerator services.KeyGenerator, nodeMetadataRepository repositories.NodeMetadataRepository) *MetadataService {
	return &MetadataService{
		keyGenerator:           keyGenerator,
		nodeMetadataRepository: nodeMetadataRepository,
	}
}

func (service *MetadataService) CreateOrUpdateKeys(privateKeyHex string) (viewmodels.WalletKeysVM, error) {
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)

	keys, err := ecdsa.GenerateKey(elliptic.P256(), nil)
	if err != nil {
		return viewmodels.WalletKeysVM{}, err
	}

	keys.D = new(big.Int).SetBytes(privateKeyBytes)
	publicKeyExport := elliptic.Marshal(keys.Curve, keys.PublicKey.X, keys.PublicKey.Y)

	publicKey := hex.EncodeToString(publicKeyExport)
	privateKeysExport := hex.EncodeToString(privateKeyBytes)

	// Store in-memory cache
	service.cache.Store(enums.MinerPrivateKey, privateKeysExport)

	nodeId, err := service.nodeMetadataRepository.InsertMetadataPublicKey(publicKey)
	if err != nil {
		return viewmodels.WalletKeysVM{}, err
	}

	return viewmodels.WalletKeysVM{
		PublicKey:  publicKey,
		PrivateKey: privateKeysExport,
		NodeId:     nodeId,
	}, nil
}
