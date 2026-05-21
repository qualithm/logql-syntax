// Package logqlmodel contains the error types and label-name constants used by
// the LogQL syntax and evaluation packages.
//
// This is a trimmed fork of github.com/grafana/loki/v3/pkg/logqlmodel that
// omits the runtime Result/Streams types (which depend on
// github.com/grafana/loki/pkg/push and the queryrange-base definitions). Only
// the symbols required by the syntax and log packages are retained.
package logqlmodel

import (
	"errors"
	"fmt"

	"github.com/prometheus/prometheus/model/labels"
)

// Errors useful for comparison via errors.Is.
//
//	errors.Is(err, logqlmodel.ErrParse) // is this an AST parse error?
var (
	ErrParse                            = errors.New("failed to parse the log query")
	ErrPipeline                         = errors.New("failed execute pipeline")
	ErrLimit                            = errors.New("limit reached while evaluating the query")
	ErrIntervalLimit                    = errors.New("[interval] value exceeds limit")
	ErrBlocked                          = errors.New("query blocked by policy")
	ErrParseMatchers                    = errors.New("only label matchers are supported")
	ErrUnsupportedSyntaxForInstantQuery = errors.New(
		"log queries are not supported as an instant query type, please change your query to a range query type",
	)
	ErrVariantsDisabled = errors.New(
		"multi variant queries are disabled for this instance",
	)
)

// Special label names used by the pipeline error machinery.
const (
	ErrorLabel         = "__error__"
	PreserveErrorLabel = "__preserve_error__"
	ErrorDetailsLabel  = "__error_details__"
)

// ValueTypeStreams is the parser.ValueType for log streams.
const ValueTypeStreams = "streams"

// PackedEntryKey is a special JSON key used by the pack promtail stage and
// the unpack parser.
const PackedEntryKey = "_entry"

// ParseError is returned when the parser fails.
type ParseError struct {
	msg       string
	line, col int
}

// Error implements the error interface.
func (p ParseError) Error() string {
	if p.col == 0 && p.line == 0 {
		return fmt.Sprintf("parse error : %s", p.msg)
	}
	return fmt.Sprintf("parse error at line %d, col %d: %s", p.line, p.col, p.msg)
}

// Is allows errors.Is(err, ErrParse) on this error.
func (p ParseError) Is(target error) bool {
	return target == ErrParse
}

// NewParseError constructs a ParseError at the given source location.
func NewParseError(msg string, line, col int) ParseError {
	return ParseError{msg: msg, line: line, col: col}
}

// NewStageError constructs a ParseError for a pipeline stage failure.
func NewStageError(expr string, err error) ParseError {
	return ParseError{msg: fmt.Sprintf(`stage '%s' : %s`, expr, err)}
}

// PipelineError is returned when a pipeline stage fails at evaluation time.
type PipelineError struct {
	metric    labels.Labels
	errorType string
}

// NewPipelineErr constructs a PipelineError from the given series labels.
func NewPipelineErr(metric labels.Labels) *PipelineError {
	return &PipelineError{
		metric:    metric,
		errorType: metric.Get(ErrorLabel),
	}
}

// Error implements the error interface.
func (e PipelineError) Error() string {
	return fmt.Sprintf(
		"pipeline error: '%s' for series: '%s'.\n"+
			"Use a label filter to intentionally skip this error. (e.g | __error__!=\"%s\").\n"+
			"To skip all potential errors you can match empty errors.(e.g __error__=\"\")\n"+
			"The label filter can also be specified after unwrap. (e.g | unwrap latency | __error__=\"\" )\n",
		e.errorType, e.metric, e.errorType)
}

// Is allows errors.Is(err, ErrPipeline) on this error.
func (e PipelineError) Is(target error) bool {
	return target == ErrPipeline
}

// LimitError wraps an underlying error with the ErrLimit sentinel.
type LimitError struct {
	error
}

// NewSeriesLimitError constructs a LimitError for a series-cardinality limit.
func NewSeriesLimitError(limit int) *LimitError {
	return &LimitError{
		error: fmt.Errorf("maximum number of series (%d) reached for a single query; consider reducing query cardinality by adding more specific stream selectors, reducing the time range, or aggregating results with functions like sum(), count() or topk()", limit),
	}
}

// Is allows errors.Is(err, ErrLimit) on this error.
func (e LimitError) Is(target error) bool {
	return target == ErrLimit
}
