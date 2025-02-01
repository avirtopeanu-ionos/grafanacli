// internal/telemetry/telemetry.go
package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LokiWriter struct {
	endpoint string
	client   *http.Client
}

// lokiPushRequest represents the JSON payload format expected by Loki.
type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func NewLokiWriter(endpoint string) *LokiWriter {
	return &LokiWriter{
		endpoint: endpoint,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

// Write sends the log line to Loki by building a JSON payload and POSTing it.
func (lw *LokiWriter) Write(p []byte) (n int, err error) {
	logLine := strings.TrimSpace(string(p))
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	payload := lokiPushRequest{
		Streams: []lokiStream{
			{
				Stream: map[string]string{
					"job": "mycli",
				},
				Values: [][]string{
					{ts, logLine},
				},
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return len(p), err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, lw.endpoint, bytes.NewBuffer(body))
	if err != nil {
		return len(p), err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := lw.client.Do(req)
	if err != nil {
		// If Loki is unavailable, we ignore the error.
		return len(p), nil
	}
	defer resp.Body.Close()
	// Optionally, check for non-2xx responses.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// For this demo, errors are simply ignored.
	}
	return len(p), nil
}

// InitLogger sets up zerolog to log to the console and, if available, to Loki.
func InitLogger() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	lokiEndpoint := os.Getenv("LOKI_ENDPOINT")
	if lokiEndpoint != "" {
		lokiWriter := NewLokiWriter(lokiEndpoint)
		// Use MultiWriter so logs go to both console and Loki.
		multi := io.MultiWriter(consoleWriter, lokiWriter)
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
		log.Info().Msg("Loki logging enabled")
	} else {
		log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
		log.Info().Msg("Logging to console only")
	}
}

// LogCommandExecution logs details about a command execution.
// It records the command name, arguments, flags (with their values), the output, and any error.
func LogCommandExecution(command string, args []string, flags map[string]string, output string, err error) {
	event := log.Info().
		Str("command", command).
		Strs("arguments", args).
		Interface("flags", flags).
		Str("output", output).
		Time("timestamp", time.Now())

	if err != nil {
		event = event.Err(err)
		event.Msg("Command execution failed")
	} else {
		event.Msg("Command executed successfully")
	}
}

// CatchPanic recovers from a panic and logs the error.
func CatchPanic() {
	if r := recover(); r != nil {
		log.Error().Interface("panic", r).Msg("Recovered from panic")
	}
}
