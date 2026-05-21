// Package logqlsyntax is the root of a partial fork of Grafana Loki's LogQL
// parser packages, stripped of Loki's runtime dependencies (dskit, etcd,
// jaeger, queryrange, push, etc.).
//
// Use sub-packages to parse and walk LogQL expressions:
//
//	import "github.com/qualithm/logql-syntax/syntax"
//
//	expr, err := syntax.ParseExpr(`sum by (job) (rate({app="api"} |= "error" [5m]))`)
//	if err != nil {
//	    return err
//	}
//	expr.Walk(func(e syntax.Expr) bool { return true })
//
// See the README for the relationship to upstream Loki and the sync workflow.
package logqlsyntax
