package llm

import (
	"encoding/json"
	"testing"
)

func mustUnmarshalAnthropicResponse(tb testing.TB, payload string) AnthropicChatResponse {
	tb.Helper()
	var resp AnthropicChatResponse
	if err := json.Unmarshal([]byte(payload), &resp); err != nil {
		tb.Fatalf("unmarshal anthropic response: %v", err)
	}
	return resp
}

func TestAnthropicContent_String(t *testing.T) {
	resp := mustUnmarshalAnthropicResponse(t, `{"id":"abc","type":"message","model":"claude","content":"hello"}`)
	if resp.Content.Text != "hello" {
		t.Fatalf("unexpected content text: %q", resp.Content.Text)
	}
	if len(resp.Content.Blocks) != 0 {
		t.Fatalf("expected no content blocks")
	}
}

func TestAnthropicContent_Blocks(t *testing.T) {
	resp := mustUnmarshalAnthropicResponse(t, `{"id":"abc","type":"message","model":"claude","content":[{"type":"text","text":"hi"},{"type":"text","text":" there"}]}`)
	if resp.Content.Text != "hi there" {
		t.Fatalf("unexpected content text: %q", resp.Content.Text)
	}
	if len(resp.Content.Blocks) != 2 {
		t.Fatalf("expected 2 content blocks, got %d", len(resp.Content.Blocks))
	}
}

func TestAnthropicContent_MixedBlocks(t *testing.T) {
	resp := mustUnmarshalAnthropicResponse(t, `{"id":"abc","type":"message","model":"claude","content":[{"type":"text","text":"hello"},{"type":"image","text":"ignored"},{"type":"text","text":" world"}]}`)
	if resp.Content.Text != "hello world" {
		t.Fatalf("unexpected content text: %q", resp.Content.Text)
	}
	if len(resp.Content.Blocks) != 3 {
		t.Fatalf("expected 3 content blocks, got %d", len(resp.Content.Blocks))
	}
}

func TestAnthropicStreamContent_Blocks(t *testing.T) {
	var resp AnthropicChatStreamResponse
	payload := `{"id":"abc","type":"message","model":"claude","content":[{"type":"text","text":"hi"}]}`
	if err := json.Unmarshal([]byte(payload), &resp); err != nil {
		t.Fatalf("unmarshal anthropic stream response: %v", err)
	}
	if resp.Content.Text != "hi" {
		t.Fatalf("unexpected stream content text: %q", resp.Content.Text)
	}
	if len(resp.Content.Blocks) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(resp.Content.Blocks))
	}
}
