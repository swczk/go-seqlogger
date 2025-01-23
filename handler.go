package seqlogger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type SeqHandler struct {
	config Config
	client *http.Client
	attrs  []slog.Attr
}

func NewSeqHandler(config Config) *SeqHandler {
	return &SeqHandler{
		config: config,
		client: &http.Client{Timeout: config.ClientTimeout},
		attrs:  make([]slog.Attr, 0),
	}
}

func (h *SeqHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.config.LogLevel
}

func (h *SeqHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := &SeqHandler{
		config: h.config,
		client: h.client,
		attrs:  make([]slog.Attr, len(h.attrs)+len(attrs)),
	}
	copy(newHandler.attrs, h.attrs)
	copy(newHandler.attrs[len(h.attrs):], attrs)
	return newHandler
}

func (h *SeqHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *SeqHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	event := createSeqEvent(h.config, ctx, r, h.attrs)
	return h.sendToSeq(ctx, event)
}

func (h *SeqHandler) sendToSeq(ctx context.Context, event map[string]interface{}) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	jsonData = append(jsonData, '\n')

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		h.config.Endpoint+"/ingest/clef",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/vnd.serilog.clef")
	if h.config.APIKey != "" {
		req.Header.Set("X-Seq-ApiKey", h.config.APIKey)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("seq returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
