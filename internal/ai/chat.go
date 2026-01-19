package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ChatService struct {
	LLMURL string
	APIKey string
	Model  string
	Client *http.Client
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChatService(url, key, model string) *ChatService {
	return &ChatService{
		LLMURL: url,
		APIKey: key,
		Model:  model,
		Client: &http.Client{Timeout: 60 * time.Second}, // Keep timeout, but stream might exceed it if idle? Usually fine for chunks.
	}
}

// Chat (Non-streaming - kept for backward compatibility if needed, though we will likely switch Handler to Stream)
func (s *ChatService) Chat(messages []ChatMessage) (string, error) {
	reqBody := map[string]interface{}{
		"model":    s.Model,
		"messages": messages,
		"stream":   false,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", s.LLMURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.APIKey)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("LLM returned status: %s", resp.Status)
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// ChatStream streams the response chunks to the provided callback
func (s *ChatService) ChatStream(messages []ChatMessage, onChunk func(string) error) error {
	reqBody := map[string]interface{}{
		"model":    s.Model,
		"messages": messages,
		"stream":   true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", s.LLMURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.APIKey)
	}

	// Use a client without timeout for streaming or a longer one? 
	// Standard client has 60s. For streaming, we might hit it if a chunk takes long. 
	// But usually chunks are fast. 
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("LLM returned status: %s", resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		
		// Parse SSE format
		// Expected: "data: {JSON}"
		if !bytes.HasPrefix(line, []byte("data: ")) {
			continue
		}
		
		data := bytes.TrimPrefix(line, []byte("data: "))
		if string(data) == "[DONE]" {
			return nil
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(data, &chunk); err != nil {
			// Skip malformed chunks
			continue
		}

		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				if err := onChunk(content); err != nil {
					return err
				}
			}
		}
	}

	return scanner.Err()
}