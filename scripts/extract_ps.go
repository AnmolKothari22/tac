package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	commandLinePattern = regexp.MustCompile(`^=+\s*([A-Za-z0-9_]+)\s*=+$`)
	latencyLinePattern = regexp.MustCompile(`^\s*([0-9]+(\.[0-9]+)?)%\s*<=\s*([0-9]+(\.[0-9]+)?)\s+milliseconds(\s+\(cumulative count\s+([0-9]+)\))?`)
)

type distributionPoint struct {
	percentile float64
	latency    float64
	count      int64
}

type benchmarkResult struct {
	command      string
	summary      map[string]float64
	distribution []distributionPoint
}

func extractResults(input io.Reader, output io.Writer, run, server string, clientCount int) error {
	if run == "" || (server != "redis" && server != "tac") || clientCount <= 0 {
		return fmt.Errorf("run, server, and a positive client count are required")
	}

	results, err := parseBenchmarkOutput(input)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(output)
	for _, result := range results {
		p50, ok := result.summary["p50"]
		if !ok {
			return fmt.Errorf("%s output is missing p50", result.command)
		}
		p99, ok := result.summary["p99"]
		if !ok {
			return fmt.Errorf("%s output is missing p99", result.command)
		}
		p90, ok := result.summary["p90"]
		if !ok {
			p90, ok = percentileFromDistribution(result.distribution, 90)
		}
		if !ok {
			return fmt.Errorf("%s output is missing p90 latency data", result.command)
		}

		record := []string{
			run,
			server,
			result.command,
			strconv.FormatFloat(p50, 'f', 3, 64),
			strconv.FormatFloat(p90, 'f', 3, 64),
			strconv.FormatFloat(p99, 'f', 3, 64),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func parseBenchmarkOutput(input io.Reader) ([]benchmarkResult, error) {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 4096), 1024*1024)

	results := make([]benchmarkResult, 0, 2)
	byCommand := make(map[string]int)
	var current *benchmarkResult
	var summaryHeaders []string
	readingSummary := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if match := commandLinePattern.FindStringSubmatch(line); match != nil {
			command := strings.ToUpper(match[1])
			if _, exists := byCommand[command]; exists {
				return nil, fmt.Errorf("duplicate benchmark section for %s", command)
			}
			results = append(results, benchmarkResult{
				command: command,
				summary: make(map[string]float64),
			})
			current = &results[len(results)-1]
			byCommand[command] = len(results) - 1
			summaryHeaders = nil
			readingSummary = false
			continue
		}
		if current == nil {
			continue
		}

		if match := latencyLinePattern.FindStringSubmatch(line); match != nil {
			percentile, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				return nil, fmt.Errorf("parse latency percentile %q: %w", match[1], err)
			}
			latency, err := strconv.ParseFloat(match[3], 64)
			if err != nil {
				return nil, fmt.Errorf("parse latency value %q: %w", match[3], err)
			}
			point := distributionPoint{percentile: percentile, latency: latency, count: -1}
			if match[6] != "" {
				point.count, err = strconv.ParseInt(match[6], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse cumulative count %q: %w", match[6], err)
				}
			}
			current.distribution = append(current.distribution, point)
		}

		if strings.Contains(strings.ToLower(line), "latency summary") {
			readingSummary = true
			summaryHeaders = nil
			continue
		}
		if !readingSummary || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if summaryHeaders == nil {
			summaryHeaders = fields
			continue
		}
		if len(fields) != len(summaryHeaders) {
			return nil, fmt.Errorf("%s summary has %d values for %d columns", current.command, len(fields), len(summaryHeaders))
		}
		for i, header := range summaryHeaders {
			value, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				return nil, fmt.Errorf("parse %s latency %q: %w", header, fields[i], err)
			}
			current.summary[strings.ToLower(header)] = value
		}
		readingSummary = false
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read benchmark output: %w", err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no benchmark summary sections found")
	}
	for _, result := range results {
		if len(result.summary) == 0 {
			return nil, fmt.Errorf("%s output is missing a latency summary", result.command)
		}
	}
	return results, nil
}

func percentileFromDistribution(points []distributionPoint, percentile float64) (float64, bool) {
	if len(points) == 0 {
		return 0, false
	}

	var total int64
	hasCounts := false
	for _, point := range points {
		if point.count >= 0 {
			hasCounts = true
			if point.count > total {
				total = point.count
			}
		}
	}
	if hasCounts && total > 0 {
		target := int64(math.Ceil(float64(total) * percentile / 100))
		for _, point := range points {
			if point.count >= target {
				return point.latency, true
			}
		}
	}

	for _, point := range points {
		if point.percentile >= percentile {
			return point.latency, true
		}
	}
	return 0, false
}
