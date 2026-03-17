package sdk

import "context"

// EmbeddingProvider is the interface that embedding backends must implement.
type EmbeddingProvider interface {
	DoEmbed(ctx context.Context, params EmbedParams) (*EmbedResult, error)
}

// EmbeddingModel represents an embedding model bound to an EmbeddingProvider.
type EmbeddingModel struct {
	ID                   string
	Provider             EmbeddingProvider
	MaxEmbeddingsPerCall int
}

// EmbedParams holds the parameters for an embedding request.
type EmbedParams struct {
	Model      *EmbeddingModel
	Values     []string
	Dimensions *int
}

// EmbedResult holds the result of an embedding request.
type EmbedResult struct {
	Embeddings [][]float64
	Usage      EmbeddingUsage
}

// EmbeddingUsage tracks token usage for embedding requests.
type EmbeddingUsage struct {
	Tokens int
}
