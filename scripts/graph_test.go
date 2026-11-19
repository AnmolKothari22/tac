package main

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestGenerateGraphsWritesSixClientPercentileSVGs(t *testing.T) {
	root := t.TempDir()
	dataDir := filepath.Join(root, "data")
	graphDir := filepath.Join(root, "graphs")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatal(err)
	}

	for _, clients := range []int{200, 300} {
		content := "run,server,command,p50,p90,p99\n" +
			"10,redis,SET,1,2,3\n10,tac,SET,2,4,6\n" +
			"10,redis,GET,2,3,4\n10,tac,GET,4,6,8\n"
		path := filepath.Join(dataDir, "data_"+strconv.Itoa(clients)+".csv")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := generateGraphs(dataDir, graphDir); err != nil {
		t.Fatal(err)
	}
	for _, clients := range []int{200, 300} {
		for _, metric := range []string{"p50", "p90", "p99"} {
			path := filepath.Join(graphDir, "client"+strconv.Itoa(clients)+"_"+metric+".svg")
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(content), "TAC / Redis") {
				t.Errorf("%s is missing the TAC/Redis ratio title", path)
			}
		}
	}
}

func TestRunLessSortsRunAttemptsChronologically(t *testing.T) {
	runs := []string{"100-2", "99-1", "100-1"}
	sort.Slice(runs, func(i, j int) bool {
		return runLess(runs[i], runs[j])
	})
	if got, want := strings.Join(runs, ","), "99-1,100-1,100-2"; got != want {
		t.Fatalf("sorted runs = %s, want %s", got, want)
	}
}
