package main

import (
	"slices"

	"github.com/yogenyslav/ya-metrics/pkg/analyzer"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	wantChecks := map[string]bool{
		"ST1003": true,
		"S1021":  true,
		"QF1003": true,
	}

	multichecker.Main(
		slices.Concat(
			analyzer.StaticChecks(wantChecks),
			analyzer.DefaultAnalyzers,
			analyzer.ThirdPartyAnalyzers,
			analyzer.CustomAnalyzers,
		)...,
	)
}
