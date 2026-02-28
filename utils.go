// Package main — utility functions with intentional lint issues for CI demo.
package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// processItems has multiple lint issues: exported function without doc comment format,
// error handling issues, and shadowed variables.
func processItems(items []string) []string {
	var results []string // nolint - intentional var usage
	for i := 0; i < len(items); i++ {
		item := items[i]
		if item != "" {
			item := strings.ToUpper(item) // shadowed variable (lint error)
			results = append(results, item)
			_ = item
		}
	}
	return results
}

// RunCommand executes a shell command — security issue (command injection).
func RunCommand(userInput string) (string, error) {
	// gosec G204: subprocess launched with variable — security finding
	cmd := exec.Command("sh", "-c", userInput)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// formatData has inefficient string concatenation (lint warning).
func formatData(items []string) string {
	result := "" // ineffassign / inefficient concatenation
	for _, item := range items {
		result = result + item + ", " // should use strings.Builder
	}
	return result
}

// DeadCode is an exported function that's never used.
func DeadCode() {
	x := 1
	y := 2
	z := x + y
	fmt.Println(z)
}

// calculateTotal has an unreachable return.
func calculateTotal(prices []float64, taxRate float64) float64 {
	if len(prices) == 0 {
		return 0
	}

	total := 0.0
	for _, p := range prices {
		total += p
	}

	total = total * (1 + taxRate)
	return total
	fmt.Println("done") // unreachable code (lint error)
}
