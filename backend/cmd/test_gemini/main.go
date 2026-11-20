package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/deadlium/eventbeep/backend/internal/config"
	"github.com/deadlium/eventbeep/backend/internal/llm"
)

func main() {
	fmt.Println("======================================================")
	fmt.Println(" EventBeep - Gemini API Key Verification Utility")
	fmt.Println("======================================================")

	// Load .env if present
	config.LoadEnvFile()

	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	model := strings.TrimSpace(os.Getenv("GEMINI_MODEL"))
	if model == "" {
		if strings.HasPrefix(apiKey, "AQ.") {
			model = "gemini-3.8-flash"
		} else {
			model = "gemini-1.5-flash"
		}
	} else if strings.HasPrefix(apiKey, "AQ.") && model == "gemini-1.5-flash" {
		model = "gemini-3.8-flash"
	}

	if apiKey == "" {
		fmt.Println("❌ FAILED: GEMINI_API_KEY is empty or not set.")
		fmt.Println("")
		fmt.Println("👉 How to fix:")
		fmt.Println("1. Open backend/.env in your editor.")
		fmt.Println("2. Find line 37: GEMINI_API_KEY=")
		fmt.Println("3. Paste your key: GEMINI_API_KEY=AIzaSy...")
		fmt.Println("4. Save the file (Cmd+S or Ctrl+S).")
		fmt.Println("5. Run this test again:")
		fmt.Println("   go run ./cmd/test_gemini")
		os.Exit(1)
	}

	// Mask the key for display
	maskedKey := apiKey
	if len(apiKey) > 10 {
		maskedKey = apiKey[:6] + "..." + apiKey[len(apiKey)-4:]
	}

	endpointName := "v1beta/models:generateContent"
	if strings.HasPrefix(apiKey, "AQ.") || strings.Contains(model, "3.8") {
		endpointName = "v1beta/interactions"
	}
	fmt.Printf("🔑 Found API Key : %s (length: %d chars)\n", maskedKey, len(apiKey))
	fmt.Printf("🤖 Target Model  : %s\n", model)
	fmt.Printf("📡 Endpoint      : %s\n", endpointName)
	fmt.Println("📡 Connecting to Google Gemini API (timeout: 2 minutes)...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	extractor := llm.NewGeminiExtractor(apiKey, llm.WithModel(model), llm.WithTimeout(2*time.Minute))
	entity, err := extractor.ValidateKey(ctx)
	if err != nil {
		fmt.Println("")
		fmt.Println("❌ FAILED: Google Gemini API returned an error:")
		fmt.Printf("   %v\n\n", err)

		errStr := err.Error()
		if strings.Contains(errStr, "400") || strings.Contains(errStr, "API_KEY_INVALID") {
			fmt.Println("💡 Hint: The API key appears to be invalid or mistyped.")
			fmt.Println("   Verify your key in Google AI Studio: https://aistudio.google.com/app/apikey")
		} else if strings.Contains(errStr, "403") {
			fmt.Println("💡 Hint: Permission denied. Your API key might not have access to the model, or billing/quota restrictions apply.")
		} else if strings.Contains(errStr, "429") {
			fmt.Println("💡 Hint: Rate limit or quota exceeded for this Gemini API key.")
		} else if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "timeout") {
			fmt.Println("💡 Hint: Network error connecting to generativelanguage.googleapis.com. Check your internet connection.")
		}
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("======================================================")
	fmt.Println("✅ SUCCESS: Gemini API key is valid and working!")
	fmt.Println("======================================================")
	fmt.Printf("Extracted Role: %s\n", entity.Role)
	fmt.Printf("Extracted City: %s\n", entity.City)
	fmt.Println("The backend LLM entity extraction pipeline is ready to use.")
}
