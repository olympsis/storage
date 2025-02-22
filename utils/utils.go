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

	firebaseFilePath := os.Getenv("FIREBASE_FILE_PATH")
	if firebaseFilePath == "" {
		panic("firebase file path required in environment variables")
	}

	return &ServerConfig{
		Port:             port,
		LogLevel:         logLevel,
		FirebaseFilePath: firebaseFilePath,
	}
}
