package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"
	"testing"
)

func TestParseBenchmarkOutputExtractsP90FromHistogram(t *testing.T) {
	output := `====== SET ======
Latency by percentile distribution:
50.000% <= 0.100 milliseconds (cumulative count 50)
87.500% <= 0.300 milliseconds (cumulative count 88)
93.750% <= 0.400 milliseconds (cumulative count 94)
100.000% <= 1.000 milliseconds (cumulative count 100)

Summary:
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.200     0.050     0.100     0.400     0.800     1.000
====== GET ======
Latency by percentile distribution:
50.000% <= 0.050 milliseconds (cumulative count 50)
90.000% <= 0.200 milliseconds (cumulative count 90)
100.000% <= 0.500 milliseconds (cumulative count 100)

Summary:
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.100     0.020     0.050     0.300     0.400     0.500
`

	output = "SET: rps=70000.0 (overall: 70000.0) avg_msec=1.000\r" +
		"GET: rps=70000.0 (overall: 70000.0) avg_msec=1.000\r" +
		strings.ReplaceAll(output, "\n", "\r")
	results, err := parseBenchmarkOutput(strings.NewReader(output))
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].command != "SET" {
		t.Fatalf("first command = %q, want SET", results[0].command)
	}
	if p90, ok := percentileFromDistribution(results[0].distribution, 90); !ok || p90 != 0.4 {
		t.Fatalf("SET p90 = %v, %v; want 0.4, true", p90, ok)
	}
	if p90, ok := percentileFromDistribution(results[1].distribution, 90); !ok || p90 != 0.2 {
		t.Fatalf("GET p90 = %v, %v; want 0.2, true", p90, ok)
	}
}

func TestPercentileFromDistributionWithoutCounts(t *testing.T) {
	points := []distributionPoint{
		{percentile: 87.5, latency: 0.3, count: -1},
		{percentile: 93.75, latency: 0.4, count: -1},
	}

	if p90, ok := percentileFromDistribution(points, 90); !ok || p90 != 0.4 {
		t.Fatalf("p90 = %v, %v; want 0.4, true", p90, ok)
	}
}

func TestExtractResultsWritesBenchmarkCSVRecords(t *testing.T) {
	input := `====== SET ======
Latency by percentile distribution:
90.000% <= 0.200 milliseconds (cumulative count 90)
100.000% <= 1.000 milliseconds (cumulative count 100)
Summary:
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.200     0.050     0.100     0.400     0.800     1.000
`
	var output bytes.Buffer
	if err := extractResults(strings.NewReader(input), &output, "run-1", "tac", 200); err != nil {
		t.Fatal(err)
	}

	reader := csv.NewReader(&output)
	record, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"run-1", "tac", "SET", "0.100", "0.200", "0.800"}; strings.Join(record, ",") != strings.Join(want, ",") {
		t.Fatalf("record = %v, want %v", record, want)
	}
	if _, err := reader.Read(); err != io.EOF {
		t.Fatalf("second CSV read error = %v, want EOF", err)
	}
}
