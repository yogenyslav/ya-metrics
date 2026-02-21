package analyzer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/simple/s1000"
	"honnef.co/go/tools/staticcheck"
)

func TestStaticChecks(t *testing.T) {
	t.Parallel()

	saChecks := make([]*analysis.Analyzer, 0)
	for _, a := range staticcheck.Analyzers {
		if strings.HasPrefix(a.Analyzer.Name, "SA") {
			saChecks = append(saChecks, a.Analyzer)
		}
	}

	tests := []struct {
		name       string
		wantChecks map[string]bool
		want       []*analysis.Analyzer
	}{
		{
			name:       "no checks specified",
			wantChecks: map[string]bool{},
			want:       saChecks,
		},
		{
			name: "some checks are specified, but disabled",
			wantChecks: map[string]bool{
				"S1000": false,
			},
			want: saChecks,
		},
		{
			name: "some checks are specified and enabled",
			wantChecks: map[string]bool{
				"S1000": true,
			},
			want: append(saChecks, s1000.Analyzer),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			checks := StaticChecks(tt.wantChecks)
			assert.ElementsMatch(t, checks, tt.want)
		})
	}
}
