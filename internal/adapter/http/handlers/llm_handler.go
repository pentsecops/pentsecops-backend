package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/utils"
)

type LLMHandler struct{}

func NewLLMHandler() *LLMHandler {
	return &LLMHandler{}
}

func (h *LLMHandler) Query(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	role := c.Locals("role")
	if userID == nil || role == nil {
		return utils.Unauthorized(c, "Unauthorized")
	}

	var req dto.LLMQueryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	if req.Message == "" {
		return utils.BadRequest(c, "Message is required", nil)
	}

	// Call Groq LLM API
	groqApiKey := os.Getenv("GROQ_API")
	groqUrl := "https://api.groq.com/v1/chat/completions"
	groqReq := map[string]interface{}{
		"model":    "llama-3.3-70b-versatile",
		"messages": []map[string]string{{"role": "user", "content": req.Message}},
	}
	groqBody, _ := json.Marshal(groqReq)
	groqHttpReq, _ := http.NewRequest("POST", groqUrl, bytes.NewBuffer(groqBody))
	groqHttpReq.Header.Set("Authorization", "Bearer "+groqApiKey)
	groqHttpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	groqResp, err := client.Do(groqHttpReq)
	if err != nil {
		return utils.InternalServerError(c, "Failed to contact LLM API")
	}
	defer groqResp.Body.Close()
	var groqResult struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.NewDecoder(groqResp.Body).Decode(&groqResult)
	llmResponse := ""
	if len(groqResult.Choices) > 0 {
		llmResponse = groqResult.Choices[0].Message.Content
	}

	// Save to llm_messages (pseudo, implement actual DB logic as needed)
	msgID := uuid.New()
	createdAt := time.Now().Format(time.RFC3339)
	// TODO: Insert (msgID, userID, role, req.Message, llmResponse, createdAt) into llm_messages table

	return utils.Success(c, dto.LLMQueryResponse{
		ID:        msgID,
		Message:   req.Message,
		Response:  llmResponse,
		CreatedAt: createdAt,
	}, "LLM response returned successfully")
}
