package auth

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

// TokenClaims represents the claims in a PASETO token
type TokenClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"` // "access" or "refresh"
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	TokenID   string    `json:"jti"` // JWT ID for token tracking
}

// PasetoService handles PASETO token operations
type PasetoService struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// NewPasetoService creates a new PASETO service
func NewPasetoService(privateKeyHex, publicKeyHex string) (*PasetoService, error) {
	// Decode private key from hex
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	// Decode public key from hex
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	// Validate key lengths
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key length: expected %d, got %d", ed25519.PrivateKeySize, len(privateKeyBytes))
	}

	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key length: expected %d, got %d", ed25519.PublicKeySize, len(publicKeyBytes))
	}

	return &PasetoService{
		privateKey: ed25519.PrivateKey(privateKeyBytes),
		publicKey:  ed25519.PublicKey(publicKeyBytes),
	}, nil
}

// GenerateAccessToken generates a new access token
func (s *PasetoService) GenerateAccessToken(userID, email, role string, duration time.Duration) (string, error) {
	return s.generateToken(userID, email, role, "access", duration)
}

// GenerateRefreshToken generates a new refresh token
func (s *PasetoService) GenerateRefreshToken(userID, email, role string, duration time.Duration) (string, error) {
	return s.generateToken(userID, email, role, "refresh", duration)
}

// generateToken generates a PASETO v4.public token
func (s *PasetoService) generateToken(userID, email, role, tokenType string, duration time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	// Generate unique token ID (UUIDv7)
	tokenID, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate token ID: %w", err)
	}

	// Create token
	token := paseto.NewToken()

	// Set standard claims
	token.SetIssuedAt(now)
	token.SetNotBefore(now) // Add nbf claim for ValidAt() rule
	token.SetExpiration(expiresAt)
	token.SetJti(tokenID.String())
	token.SetSubject(userID)

	// Set custom claims
	token.SetString("user_id", userID)
	token.SetString("email", email)
	token.SetString("role", role)
	token.SetString("token_type", tokenType)

	// Sign token with private key
	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromEd25519(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create secret key: %w", err)
	}
	signedToken := token.V4Sign(secretKey, nil)

	return signedToken, nil
}

// ValidateToken validates a PASETO token and returns the claims
func (s *PasetoService) ValidateToken(tokenString string) (*TokenClaims, error) {
	// Create parser
	parser := paseto.NewParser()

	// Add validation rules
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	// Parse and verify token
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromEd25519(s.publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}
	token, err := parser.ParseV4Public(publicKey, tokenString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	userID, err := token.GetString("user_id")
	if err != nil {
		return nil, fmt.Errorf("failed to get user_id: %w", err)
	}

	email, err := token.GetString("email")
	if err != nil {
		return nil, fmt.Errorf("failed to get email: %w", err)
	}

	role, err := token.GetString("role")
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	tokenType, err := token.GetString("token_type")
	if err != nil {
		return nil, fmt.Errorf("failed to get token_type: %w", err)
	}

	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		return nil, fmt.Errorf("failed to get issued_at: %w", err)
	}

	expiresAt, err := token.GetExpiration()
	if err != nil {
		return nil, fmt.Errorf("failed to get expiration: %w", err)
	}

	tokenID, err := token.GetJti()
	if err != nil {
		return nil, fmt.Errorf("failed to get token_id: %w", err)
	}

	return &TokenClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		TokenID:   tokenID,
	}, nil
}

// GenerateKeyPair generates a new Ed25519 key pair for PASETO v4
func GenerateKeyPair() (privateKeyHex, publicKeyHex string, err error) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key pair: %w", err)
	}

	privateKeyHex = hex.EncodeToString(privateKey)
	publicKeyHex = hex.EncodeToString(publicKey)

	return privateKeyHex, publicKeyHex, nil
}
