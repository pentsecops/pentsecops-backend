package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"io"

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
	       userRole := c.Locals("user_role")
	       if userID == nil || userRole == nil {
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
	       groqUrl := "https://api.groq.com/openai/v1/responses"
	       groqReq := map[string]interface{}{
		       "model": "openai/gpt-oss-20b",
		       "input": req.Message,
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
		       // Debug: log the full Groq API response body
		       var rawBody bytes.Buffer
		       tee := io.TeeReader(groqResp.Body, &rawBody)
			       var groqResult struct {
				       Output []struct {
					       Type    string `json:"type"`
					       Role    string `json:"role"`
					       Content []struct {
						       Type string `json:"type"`
						       Text string `json:"text"`
					       } `json:"content"`
				       } `json:"output"`
				       Error struct {
					       Message string `json:"message"`
				       } `json:"error"`
			       }
			       if err := json.NewDecoder(tee).Decode(&groqResult); err != nil {
				       println("[Groq API Debug] Raw response:", rawBody.String())
				       return utils.InternalServerError(c, "Failed to parse LLM API response")
			       }
			       if groqResult.Error.Message != "" {
				       println("[Groq API Debug] Error message:", groqResult.Error.Message)
				       println("[Groq API Debug] Raw response:", rawBody.String())
				       return utils.InternalServerError(c, "LLM API error: "+groqResult.Error.Message)
			       }
			       llmResponse := ""
			       // Extract assistant message text from output
			       for _, out := range groqResult.Output {
				       if out.Type == "message" && out.Role == "assistant" {
					       for _, c := range out.Content {
						       if c.Type == "output_text" && c.Text != "" {
							       llmResponse = c.Text
							       break
						       }
					       }
				       }
				       if llmResponse != "" {
					       break
				       }
			       }
			       if llmResponse == "" {
				       println("[Groq API Debug] Could not extract assistant message. Raw body:", rawBody.String())
				       return utils.InternalServerError(c, "LLM API returned no usable response. Please check your API key, quota, or request format.")
			       }

	// Save to llm_messages (pseudo, implement actual DB logic as needed)
	       msgID := uuid.New()
	       createdAt := time.Now().Format(time.RFC3339)
	       // TODO: Insert (msgID, userID, userRole, req.Message, llmResponse, createdAt) into llm_messages table

	return utils.Success(c, dto.LLMQueryResponse{
		       ID:        msgID,
		       Message:   req.Message,
		       Response:  llmResponse,
		       CreatedAt: createdAt,
	       }, "LLM response returned successfully")
}
