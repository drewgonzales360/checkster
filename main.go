package main

import (
	"github.com/drewgonzales360/checkster/internal/analyzers/funclen"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(funclen.Analyzer)
}
