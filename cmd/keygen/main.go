package main

import (
	"fmt"
	"log"

	"github.com/pentsecops/backend/pkg/auth"
)

func main() {
	fmt.Println("Generating PASETO v4 Ed25519 Key Pair...")
	fmt.Println("========================================")

	privateKey, publicKey, err := auth.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	fmt.Println("\n✅ Key pair generated successfully!")
	fmt.Println("\nAdd these to your .env file:")
	fmt.Println("========================================")
	fmt.Printf("PASETO_PRIVATE_KEY=%s\n", privateKey)
	fmt.Printf("PASETO_PUBLIC_KEY=%s\n", publicKey)
	fmt.Println("========================================")
	fmt.Println("\n⚠️  IMPORTANT: Keep the private key secret!")
	fmt.Println("⚠️  Never commit the private key to version control!")
}

