package logqlsyntax_test

import (
	"testing"

	"github.com/qualithm/logql-syntax/syntax"
)

// Seeds span the grammar: selectors, line and label filters, parsers,
// formatters, range and vector aggregations, and binary operations.
var logqlSeeds = []string{
	`{app="api"}`,
	`{app="api", env!="dev", team=~"a|b", zone!~"x.*"}`,
	`{app="api"} |= "error" != "timeout" |~ "5\\d\\d" !~ "health"`,
	`{app="api"} | json | status >= 500 and duration > 1s`,
	`{app="api"} | logfmt | level="error" or level="warn"`,
	`{app="api"} | regexp "(?P<method>\\w+) (?P<path>\\S+)"`,
	`{app="api"} | pattern "<ip> - - <_> \"<method> <uri> <_>\" <status>"`,
	`{app="api"} | line_format "{{.status}} {{.path}}" | label_format dst=src`,
	`{app="api"} | unpack | drop foo, bar | keep baz`,
	`{app="api"} | decolorize`,
	`rate({app="api"}[5m])`,
	`sum by (status) (count_over_time({app="api"} | json [1m]))`,
	`quantile_over_time(0.99, {app="api"} | json | unwrap duration(latency) [5m]) by (path)`,
	`topk(5, sum without (pod) (rate({app="api"}[1m] offset 1h)))`,
	`sum(rate({a="b"}[1m])) / sum(rate({c="d"}[1m])) > bool 0.5`,
	`label_replace(rate({a="b"}[1m]), "dst", "$1", "src", "(.*)")`,
	`vector(1) + on (x) group_left (y) vector(2)`,
	`absent_over_time({app="api"}[10m])`,
}

// FuzzParseExpr checks that the LogQL parser never panics, and that a parsed
// expression prints and reparses without panicking.
func FuzzParseExpr(f *testing.F) {
	for _, seed := range logqlSeeds {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, query string) {
		expr, err := syntax.ParseExpr(query)
		if err != nil {
			return
		}
		_, _ = syntax.ParseExpr(expr.String())
	})
}
