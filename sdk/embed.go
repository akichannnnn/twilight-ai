package sdk

import (
	"context"
	"fmt"
)

type embedConfig struct {
	Params EmbedParams
}

// EmbedOption configures an embedding request.
type EmbedOption func(*embedConfig)

func WithEmbeddingModel(model *EmbeddingModel) EmbedOption {
	return func(c *embedConfig) { c.Params.Model = model }
}

func WithDimensions(d int) EmbedOption {
	return func(c *embedConfig) { c.Params.Dimensions = &d }
}

func buildEmbedConfig(values []string, options []EmbedOption) (*embedConfig, EmbeddingProvider, error) {
	cfg := &embedConfig{}
	for _, opt := range options {
		opt(cfg)
	}
	cfg.Params.Values = values

	if cfg.Params.Model == nil {
		return nil, nil, fmt.Errorf("twilightai: embedding model is required (use WithEmbeddingModel)")
	}
	if cfg.Params.Model.Provider == nil {
		return nil, nil, fmt.Errorf("twilightai: embedding model %q has no provider", cfg.Params.Model.ID)
	}
	return cfg, cfg.Params.Model.Provider, nil
}

// EmbedMany generates embeddings for multiple values.
func (c *Client) EmbedMany(ctx context.Context, values []string, options ...EmbedOption) (*EmbedResult, error) {
	cfg, prov, err := buildEmbedConfig(values, options)
	if err != nil {
		return nil, err
	}
	return prov.DoEmbed(ctx, cfg.Params)
}

// Embed generates an embedding for a single value.
// It returns the first embedding vector from the result.
func (c *Client) Embed(ctx context.Context, value string, options ...EmbedOption) ([]float64, error) {
	result, err := c.EmbedMany(ctx, []string{value}, options...)
	if err != nil {
		return nil, err
	}
	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("twilightai: no embedding returned")
	}
	return result.Embeddings[0], nil
}
