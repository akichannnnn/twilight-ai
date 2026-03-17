package twilightai_test

import (
	"context"
	"os"
	"testing"

	twilightai "github.com/memohai/twilight-ai"
	"github.com/memohai/twilight-ai/internal/testutil"
	"github.com/memohai/twilight-ai/provider/openai"
	"github.com/memohai/twilight-ai/types"
)

func TestMain(m *testing.M) {
	testutil.LoadEnv()
	os.Exit(m.Run())
}

func envOrSkip(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("skipping: %s not set", key)
	}
	return v
}

func newProvider(t *testing.T) *openai.OpenAICompletionsProvider {
	t.Helper()
	apiKey := envOrSkip(t, "OPENAI_API_KEY")
	opts := []openai.OpenAICompletionsProviderOption{openai.WithAPIKey(apiKey)}
	if base := os.Getenv("OPENAI_BASE_URL"); base != "" {
		opts = append(opts, openai.WithBaseURL(base))
	}
	return openai.NewCompletions(opts...)
}

func model(t *testing.T) *types.Model {
	t.Helper()
	id := os.Getenv("OPENAI_MODEL")
	if id == "" {
		id = "gpt-4o-mini"
	}
	return newProvider(t).ChatModel(id)
}

func TestClient_GenerateText(t *testing.T) {
	text, err := twilightai.GenerateText(context.Background(),
		twilightai.WithModel(model(t)),
		twilightai.WithMessages([]types.Message{
			types.UserMessage("Say hi in one word."),
		}),
	)
	if err != nil {
		t.Fatalf("GenerateText: %v", err)
	}
	t.Logf("response: %q", text)
	if text == "" {
		t.Error("expected non-empty response")
	}
}

func TestClient_GenerateTextResult(t *testing.T) {
	result, err := twilightai.GenerateTextResult(context.Background(),
		twilightai.WithModel(model(t)),
		twilightai.WithMessages([]types.Message{
			types.UserMessage("Say hi in one word."),
		}),
	)
	if err != nil {
		t.Fatalf("GenerateTextResult: %v", err)
	}
	t.Logf("text=%q finish=%s input=%d output=%d",
		result.Text, result.FinishReason,
		result.Usage.InputTokens, result.Usage.OutputTokens)

	if result.Text == "" {
		t.Error("expected non-empty text")
	}
	if result.FinishReason != types.FinishReasonStop {
		t.Errorf("expected stop, got %s", result.FinishReason)
	}
}

func TestClient_StreamText(t *testing.T) {
	sr, err := twilightai.StreamText(context.Background(),
		twilightai.WithModel(model(t)),
		twilightai.WithMessages([]types.Message{
			types.UserMessage("Count from 1 to 3."),
		}),
	)
	if err != nil {
		t.Fatalf("StreamText: %v", err)
	}

	var text string
	for part := range sr.Stream {
		switch p := part.(type) {
		case *types.TextDeltaPart:
			text += p.Text
		case *types.ErrorPart:
			t.Fatalf("stream error: %v", p.Error)
		case *types.FinishPart:
			t.Logf("finish=%s tokens=%d", p.FinishReason, p.TotalUsage.TotalTokens)
		}
	}
	t.Logf("streamed: %q", text)
	if text == "" {
		t.Error("expected non-empty streamed text")
	}
}

func TestClient_StreamText_ToResult(t *testing.T) {
	sr, err := twilightai.StreamText(context.Background(),
		twilightai.WithModel(model(t)),
		twilightai.WithMessages([]types.Message{
			types.UserMessage("Say hello in one word."),
		}),
	)
	if err != nil {
		t.Fatalf("StreamText: %v", err)
	}

	result, err := sr.ToResult()
	if err != nil {
		t.Fatalf("ToResult: %v", err)
	}
	t.Logf("text=%q finish=%s", result.Text, result.FinishReason)
	if result.Text == "" {
		t.Error("expected non-empty text")
	}
}

func TestClient_WithSystem(t *testing.T) {
	text, err := twilightai.GenerateText(context.Background(),
		twilightai.WithModel(model(t)),
		twilightai.WithSystem("You always respond with exactly one word."),
		twilightai.WithMessages([]types.Message{
			types.UserMessage("Greet me."),
		}),
	)
	if err != nil {
		t.Fatalf("GenerateText: %v", err)
	}
	t.Logf("response: %q", text)
	if text == "" {
		t.Error("expected non-empty response")
	}
}

func TestClient_NoModel(t *testing.T) {
	_, err := twilightai.GenerateText(context.Background(),
		twilightai.WithMessages([]types.Message{
			types.UserMessage("Hi"),
		}),
	)
	if err == nil {
		t.Fatal("expected error for nil model")
	}
}
