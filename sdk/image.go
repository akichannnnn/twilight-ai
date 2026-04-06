package sdk

import "context"

// ImageGenerationProvider is the interface that image generation backends must implement.
type ImageGenerationProvider interface {
	DoGenerate(ctx context.Context, params *ImageGenerationParams) (*ImageResult, error)
}

// ImageEditProvider is the interface that image editing backends must implement.
type ImageEditProvider interface {
	DoEdit(ctx context.Context, params *ImageEditParams) (*ImageResult, error)
}

// ImageGenerationModel represents an image generation model bound to an ImageGenerationProvider.
type ImageGenerationModel struct {
	ID       string
	Provider ImageGenerationProvider
}

// ImageEditModel represents an image edit model bound to an ImageEditProvider.
type ImageEditModel struct {
	ID       string
	Provider ImageEditProvider
}

// ImageGenerationParams holds the parameters for an image generation request.
type ImageGenerationParams struct {
	Model             *ImageGenerationModel
	Prompt            string
	N                 *int
	Size              string // e.g. "1024x1024", "1536x1024", "256x256"
	Quality           string // "auto", "low", "medium", "high", "standard", "hd"
	Style             string // dall-e-3 only: "vivid", "natural"
	ResponseFormat    string // dall-e-2/3: "url", "b64_json"
	Background        string // gpt-image: "transparent", "opaque", "auto"
	OutputFormat      string // gpt-image: "png", "jpeg", "webp"
	OutputCompression *int   // gpt-image, jpeg/webp only: 0-100
	Moderation        string // gpt-image: "low", "auto"
	User              string
}

// ImageEditParams holds the parameters for an image edit request.
type ImageEditParams struct {
	Model             *ImageEditModel
	Images            []ImageInput
	Prompt            string
	Mask              *ImageInput
	N                 *int
	Size              string
	Quality           string // gpt-image: "auto", "low", "medium", "high"
	Background        string // "transparent", "opaque", "auto"
	OutputFormat      string // gpt-image: "png", "jpeg", "webp"
	OutputCompression *int   // gpt-image, jpeg/webp only: 0-100
	InputFidelity     string // gpt-image: "high", "low"
	Moderation        string // gpt-image: "low", "auto"
	ResponseFormat    string // dall-e-2: "url", "b64_json"
	User              string
}

// ImageInput represents an image provided as input to an edit request.
// Exactly one of Data, URL, or FileID should be set.
type ImageInput struct {
	Data      []byte // raw file bytes (for multipart upload)
	MediaType string // e.g. "image/png"
	Filename  string
	URL       string // image URL reference
	FileID    string // OpenAI file ID reference
}

// ImageResult holds the result of an image generation or edit request.
type ImageResult struct {
	Created int64
	Data    []ImageData
	Usage   ImageUsage
}

// ImageData represents a single generated or edited image.
type ImageData struct {
	B64JSON       string
	URL           string
	RevisedPrompt string // dall-e-3 only
}

// ImageUsage tracks token usage for image requests (GPT Image models).
type ImageUsage struct {
	TotalTokens       int
	InputTokens       int
	OutputTokens      int
	InputTokenDetails *ImageInputTokenDetails
}

// ImageInputTokenDetails breaks down input token usage by modality.
type ImageInputTokenDetails struct {
	TextTokens  int
	ImageTokens int
}
