package logqlsyntax_test

import (
	"math/rand"
	"testing"
	"testing/quick"

	"github.com/qualithm/logql-syntax/syntax"
)

// Property: syntax.ParseExpr must never panic, regardless of input.
//
// This test lives outside syntax/ so the upstream-sync script does not
// overwrite it. It exercises the same parser the vendored fuzz harness does,
// but as a regular `go test` target consistent with our property-based
// testing preference.
func TestProperty_ParseExpr_NoPanic(t *testing.T) {
	t.Parallel()
	f := func(s string) bool {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic on input %q: %v", s, r)
			}
		}()
		_, _ = syntax.ParseExpr(s)
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 500}); err != nil {
		t.Fatal(err)
	}
}

// Property: parsing then re-parsing the printed form of a successfully parsed
// expression must succeed. This is a round-trip / idempotence check over a
// small grammar of plausible LogQL expressions.
func TestProperty_ParseExpr_RoundTrip(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 300; i++ {
		expr := randLogQL(rng)
		parsed, err := syntax.ParseExpr(expr)
		if err != nil {
			continue
		}
		printed := parsed.String()
		if _, err := syntax.ParseExpr(printed); err != nil {
			t.Fatalf("re-parse failed for %q -> %q: %v", expr, printed, err)
		}
	}
}

// Property: ParseExpr is deterministic — same input produces equal printed form.
func TestProperty_ParseExpr_Deterministic(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(7))
	for i := 0; i < 200; i++ {
		expr := randLogQL(rng)
		a, errA := syntax.ParseExpr(expr)
		b, errB := syntax.ParseExpr(expr)
		if (errA == nil) != (errB == nil) {
			t.Fatalf("nondeterministic error for %q: %v vs %v", expr, errA, errB)
		}
		if errA != nil {
			continue
		}
		if a.String() != b.String() {
			t.Fatalf("nondeterministic print for %q:\n %s\n vs\n %s", expr, a.String(), b.String())
		}
	}
}

func randLogQL(r *rand.Rand) string {
	labels := []string{"app", "env", "job"}
	values := []string{"api", "web", "prod", "staging"}
	lineOps := []string{"|=", "!=", "|~", "!~"}
	parsers := []string{"| json", "| logfmt"}

	pick := func(s []string) string { return s[r.Intn(len(s))] }

	expr := `{` + pick(labels) + `="` + pick(values) + `"}`
	if r.Intn(2) == 0 {
		expr += " " + pick(lineOps) + ` "` + pick(values) + `"`
	}
	if r.Intn(2) == 0 {
		expr += " " + pick(parsers)
	}
	if r.Intn(3) == 0 {
		expr = "rate(" + expr + " [5m])"
	}
	if r.Intn(4) == 0 {
		expr = "sum by (" + pick(labels) + ") (" + expr + ")"
	}
	return expr
}
