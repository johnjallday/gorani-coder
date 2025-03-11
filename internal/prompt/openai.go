package prompt

import (
	"context"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

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
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var codeResponseSchema = GenerateSchema[CodeResponse]()

// PromptOpenai sends the given input to OpenAI and expects a structured JSON response.
// The response should follow this schema:
//
//	{
//	  "filename": "string",
//	  "scripts": ["string", ...]
//	}
func PromptOpenai(input string) {
	// Define the expected structured output.
	// Generate the JSON schema for CodeResponse.
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("code_response"),
		Description: openai.F("Response containing a filename and scripts"),
		Schema:      openai.F(codeResponseSchema),
		Strict:      openai.Bool(true),
	}

	// Set up the OpenAI client and context.
	client := openai.NewClient()
	ctx := context.Background()

	// Call the Chat Completions API with the structured response format.
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
		// Use a model that supports structured outputs.
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})
	if err != nil {
		fmt.Println("OpenAI request failed:", err)
		return
	}

	// Save the raw response to output.md.
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
	PromptOpenai(input) // Call PromptOpenai with the edited input.
}

type FileResponse struct {
	Files []string `json:"files"`
}

// promptOpenAIfor Files output
func PromptOpenaiFiles(input string) {
	// Generate JSON schema for FileResponse using the GenerateSchema function.
	codeResponseSchema := GenerateSchema[FileResponse]()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("code_response"),
		Description: openai.F("Response containing only a list of file names"),
		Schema:      openai.F(codeResponseSchema),
		Strict:      openai.Bool(true),
	}

	// Set up the OpenAI client and context.
	client := openai.NewClient()
	ctx := context.Background()

	// Call the Chat Completions API with the structured response format.
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
		// Use a model that supports structured outputs.
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})
	if err != nil {
		fmt.Println("OpenAI request failed:", err)
		return
	}

	// Save the raw response to output.md.
	if err := SaveOutputToFile(chat.Choices[0].Message.Content); err != nil {
		fmt.Println("Error saving output:", err)
	} else {
		fmt.Println("\n📄 Response saved to output.md")
	}
}
