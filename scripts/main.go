package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "extract":
		err = runExtract(os.Args[2:])
	case "graph":
		err = runGraph(os.Args[2:])
	default:
		printUsage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runExtract(args []string) error {
	flags := flag.NewFlagSet("extract", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	run := flags.String("run", "", "unique benchmark run identifier")
	server := flags.String("server", "", "benchmark target: redis or tac")
	clientCount := flags.Int("client", 0, "number of benchmark clients")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("extract does not accept positional arguments")
	}
	return extractResults(os.Stdin, os.Stdout, *run, *server, *clientCount)
}

func runGraph(args []string) error {
	flags := flag.NewFlagSet("graph", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	dataDir := flags.String("data-dir", "data_on_run", "directory containing benchmark CSV files")
	outputDir := flags.String("output-dir", "graphs", "directory where SVG graphs are written")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("graph does not accept positional arguments")
	}
	return generateGraphs(*dataDir, *outputDir)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: benchmark-tools {extract|graph} [options]")
}
