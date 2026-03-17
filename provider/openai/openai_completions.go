package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/memohai/twilight-ai/internal/utils"
	"github.com/memohai/twilight-ai/types"
)

const defaultBaseURL = "https://api.openai.com/v1"

type OpenAICompletionsProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type OpenAICompletionsProviderOption func(*OpenAICompletionsProvider)

func WithAPIKey(apiKey string) OpenAICompletionsProviderOption {
	return func(p *OpenAICompletionsProvider) {
		p.apiKey = apiKey
	}
}

func WithBaseURL(baseURL string) OpenAICompletionsProviderOption {
	return func(p *OpenAICompletionsProvider) {
		p.baseURL = baseURL
	}
}

func WithHTTPClient(client *http.Client) OpenAICompletionsProviderOption {
	return func(p *OpenAICompletionsProvider) {
		p.httpClient = client
	}
}

func NewCompletions(options ...OpenAICompletionsProviderOption) *OpenAICompletionsProvider {
	provider := &OpenAICompletionsProvider{
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
	}
	for _, option := range options {
		option(provider)
	}
	return provider
}

func (p *OpenAICompletionsProvider) Name() string {
	return "openai-completions"
}

func (p *OpenAICompletionsProvider) GetModels() ([]types.Model, error) {
	return nil, nil
}

// ChatModel creates a Model bound to this provider.
func (p *OpenAICompletionsProvider) ChatModel(id string) *types.Model {
	return &types.Model{
		ID:       id,
		Provider: p,
		Type:     types.ModelTypeChat,
	}
}

// ---------- DoGenerate ----------

func (p *OpenAICompletionsProvider) DoGenerate(ctx context.Context, params types.GenerateParams) (*types.GenerateResult, error) {
	if params.Model == nil {
		return nil, fmt.Errorf("openai: model is required")
	}

	req := p.buildRequest(params)

	resp, err := utils.FetchJSON[chatResponse](ctx, p.httpClient, utils.RequestOptions{
		Method:  http.MethodPost,
		BaseURL: p.baseURL,
		Path:    "/chat/completions",
		Headers: utils.AuthHeader(p.apiKey),
		Body:    req,
	})
	if err != nil {
		return nil, fmt.Errorf("openai: chat completions request failed: %w", err)
	}

	return p.parseResponse(resp), nil
}

// ---------- buildRequest ----------

func (p *OpenAICompletionsProvider) buildRequest(params types.GenerateParams) *chatRequest {
	req := &chatRequest{
		Model:               params.Model.ID,
		Messages:            convertMessages(params),
		Temperature:         params.Temperature,
		TopP:                params.TopP,
		MaxCompletionTokens: params.MaxTokens,
		FrequencyPenalty:    params.FrequencyPenalty,
		PresencePenalty:     params.PresencePenalty,
		Seed:                params.Seed,
		ReasoningEffort:     params.ReasoningEffort,
	}
	if len(params.StopSequences) > 0 {
		req.Stop = params.StopSequences
	}
	if len(params.Tools) > 0 {
		req.Tools = convertTools(params.Tools)
		req.ToolChoice = params.ToolChoice
	}
	if params.ResponseFormat != nil {
		req.ResponseFormat = &chatRespFormat{
			Type:       string(params.ResponseFormat.Type),
			JSONSchema: params.ResponseFormat.JSONSchema,
		}
	}
	return req
}

func convertTools(tools []types.Tool) []chatTool {
	out := make([]chatTool, 0, len(tools))
	for _, t := range tools {
		out = append(out, chatTool{
			Type: "function",
			Function: chatFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		})
	}
	return out
}

// ---------- message conversion ----------

func convertMessages(params types.GenerateParams) []chatMessage {
	var out []chatMessage

	if params.System != "" {
		out = append(out, chatMessage{
			Role:    "system",
			Content: params.System,
		})
	}

	for _, msg := range params.Messages {
		out = append(out, convertMessage(msg)...)
	}
	return out
}

func convertMessage(msg types.Message) []chatMessage {
	switch msg.Role {
	case types.MessageRoleTool:
		return convertToolResultMessages(msg)
	case types.MessageRoleAssistant:
		return []chatMessage{convertAssistantMessage(msg)}
	default:
		return []chatMessage{{
			Role:    string(msg.Role),
			Content: convertContent(msg.Content),
		}}
	}
}

func convertAssistantMessage(msg types.Message) chatMessage {
	cm := chatMessage{Role: "assistant"}

	var contentParts []types.MessagePart
	var toolCalls []chatToolCall

	for _, part := range msg.Content {
		switch p := part.(type) {
		case types.ToolCallPart:
			args, _ := json.Marshal(p.Input)
			toolCalls = append(toolCalls, chatToolCall{
				ID:   p.ToolCallID,
				Type: "function",
				Function: chatFunctionCall{
					Name:      p.ToolName,
					Arguments: string(args),
				},
			})
		default:
			contentParts = append(contentParts, part)
		}
	}

	if len(contentParts) > 0 {
		cm.Content = convertContent(contentParts)
	}
	if len(toolCalls) > 0 {
		cm.ToolCalls = toolCalls
	}

	return cm
}

func convertToolResultMessages(msg types.Message) []chatMessage {
	var out []chatMessage
	for _, part := range msg.Content {
		if trp, ok := part.(types.ToolResultPart); ok {
			content, _ := json.Marshal(trp.Result)
			out = append(out, chatMessage{
				Role:       "tool",
				ToolCallID: trp.ToolCallID,
				Content:    string(content),
			})
		}
	}
	return out
}

func convertContent(parts []types.MessagePart) any {
	if len(parts) == 1 {
		if tp, ok := parts[0].(types.TextPart); ok {
			return tp.Text
		}
	}

	out := make([]any, 0, len(parts))
	for _, part := range parts {
		switch p := part.(type) {
		case types.TextPart:
			out = append(out, chatContentPartText{Type: "text", Text: p.Text})
		case types.ReasoningPart:
			out = append(out, chatContentPartText{Type: "text", Text: p.Text})
		case types.ImagePart:
			out = append(out, chatContentPartImage{
				Type:     "image_url",
				ImageURL: chatImageURL{URL: p.Image},
			})
		case types.FilePart:
			out = append(out, chatContentPartText{Type: "text", Text: p.Data})
		}
	}
	return out
}

// ---------- parseResponse ----------

func (p *OpenAICompletionsProvider) parseResponse(resp *chatResponse) *types.GenerateResult {
	result := &types.GenerateResult{
		Usage: convertUsage(&resp.Usage),
		Response: types.ResponseMetadata{
			ID:        resp.ID,
			ModelID:   resp.Model,
			Timestamp: time.Unix(resp.Created, 0),
		},
	}

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		result.Text = choice.Message.Content
		result.Reasoning = choice.Message.ReasoningContent
		result.FinishReason = mapFinishReason(choice.FinishReason)
		result.RawFinishReason = choice.FinishReason

		for _, tc := range choice.Message.ToolCalls {
			var input any
			json.Unmarshal([]byte(tc.Function.Arguments), &input)
			result.ToolCalls = append(result.ToolCalls, types.ToolCall{
				ToolCallID: tc.ID,
				ToolName:   tc.Function.Name,
				Input:      input,
			})
		}
	}

	return result
}

// ---------- DoStream ----------

func (p *OpenAICompletionsProvider) DoStream(ctx context.Context, params types.GenerateParams) (*types.StreamResult, error) {
	if params.Model == nil {
		return nil, fmt.Errorf("openai: model is required")
	}

	req := p.buildRequest(params)
	req.Stream = true
	req.StreamOptions = &chatStreamOptions{IncludeUsage: true}

	ch := make(chan types.StreamPart, 64)

	go func() {
		defer close(ch)

		var (
			textStartSent      bool
			reasoningStartSent bool
			rawFinishReason    string
			finishReason       types.FinishReason
			usage              types.Usage
			chunkID            string
			chunkModel         string
			chunkCreated       int64
			pendingToolCalls   = map[int]*streamingToolCall{}
		)

		send := func(part types.StreamPart) bool {
			select {
			case ch <- part:
				return true
			case <-ctx.Done():
				return false
			}
		}

		if !send(&types.StartPart{}) {
			return
		}
		if !send(&types.StartStepPart{}) {
			return
		}

		err := utils.FetchSSE(ctx, p.httpClient, utils.RequestOptions{
			Method:  http.MethodPost,
			BaseURL: p.baseURL,
			Path:    "/chat/completions",
			Headers: utils.AuthHeader(p.apiKey),
			Body:    req,
		}, func(ev *utils.SSEEvent) error {
			if ev.Data == "[DONE]" {
				return utils.ErrStreamDone
			}

			var chunk chatChunkResponse
			if err := json.Unmarshal([]byte(ev.Data), &chunk); err != nil {
				send(&types.ErrorPart{Error: fmt.Errorf("openai: unmarshal chunk: %w", err)})
				return err
			}

			if chunkID == "" {
				chunkID = chunk.ID
				chunkModel = chunk.Model
				chunkCreated = chunk.Created
			}

			if chunk.Usage != nil {
				usage = convertUsage(chunk.Usage)
			}

			if len(chunk.Choices) == 0 {
				return nil
			}
			choice := chunk.Choices[0]

			// reasoning content (e.g. DeepSeek, o1-compatible providers)
			if choice.Delta.ReasoningContent != "" {
				if !reasoningStartSent {
					send(&types.ReasoningStartPart{ID: chunk.ID})
					reasoningStartSent = true
				}
				send(&types.ReasoningDeltaPart{ID: chunk.ID, Text: choice.Delta.ReasoningContent})
			}

			// text content
			if choice.Delta.Content != "" {
				if reasoningStartSent {
					send(&types.ReasoningEndPart{ID: chunk.ID})
					reasoningStartSent = false
				}
				if !textStartSent {
					send(&types.TextStartPart{ID: chunk.ID})
					textStartSent = true
				}
				send(&types.TextDeltaPart{ID: chunk.ID, Text: choice.Delta.Content})
			}

			// tool call deltas
			for _, tc := range choice.Delta.ToolCalls {
				stc, exists := pendingToolCalls[tc.Index]
				if !exists {
					stc = &streamingToolCall{}
					pendingToolCalls[tc.Index] = stc
					stc.id = tc.ID
					stc.name = tc.Function.Name
					send(&types.ToolInputStartPart{
						ID:       tc.ID,
						ToolName: tc.Function.Name,
					})
				}
				if tc.Function.Arguments != "" {
					stc.args += tc.Function.Arguments
					send(&types.ToolInputDeltaPart{
						ID:    stc.id,
						Delta: tc.Function.Arguments,
					})
				}
			}

			// finish
			if choice.FinishReason != nil && *choice.FinishReason != "" {
				rawFinishReason = *choice.FinishReason
				finishReason = mapFinishReason(rawFinishReason)

				if reasoningStartSent {
					send(&types.ReasoningEndPart{ID: chunk.ID})
				}
				if textStartSent {
					send(&types.TextEndPart{ID: chunk.ID})
				}

				for _, stc := range pendingToolCalls {
					send(&types.ToolInputEndPart{ID: stc.id})
					var input any
					json.Unmarshal([]byte(stc.args), &input)
					send(&types.StreamToolCallPart{
						ToolCallID: stc.id,
						ToolName:   stc.name,
						Input:      input,
					})
				}

				send(&types.FinishStepPart{
					FinishReason:    finishReason,
					RawFinishReason: rawFinishReason,
					Usage:           usage,
					Response: types.ResponseMetadata{
						ID:        chunkID,
						ModelID:   chunkModel,
						Timestamp: time.Unix(chunkCreated, 0),
					},
				})
			}

			return nil
		})

		if err != nil {
			send(&types.ErrorPart{Error: fmt.Errorf("openai: stream failed: %w", err)})
		}

		send(&types.FinishPart{
			FinishReason:    finishReason,
			RawFinishReason: rawFinishReason,
			TotalUsage:      usage,
		})
	}()

	return &types.StreamResult{Stream: ch}, nil
}

type streamingToolCall struct {
	id   string
	name string
	args string
}

// ---------- helpers ----------

func convertUsage(u *chatUsage) types.Usage {
	usage := types.Usage{
		InputTokens:  u.PromptTokens,
		OutputTokens: u.CompletionTokens,
		TotalTokens:  u.TotalTokens,
	}
	if u.PromptTokensDetails != nil {
		usage.CachedInputTokens = u.PromptTokensDetails.CachedTokens
		usage.InputTokenDetails.CacheReadTokens = u.PromptTokensDetails.CachedTokens
	}
	if u.CompletionTokensDetails != nil {
		usage.ReasoningTokens = u.CompletionTokensDetails.ReasoningTokens
		usage.OutputTokenDetails.ReasoningTokens = u.CompletionTokensDetails.ReasoningTokens
		usage.OutputTokenDetails.TextTokens = u.CompletionTokensDetails.TextTokens
	}
	return usage
}

func mapFinishReason(reason string) types.FinishReason {
	switch reason {
	case "stop":
		return types.FinishReasonStop
	case "length":
		return types.FinishReasonLength
	case "content_filter":
		return types.FinishReasonContentFilter
	case "tool_calls":
		return types.FinishReasonToolCalls
	default:
		return types.FinishReasonUnknown
	}
}
