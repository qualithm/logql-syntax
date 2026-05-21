package logqlmodel

import (
	"errors"
	"strings"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
)

func TestParseError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  ParseError
		want string
	}{
		{
			name: "no location",
			err:  NewParseError("unexpected token", 0, 0),
			want: "parse error : unexpected token",
		},
		{
			name: "with location",
			err:  NewParseError("unexpected token", 3, 7),
			want: "parse error at line 3, col 7: unexpected token",
		},
		{
			name: "line only",
			err:  NewParseError("bad", 2, 0),
			want: "parse error at line 2, col 0: bad",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseError_IsErrParse(t *testing.T) {
	err := NewParseError("boom", 1, 2)
	if !errors.Is(err, ErrParse) {
		t.Fatalf("ParseError should match ErrParse via errors.Is")
	}
	if errors.Is(err, ErrPipeline) {
		t.Fatalf("ParseError should not match ErrPipeline")
	}
}

func TestNewStageError(t *testing.T) {
	inner := errors.New("oops")
	err := NewStageError("| json", inner)
	if !errors.Is(err, ErrParse) {
		t.Fatalf("StageError should match ErrParse")
	}
	got := err.Error()
	if !strings.Contains(got, "| json") || !strings.Contains(got, "oops") {
		t.Fatalf("stage error message missing pieces: %q", got)
	}
}

func TestPipelineError(t *testing.T) {
	lbls := labels.FromStrings(ErrorLabel, "JSONParserErr", "app", "foo")
	pe := NewPipelineErr(lbls)
	if pe.errorType != "JSONParserErr" {
		t.Fatalf("errorType = %q, want JSONParserErr", pe.errorType)
	}
	msg := pe.Error()
	for _, want := range []string{"pipeline error", "JSONParserErr", "app", "foo"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("pipeline error message missing %q: %s", want, msg)
		}
	}
	if !errors.Is(*pe, ErrPipeline) {
		t.Fatalf("PipelineError should match ErrPipeline")
	}
	if errors.Is(*pe, ErrParse) {
		t.Fatalf("PipelineError should not match ErrParse")
	}
}

func TestLimitError(t *testing.T) {
	le := NewSeriesLimitError(1000)
	if le == nil {
		t.Fatal("NewSeriesLimitError returned nil")
	}
	if !errors.Is(*le, ErrLimit) {
		t.Fatalf("LimitError should match ErrLimit")
	}
	if errors.Is(*le, ErrParse) {
		t.Fatalf("LimitError should not match ErrParse")
	}
	if !strings.Contains(le.Error(), "1000") {
		t.Fatalf("limit error should mention limit value: %s", le.Error())
	}
}

func TestConstants(t *testing.T) {
	if ValueTypeStreams != "streams" {
		t.Errorf("ValueTypeStreams = %q", ValueTypeStreams)
	}
	if PackedEntryKey != "_entry" {
		t.Errorf("PackedEntryKey = %q", PackedEntryKey)
	}
	if ErrorLabel != "__error__" || PreserveErrorLabel != "__preserve_error__" || ErrorDetailsLabel != "__error_details__" {
		t.Errorf("error label constants drifted: %q %q %q", ErrorLabel, PreserveErrorLabel, ErrorDetailsLabel)
	}
}
