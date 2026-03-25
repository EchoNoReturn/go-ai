package llm

import "testing"

func TestNormalizeResponseBody_TrimsAndBom(t *testing.T) {
	input := []byte(" \n\t\xEF\xBB\xBF{\"ok\":true}\n")
	got := string(normalizeResponseBody(input))
	want := "{\"ok\":true}"
	if got != want {
		t.Fatalf("normalizeResponseBody mismatch: got=%q want=%q", got, want)
	}
}

func TestDecodeResponse_EnvelopeData(t *testing.T) {
	body := []byte(`{"data":{"id":"abc","object":"chat.completion","created":1,"model":"gpt","choices":[]}}`)
	resp, err := decodeResponse(OpenAI, body)
	if err != nil {
		t.Fatalf("decodeResponse error: %v", err)
	}
	result, ok := resp.(*OpenAIChatResponseBody)
	if !ok {
		t.Fatalf("unexpected response type: %T", resp)
	}
	if result.ID != "abc" {
		t.Fatalf("unexpected id: %q", result.ID)
	}
}

func TestDecodeResponse_ParseErrorIncludesProviderAndPreview(t *testing.T) {
	body := []byte("{bad json}")
	_, err := decodeResponse(OpenAI, body)
	if err == nil {
		t.Fatalf("expected parse error")
	}
	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("expected ParseError, got %T", err)
	}
	if parseErr.Provider != "openai" {
		t.Fatalf("unexpected provider: %q", parseErr.Provider)
	}
	if parseErr.Preview == "" {
		t.Fatalf("expected preview to be set")
	}
}
