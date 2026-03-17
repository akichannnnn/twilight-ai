package types_test

import (
	"encoding/json"
	"testing"

	"github.com/memohai/twilight-ai/types"
)

func TestMessage_JSON_TextOnly(t *testing.T) {
	msg := types.Message{
		Role:    types.MessageRoleUser,
		Content: []types.MessagePart{types.TextPart{Text: "Hello"}},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	t.Logf("json: %s", data)

	var got types.Message
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Role != types.MessageRoleUser {
		t.Errorf("role: got %q, want %q", got.Role, types.MessageRoleUser)
	}
	if len(got.Content) != 1 {
		t.Fatalf("parts: got %d, want 1", len(got.Content))
	}
	tp, ok := got.Content[0].(types.TextPart)
	if !ok {
		t.Fatalf("part type: got %T, want TextPart", got.Content[0])
	}
	if tp.Text != "Hello" {
		t.Errorf("text: got %q, want %q", tp.Text, "Hello")
	}
}

func TestMessage_JSON_MultiPart(t *testing.T) {
	msg := types.Message{
		Role: types.MessageRoleUser,
		Content: []types.MessagePart{
			types.TextPart{Text: "Describe this image"},
			types.ImagePart{Image: "https://example.com/cat.png", MediaType: "image/png"},
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	t.Logf("json: %s", data)

	var got types.Message
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(got.Content) != 2 {
		t.Fatalf("parts: got %d, want 2", len(got.Content))
	}
	if tp, ok := got.Content[0].(types.TextPart); !ok || tp.Text != "Describe this image" {
		t.Errorf("part[0]: got %+v", got.Content[0])
	}
	if ip, ok := got.Content[1].(types.ImagePart); !ok || ip.Image != "https://example.com/cat.png" || ip.MediaType != "image/png" {
		t.Errorf("part[1]: got %+v", got.Content[1])
	}
}

func TestMessage_JSON_AllPartTypes(t *testing.T) {
	msg := types.Message{
		Role: types.MessageRoleAssistant,
		Content: []types.MessagePart{
			types.TextPart{Text: "answer"},
			types.ReasoningPart{Text: "thinking", Signature: "sig123"},
			types.ImagePart{Image: "data:image/png;base64,abc", MediaType: "image/png"},
			types.FilePart{Data: "base64data", MediaType: "application/pdf", Filename: "doc.pdf"},
			types.ToolCallPart{ToolCallID: "tc1", ToolName: "search", Input: map[string]any{"q": "go"}},
			types.ToolResultPart{ToolCallID: "tc1", ToolName: "search", Result: "found", IsError: false},
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	t.Logf("json: %s", data)

	var got types.Message
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(got.Content) != 6 {
		t.Fatalf("parts: got %d, want 6", len(got.Content))
	}

	expectTypes := []types.MessagePartType{
		types.MessagePartTypeText,
		types.MessagePartTypeReasoning,
		types.MessagePartTypeImage,
		types.MessagePartTypeFile,
		types.MessagePartTypeToolCall,
		types.MessagePartTypeToolResult,
	}
	for i, want := range expectTypes {
		if got.Content[i].PartType() != want {
			t.Errorf("part[%d]: type got %q, want %q", i, got.Content[i].PartType(), want)
		}
	}

	rp := got.Content[1].(types.ReasoningPart)
	if rp.Text != "thinking" || rp.Signature != "sig123" {
		t.Errorf("reasoning: got %+v", rp)
	}

	fp := got.Content[3].(types.FilePart)
	if fp.Data != "base64data" || fp.Filename != "doc.pdf" {
		t.Errorf("file: got %+v", fp)
	}

	tcp := got.Content[4].(types.ToolCallPart)
	if tcp.ToolCallID != "tc1" || tcp.ToolName != "search" {
		t.Errorf("tool call: got %+v", tcp)
	}

	trp := got.Content[5].(types.ToolResultPart)
	if trp.ToolCallID != "tc1" || trp.Result != "found" {
		t.Errorf("tool result: got %+v", trp)
	}
}

func TestMessage_JSON_FromRawJSON(t *testing.T) {
	raw := `{
		"role": "user",
		"content": [
			{"type": "text", "text": "What is this?"},
			{"type": "image", "image": "https://example.com/photo.jpg", "mediaType": "image/jpeg"}
		]
	}`

	var msg types.Message
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if msg.Role != types.MessageRoleUser {
		t.Errorf("role: got %q", msg.Role)
	}
	if len(msg.Content) != 2 {
		t.Fatalf("parts: got %d", len(msg.Content))
	}
	if tp, ok := msg.Content[0].(types.TextPart); !ok || tp.Text != "What is this?" {
		t.Errorf("part[0]: %+v", msg.Content[0])
	}
	if ip, ok := msg.Content[1].(types.ImagePart); !ok || ip.Image != "https://example.com/photo.jpg" {
		t.Errorf("part[1]: %+v", msg.Content[1])
	}
}

func TestUsage_JSON(t *testing.T) {
	u := types.Usage{
		InputTokens:  10,
		OutputTokens: 20,
		TotalTokens:  30,
		InputTokenDetails: types.InputTokenDetail{
			CacheReadTokens: 5,
		},
		OutputTokenDetails: types.OutputTokenDetail{
			ReasoningTokens: 8,
			TextTokens:      12,
		},
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	t.Logf("json: %s", data)

	var got types.Usage
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.InputTokens != 10 || got.OutputTokens != 20 || got.TotalTokens != 30 {
		t.Errorf("tokens: %+v", got)
	}
	if got.InputTokenDetails.CacheReadTokens != 5 {
		t.Errorf("cache read: got %d", got.InputTokenDetails.CacheReadTokens)
	}
	if got.OutputTokenDetails.ReasoningTokens != 8 {
		t.Errorf("reasoning: got %d", got.OutputTokenDetails.ReasoningTokens)
	}
}

func TestGenerateResult_JSON(t *testing.T) {
	r := types.GenerateResult{
		Text:         "Hello world",
		FinishReason: types.FinishReasonStop,
		Usage: types.Usage{
			InputTokens:  5,
			OutputTokens: 2,
			TotalTokens:  7,
		},
		ToolCalls: []types.ToolCall{{
			ToolCallID: "tc1",
			ToolName:   "search",
			Input:      map[string]any{"query": "go"},
		}},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	t.Logf("json: %s", data)

	var got types.GenerateResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Text != "Hello world" {
		t.Errorf("text: got %q", got.Text)
	}
	if got.FinishReason != types.FinishReasonStop {
		t.Errorf("finish: got %q", got.FinishReason)
	}
	if len(got.ToolCalls) != 1 || got.ToolCalls[0].ToolName != "search" {
		t.Errorf("tool calls: %+v", got.ToolCalls)
	}
}
