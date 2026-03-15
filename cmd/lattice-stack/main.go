// Package main provides the lattice-stack CLI tool for scaffolding Department Stacks.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: lattice-stack init <stack-name>\n")
			os.Exit(1)
		}
		if err := initStack(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`lattice-stack — Scaffold Department Stacks for Lattice Runtime

Usage:
  lattice-stack init <stack-name>    Create a new Department Stack project

Commands:
  init     Create a new stack project with boilerplate code
  help     Show this help message

Examples:
  lattice-stack init hr-stack
  lattice-stack init my-finance-stack`)
}
