package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"golang.org/x/crypto/argon2"
)

const (
	Argon2Time    = 1
	Argon2Memory  = 64 * 1024
	Argon2Threads = 4
	Argon2KeyLen  = 32
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hashpass <password>")
		os.Exit(1)
	}

	password := os.Args[1]
	salt := []byte("pentsecops-salt-2025")
	hash := argon2.IDKey([]byte(password), salt, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLen)
	hashStr := hex.EncodeToString(hash)

	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", hashStr)
}

