package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_runExitAnalyzer(t *testing.T) {
	t.Parallel()

	analysistest.Run(t, analysistest.TestData(), ExitAnalyzer, "./...")
}
