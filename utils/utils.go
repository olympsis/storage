package utils

import "os"

func GetServerConfig() *ServerConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}

	filePath := os.Getenv("CREDENTIALS_FILE_PATH")
	if filePath == "" {
		panic("Credentials file path required in environment variables")
	}

	return &ServerConfig{
		Port:                port,
		LogLevel:            logLevel,
		CredentialsFilePath: filePath,
	}
}
