package main

import (
	"fmt"
	"os"
	
	"awsomeshop/backend/pkg/config"
)

func main() {
	// Create a test .env file
	envContent := `DB_HOST=testhost
DB_PORT=3307
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
JWT_SECRET=test-secret-key
SERVER_PORT=8081
UPLOAD_DIR=/test/uploads/
MAX_UPLOAD_SIZE=2097152
LOG_LEVEL=DEBUG
`
	
	// Write test .env file
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		fmt.Printf("Failed to create .env file: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(".env")
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}
	
	// Verify values from .env file
	fmt.Println("Configuration loaded successfully from .env file:")
	fmt.Printf("DB_HOST: %s\n", cfg.DBHost)
	fmt.Printf("DB_PORT: %s\n", cfg.DBPort)
	fmt.Printf("DB_USER: %s\n", cfg.DBUser)
	fmt.Printf("DB_PASSWORD: %s\n", cfg.DBPassword)
	fmt.Printf("DB_NAME: %s\n", cfg.DBName)
	fmt.Printf("JWT_SECRET: %s\n", cfg.JWTSecret)
	fmt.Printf("SERVER_PORT: %s\n", cfg.ServerPort)
	fmt.Printf("UPLOAD_DIR: %s\n", cfg.UploadDir)
	fmt.Printf("MAX_UPLOAD_SIZE: %d\n", cfg.MaxUploadSize)
	fmt.Printf("LOG_LEVEL: %s\n", cfg.LogLevel)
	
	// Verify the values match what we set in .env
	if cfg.DBHost != "testhost" {
		fmt.Printf("ERROR: Expected DBHost='testhost', got '%s'\n", cfg.DBHost)
		os.Exit(1)
	}
	if cfg.DBPort != "3307" {
		fmt.Printf("ERROR: Expected DBPort='3307', got '%s'\n", cfg.DBPort)
		os.Exit(1)
	}
	if cfg.JWTSecret != "test-secret-key" {
		fmt.Printf("ERROR: Expected JWTSecret='test-secret-key', got '%s'\n", cfg.JWTSecret)
		os.Exit(1)
	}
	
	fmt.Println("\n✓ All values from .env file loaded correctly!")
}
