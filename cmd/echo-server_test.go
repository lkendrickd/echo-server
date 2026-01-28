package main

import (
	"log/slog"
	"testing"
)

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		wantLevel slog.Level
	}{
		{
			name:      "debug level",
			level:     "debug",
			wantLevel: slog.LevelDebug,
		},
		{
			name:      "info level",
			level:     "info",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "warn level",
			level:     "warn",
			wantLevel: slog.LevelWarn,
		},
		{
			name:      "error level",
			level:     "error",
			wantLevel: slog.LevelError,
		},
		{
			name:      "unknown level defaults to info",
			level:     "unknown",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "empty string defaults to info",
			level:     "",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "uppercase DEBUG defaults to info",
			level:     "DEBUG",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "mixed case Info defaults to info",
			level:     "Info",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "invalid level defaults to info",
			level:     "fatal",
			wantLevel: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setLogLevel(tt.level)

			if got != tt.wantLevel {
				t.Errorf("setLogLevel(%q) = %v, want %v", tt.level, got, tt.wantLevel)
			}
		})
	}
}
