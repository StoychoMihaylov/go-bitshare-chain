package services

import (
	repositories "bitshare-chain/application/data-access/repositories"
	enums "bitshare-chain/domain/enums"
	viewmodels "bitshare-chain/domain/view-models"
	services "bitshare-chain/infrastructure/utilities"
	ecdsa "crypto/ecdsa"
	elliptic "crypto/elliptic"
	"crypto/rand"
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
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return viewmodels.WalletKeysVM{}, err
	}

	keys := &ecdsa.PrivateKey{}
	keys.Curve = elliptic.P256()
	keys.D = new(big.Int).SetBytes(privateKeyBytes)

	random := rand.Reader
	_, x, y, err := elliptic.GenerateKey(keys.Curve, random)
	if err != nil {
		return viewmodels.WalletKeysVM{}, err
	}

	keys.PublicKey = ecdsa.PublicKey{Curve: keys.Curve, X: x, Y: y}

	publicKeyExport := elliptic.Marshal(keys.Curve, keys.PublicKey.X, keys.PublicKey.Y)

	publicKey := hex.EncodeToString(publicKeyExport)
	privateKeysExport := hex.EncodeToString(privateKeyBytes)

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
