package prompt

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
)

// init loads environment variables from a .env file when the package initializes.
// It first checks for the OPENAI_API_KEY. If not found, it prompts the user for it.
func init() {
	// Try to load .env file (if it exists).
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: no .env file found.")
	}

	// Check if the OPENAI_API_KEY is already set.
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		fmt.Println("OPENAI_API_KEY found in environment.")
		return
	}

	// If not set, prompt the user for the API key.
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your OpenAI API key: ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading API key:", err)
		return
	}
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		fmt.Println("No API key provided. OpenAI functionality may not work correctly.")
		return
	}

	// Write the provided API key to a .env file.
	if err := WriteEnvFile(apiKey); err != nil {
		fmt.Println("Error writing .env file:", err)
	} else {
		fmt.Println(".env file written successfully.")
		// Optionally, reload the environment variables so that the key is available.
		_ = godotenv.Load()
	}
}

// WriteEnvFile writes a .env file in the current directory with the provided OpenAI API key.
func WriteEnvFile(apiKey string) error {
	envContent := fmt.Sprintf("OPENAI_API_KEY=%s\n", apiKey)
	if err := os.WriteFile(".env", []byte(envContent), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %v", err)
	}
	return nil
}

// SaveOutputToFile saves the given response to output.md.
func SaveOutputToFile(response string) error {
	filePath := "output.md"
	err := os.WriteFile(filePath, []byte(response), 0644)
	if err != nil {
		return fmt.Errorf("failed to save response to output.md: %v", err)
	}
	return nil
}

type CodeResponse struct {
	Filename string   `json:"filename" jsonschema_description:"Name of the output file"`
	Scripts  []string `json:"scripts" jsonschema_description:"List of code scripts"`
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	return reflector.Reflect(v)
}

var codeResponseSchema = GenerateSchema[CodeResponse]()

// PromptOpenai sends the given input to OpenAI and expects a structured JSON response.
func PromptOpenai(input string) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("code_response"),
		Description: openai.F("Response containing a filename and scripts"),
		Schema:      openai.F(codeResponseSchema),
		Strict:      openai.Bool(true),
	}

	client := openai.NewClient()
	ctx := context.Background()

	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})
	if err != nil {
		fmt.Println("OpenAI request failed:", err)
		return
	}

	if err := SaveOutputToFile(chat.Choices[0].Message.Content); err != nil {
		fmt.Println("Error saving output:", err)
	} else {
		fmt.Println("\n📄 Response saved to output.md")
	}
}

// PromptFromNeovim opens input.md in Neovim, reads the content, and sends it to OpenAI.
func PromptFromNeovim() {
	input, err := OpenInputInNeovim()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PromptOpenai(input)
}

type FileResponse struct {
	Files []string `json:"files"`
}

// PromptOpenaiFiles sends the given input to OpenAI expecting a response with only a list of file names.
func PromptOpenaiFiles(input string) {
	codeResponseSchema := GenerateSchema[FileResponse]()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("code_response"),
		Description: openai.F("Response containing only a list of file names"),
		Schema:      openai.F(codeResponseSchema),
		Strict:      openai.Bool(true),
	}

	client := openai.NewClient()
	ctx := context.Background()

	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})
	if err != nil {
		fmt.Println("OpenAI request failed:", err)
		return
	}

	if err := SaveOutputToFile(chat.Choices[0].Message.Content); err != nil {
		fmt.Println("Error saving output:", err)
	} else {
		fmt.Println("\n📄 Response saved to output.md")
	}
}
