package config

import "os"

func GetInfuraURL() string {
	return os.Getenv("INFURA_URL")
}

func GetAlchemyURL() string {
	return os.Getenv("ALCHEMY_URL")
}
