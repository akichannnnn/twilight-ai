package embedding

import (
	"context"
	"fmt"
	"net/http"

	"github.com/memohai/twilight-ai/internal/utils"
	"github.com/memohai/twilight-ai/sdk"
)

const defaultBaseURL = "https://api.openai.com/v1"

type Provider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type Option func(*Provider)

func WithAPIKey(apiKey string) Option {
	return func(p *Provider) { p.apiKey = apiKey }
}

func WithBaseURL(baseURL string) Option {
	return func(p *Provider) { p.baseURL = baseURL }
}

func WithHTTPClient(client *http.Client) Option {
	return func(p *Provider) { p.httpClient = client }
}

func New(options ...Option) *Provider {
	p := &Provider{
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}

// EmbeddingModel creates an EmbeddingModel bound to this provider.
func (p *Provider) EmbeddingModel(id string) *sdk.EmbeddingModel {
	return &sdk.EmbeddingModel{
		ID:                   id,
		Provider:             p,
		MaxEmbeddingsPerCall: 2048,
	}
}

// DoEmbed implements sdk.EmbeddingProvider.
func (p *Provider) DoEmbed(ctx context.Context, params sdk.EmbedParams) (*sdk.EmbedResult, error) {
	if params.Model == nil {
		return nil, fmt.Errorf("openai: embedding model is required")
	}

	req := &embeddingRequest{
		Model:          params.Model.ID,
		Input:          params.Values,
		EncodingFormat: "float",
		Dimensions:     params.Dimensions,
	}

	resp, err := utils.FetchJSON[embeddingResponse](ctx, p.httpClient, &utils.RequestOptions{
		Method:  http.MethodPost,
		BaseURL: p.baseURL,
		Path:    "/embeddings",
		Headers: utils.AuthHeader(p.apiKey),
		Body:    req,
	})
	if err != nil {
		return nil, fmt.Errorf("openai: embeddings request failed: %w", err)
	}

	embeddings := make([][]float64, len(resp.Data))
	for i, d := range resp.Data {
		embeddings[i] = d.Embedding
	}

	return &sdk.EmbedResult{
		Embeddings: embeddings,
		Usage: sdk.EmbeddingUsage{
			Tokens: resp.Usage.PromptTokens,
		},
	}, nil
}
