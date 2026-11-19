package main

import (
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	graphWidth   = 1000
	graphHeight  = 600
	marginLeft   = 90
	marginRight  = 40
	marginTop    = 70
	marginBottom = 80
)

type sampleKey struct {
	run     string
	server  string
	command string
}

type latencySample struct {
	p50 float64
	p90 float64
	p99 float64
}

func generateGraphs(dataDir, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create graph directory: %w", err)
	}
	for _, clients := range []int{200, 300} {
		path := filepath.Join(dataDir, fmt.Sprintf("data_%d.csv", clients))
		samples, runs, err := loadSamples(path)
		if err != nil {
			return err
		}
		for _, metric := range []string{"p50", "p90", "p99"} {
			svg, err := renderGraph(clients, metric, samples, runs)
			if err != nil {
				return err
			}
			name := fmt.Sprintf("client%d_%s.svg", clients, metric)
			if err := os.WriteFile(filepath.Join(outputDir, name), []byte(svg), 0644); err != nil {
				return fmt.Errorf("write %s: %w", name, err)
			}
		}
	}
	return nil
}

func loadSamples(path string) (map[sampleKey]latencySample, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("read %s header: %w", path, err)
	}
	expectedHeader := []string{"run", "server", "command", "p50", "p90", "p99"}
	if strings.Join(header, ",") != strings.Join(expectedHeader, ",") {
		return nil, nil, fmt.Errorf("%s has unexpected header %q", path, header)
	}

	samples := make(map[sampleKey]latencySample)
	runSet := make(map[string]struct{})
	for line := 2; ; line++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("read %s row %d: %w", path, line, err)
		}
		if len(record) != len(expectedHeader) {
			return nil, nil, fmt.Errorf("%s row %d has %d columns, want %d", path, line, len(record), len(expectedHeader))
		}
		if record[0] == "" || (record[1] != "redis" && record[1] != "tac") ||
			(record[2] != "SET" && record[2] != "GET") {
			return nil, nil, fmt.Errorf("%s row %d has invalid run, server, or command", path, line)
		}

		values := [3]float64{}
		for i, field := range record[3:] {
			values[i], err = strconv.ParseFloat(field, 64)
			if err != nil || values[i] <= 0 {
				return nil, nil, fmt.Errorf("%s row %d has invalid %s value %q", path, line, expectedHeader[i+3], field)
			}
		}
		key := sampleKey{run: record[0], server: record[1], command: record[2]}
		if _, exists := samples[key]; exists {
			return nil, nil, fmt.Errorf("%s row %d duplicates run %s, %s, %s", path, line, key.run, key.server, key.command)
		}
		samples[key] = latencySample{p50: values[0], p90: values[1], p99: values[2]}
		runSet[key.run] = struct{}{}
	}

	runs := make([]string, 0, len(runSet))
	for run := range runSet {
		runs = append(runs, run)
	}
	sort.Slice(runs, func(i, j int) bool {
		return runLess(runs[i], runs[j])
	})

	for _, run := range runs {
		for _, command := range []string{"SET", "GET"} {
			for _, server := range []string{"redis", "tac"} {
				key := sampleKey{run: run, server: server, command: command}
				if _, exists := samples[key]; !exists {
					return nil, nil, fmt.Errorf("%s is missing run %s, %s, %s", path, run, server, command)
				}
			}
		}
	}
	return samples, runs, nil
}

func runLess(left, right string) bool {
	parseRun := func(run string) (uint64, uint64, bool) {
		parts := strings.SplitN(run, "-", 2)
		number, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return 0, 0, false
		}
		var attempt uint64
		if len(parts) == 2 {
			attempt, err = strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return 0, 0, false
			}
		}
		return number, attempt, true
	}

	leftNumber, leftAttempt, leftOK := parseRun(left)
	rightNumber, rightAttempt, rightOK := parseRun(right)
	if leftOK && rightOK {
		if leftNumber != rightNumber {
			return leftNumber < rightNumber
		}
		return leftAttempt < rightAttempt
	}
	return left < right
}

func renderGraph(clients int, metric string, samples map[sampleKey]latencySample, runs []string) (string, error) {
	if metric != "p50" && metric != "p90" && metric != "p99" {
		return "", fmt.Errorf("unsupported percentile %q", metric)
	}
	if len(runs) == 0 {
		return "", fmt.Errorf("cannot generate client%d_%s.svg: no benchmark runs", clients, metric)
	}

	value := func(sample latencySample) float64 {
		switch metric {
		case "p50":
			return sample.p50
		case "p90":
			return sample.p90
		default:
			return sample.p99
		}
	}

	maxRatio := 1.0
	for _, run := range runs {
		for _, command := range []string{"SET", "GET"} {
			redis := value(samples[sampleKey{run: run, server: "redis", command: command}])
			tac := value(samples[sampleKey{run: run, server: "tac", command: command}])
			if redis <= 0 {
				return "", fmt.Errorf("Redis %s latency for run %s must be positive", metric, run)
			}
			if ratio := tac / redis; ratio > maxRatio {
				maxRatio = ratio
			}
		}
	}
	yMax := maxRatio * 1.1
	if yMax < 1 {
		yMax = 1
	}

	plotWidth := graphWidth - marginLeft - marginRight
	plotHeight := graphHeight - marginTop - marginBottom
	xForRun := func(index int) float64 {
		if len(runs) == 1 {
			return float64(marginLeft + plotWidth/2)
		}
		return float64(marginLeft) + float64(index)*float64(plotWidth)/float64(len(runs)-1)
	}
	yForRatio := func(ratio float64) float64 {
		return float64(marginTop+plotHeight) - ratio/yMax*float64(plotHeight)
	}

	var svg strings.Builder
	fmt.Fprintf(&svg, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, graphWidth, graphHeight, graphWidth, graphHeight)
	svg.WriteString(`<rect width="100%" height="100%" fill="white"/>`)
	fmt.Fprintf(&svg, `<text x="%d" y="35" text-anchor="middle" font-family="sans-serif" font-size="22">%d clients - %s latency ratio (TAC / Redis)</text>`, graphWidth/2, clients, html.EscapeString(metric))

	for tick := 0; tick <= 5; tick++ {
		ratio := yMax * float64(tick) / 5
		y := yForRatio(ratio)
		fmt.Fprintf(&svg, `<line x1="%d" y1="%.2f" x2="%d" y2="%.2f" stroke="#dddddd"/>`, marginLeft, y, graphWidth-marginRight, y)
		fmt.Fprintf(&svg, `<text x="%d" y="%.2f" text-anchor="end" dominant-baseline="middle" font-family="sans-serif" font-size="13">%.2fx</text>`, marginLeft-12, y, ratio)
	}
	if yMax >= 1 {
		y := yForRatio(1)
		fmt.Fprintf(&svg, `<line x1="%d" y1="%.2f" x2="%d" y2="%.2f" stroke="#777777" stroke-dasharray="6 4"/>`, marginLeft, y, graphWidth-marginRight, y)
	}

	fmt.Fprintf(&svg, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#222222"/>`, marginLeft, marginTop, marginLeft, graphHeight-marginBottom)
	fmt.Fprintf(&svg, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#222222"/>`, marginLeft, graphHeight-marginBottom, graphWidth-marginRight, graphHeight-marginBottom)
	for i := range runs {
		x := xForRun(i)
		fmt.Fprintf(&svg, `<line x1="%.2f" y1="%d" x2="%.2f" y2="%d" stroke="#222222"/>`, x, graphHeight-marginBottom, x, graphHeight-marginBottom+5)
		fmt.Fprintf(&svg, `<text x="%.2f" y="%d" text-anchor="middle" font-family="sans-serif" font-size="12">%d</text>`, x, graphHeight-marginBottom+24, i+1)
	}
	fmt.Fprintf(&svg, `<text x="%d" y="%d" text-anchor="middle" font-family="sans-serif" font-size="15">Successful run number (lifetime)</text>`, marginLeft+plotWidth/2, graphHeight-20)
	fmt.Fprintf(&svg, `<text x="22" y="%d" transform="rotate(-90 22 %d)" text-anchor="middle" font-family="sans-serif" font-size="15">TAC / Redis latency ratio</text>`, marginTop+plotHeight/2, marginTop+plotHeight/2)

	colors := map[string]string{"SET": "#1769aa", "GET": "#d1495b"}
	for _, command := range []string{"SET", "GET"} {
		points := make([]string, 0, len(runs))
		for i, run := range runs {
			redis := value(samples[sampleKey{run: run, server: "redis", command: command}])
			tac := value(samples[sampleKey{run: run, server: "tac", command: command}])
			points = append(points, fmt.Sprintf("%.2f,%.2f", xForRun(i), yForRatio(tac/redis)))
		}
		fmt.Fprintf(&svg, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2.5"/>`, strings.Join(points, " "), colors[command])
		for _, point := range strings.Split(strings.Join(points, " "), " ") {
			coordinates := strings.Split(point, ",")
			fmt.Fprintf(&svg, `<circle cx="%s" cy="%s" r="4" fill="%s"/>`, coordinates[0], coordinates[1], colors[command])
		}
	}
	fmt.Fprintf(&svg, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="3"/><text x="%d" y="%d" font-family="sans-serif" font-size="14">SET</text>`,
		graphWidth-190, marginTop-14, graphWidth-160, marginTop-14, colors["SET"], graphWidth-152, marginTop-9)
	fmt.Fprintf(&svg, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="3"/><text x="%d" y="%d" font-family="sans-serif" font-size="14">GET</text>`,
		graphWidth-100, marginTop-14, graphWidth-70, marginTop-14, colors["GET"], graphWidth-62, marginTop-9)
	svg.WriteString(`</svg>`)
	return svg.String(), nil
}
