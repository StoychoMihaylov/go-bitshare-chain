package services

type KeyGeneratorService interface {
	GeneratePublicAndPrivateKey() (publicKey string, privateKey string)
}

type service struct{}
