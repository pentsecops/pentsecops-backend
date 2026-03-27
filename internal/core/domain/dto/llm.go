package dto

import "github.com/google/uuid"

// LLMQueryRequest represents a request to the LLM chat endpoint
type LLMQueryRequest struct {
	Message string `json:"message" validate:"required,min=1,max=2000"`
}

// LLMQueryResponse represents a response from the LLM chat endpoint
type LLMQueryResponse struct {
	ID        uuid.UUID `json:"id"`
	Message   string    `json:"message"`
	Response  string    `json:"response"`
	CreatedAt string    `json:"created_at"`
}
