package config

import "os" // Import the os package to access environment variables

// GetInfuraURL retrieves the Infura URL from the environment variables.
// It expects that the INFURA_URL environment variable is set.
func GetInfuraURL() string {
	return os.Getenv("INFURA_URL") // Returns the value of the INFURA_URL environment variable
}

// GetAlchemyURL retrieves the Alchemy URL from the environment variables.
// It expects that the ALCHEMY_URL environment variable is set.
func GetAlchemyURL() string {
	return os.Getenv("ALCHEMY_URL") // Returns the value of the ALCHEMY_URL environment variable
}
