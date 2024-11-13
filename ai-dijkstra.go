package main

import (
	"ai-dijkstra/cmd"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

func main() {
	usage := `AI - dijkstra
Usage:
	ai-dijkstra prepare_data
	ai-dijkstra test_prepare_data
	ai-dijkstra train [--epochs <epochs>] [--cpu]
	ai-dijkstra run [--epochs <epochs>] [--cpu]
	ai-dijkstra evaluate_file <path> [--cpu]
	ai-dijkstra evaluate_dir [--cache] <path> [--cpu]
	ai-dijkstra pull_codes
Options:
	--epochs <epochs>
	--cache
`
	color.Output = ansi.NewAnsiStdout()

	opts, _ := docopt.ParseArgs(usage, os.Args[1:], "AI - dijkstra")

	err := cmd.Eval(opts)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	color.Unset()
}
