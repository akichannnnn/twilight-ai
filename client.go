package twilightai

import (
	"context"
	"fmt"

	"github.com/memohai/twilight-ai/types"
)

// Client provides text generation methods.
// The provider is resolved from the Model passed via WithModel.
type Client struct{}

func NewClient() *Client {
	return &Client{}
}

// GenerateOption configures a text generation request.
type GenerateOption func(*types.GenerateParams)

// WithModel sets the model (which carries its Provider) for the request.
func WithModel(model *types.Model) GenerateOption {
	return func(p *types.GenerateParams) {
		p.Model = model
	}
}

func WithMessages(messages []types.Message) GenerateOption {
	return func(p *types.GenerateParams) {
		p.Messages = messages
	}
}

func WithSystem(text string) GenerateOption {
	return func(p *types.GenerateParams) {
		p.System = text
	}
}

func WithTools(tools []types.Tool) GenerateOption {
	return func(p *types.GenerateParams) {
		p.Tools = tools
	}
}

func WithToolChoice(choice any) GenerateOption {
	return func(p *types.GenerateParams) {
		p.ToolChoice = choice
	}
}

func WithResponseFormat(rf types.ResponseFormat) GenerateOption {
	return func(p *types.GenerateParams) {
		p.ResponseFormat = &rf
	}
}

func WithTemperature(t float64) GenerateOption {
	return func(p *types.GenerateParams) {
		p.Temperature = &t
	}
}

func WithTopP(topP float64) GenerateOption {
	return func(p *types.GenerateParams) {
		p.TopP = &topP
	}
}

func WithMaxTokens(n int) GenerateOption {
	return func(p *types.GenerateParams) {
		p.MaxTokens = &n
	}
}

func WithStopSequences(s []string) GenerateOption {
	return func(p *types.GenerateParams) {
		p.StopSequences = s
	}
}

func WithFrequencyPenalty(penalty float64) GenerateOption {
	return func(p *types.GenerateParams) {
		p.FrequencyPenalty = &penalty
	}
}

func WithPresencePenalty(penalty float64) GenerateOption {
	return func(p *types.GenerateParams) {
		p.PresencePenalty = &penalty
	}
}

func WithSeed(s int) GenerateOption {
	return func(p *types.GenerateParams) {
		p.Seed = &s
	}
}

func WithReasoningEffort(effort string) GenerateOption {
	return func(p *types.GenerateParams) {
		p.ReasoningEffort = &effort
	}
}

func buildParams(options []GenerateOption) (types.GenerateParams, types.Provider, error) {
	var params types.GenerateParams
	for _, opt := range options {
		opt(&params)
	}

	if params.Model == nil {
		return params, nil, fmt.Errorf("twilightai: model is required (use WithModel)")
	}
	if params.Model.Provider == nil {
		return params, nil, fmt.Errorf("twilightai: model %q has no provider", params.Model.ID)
	}

	return params, params.Model.Provider, nil
}

// GenerateText returns only the generated text content.
func (c *Client) GenerateText(ctx context.Context, options ...GenerateOption) (string, error) {
	params, prov, err := buildParams(options)
	if err != nil {
		return "", err
	}
	result, err := prov.DoGenerate(ctx, params)
	if err != nil {
		return "", err
	}
	return result.Text, nil
}

// GenerateTextResult returns the full generation result.
func (c *Client) GenerateTextResult(ctx context.Context, options ...GenerateOption) (*types.GenerateResult, error) {
	params, prov, err := buildParams(options)
	if err != nil {
		return nil, err
	}
	return prov.DoGenerate(ctx, params)
}

// StreamText returns a streaming result with a channel of StreamPart events.
func (c *Client) StreamText(ctx context.Context, options ...GenerateOption) (*types.StreamResult, error) {
	params, prov, err := buildParams(options)
	if err != nil {
		return nil, err
	}
	return prov.DoStream(ctx, params)
}

// --- Package-level convenience functions ---

// GenerateText is a package-level shortcut for text generation.
func GenerateText(ctx context.Context, options ...GenerateOption) (string, error) {
	return defaultClient.GenerateText(ctx, options...)
}

// GenerateTextResult is a package-level shortcut returning the full result.
func GenerateTextResult(ctx context.Context, options ...GenerateOption) (*types.GenerateResult, error) {
	return defaultClient.GenerateTextResult(ctx, options...)
}

// StreamText is a package-level shortcut for streaming generation.
func StreamText(ctx context.Context, options ...GenerateOption) (*types.StreamResult, error) {
	return defaultClient.StreamText(ctx, options...)
}

var defaultClient = &Client{}
