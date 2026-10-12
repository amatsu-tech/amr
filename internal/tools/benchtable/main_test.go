package main

import (
	"bytes"
	"slices"
	"strings"
	"testing"
)

func TestVerdict(t *testing.T) {
	tests := []struct {
		name        string
		significant bool
		change      float64
		want        string
	}{
		{name: "not significant", significant: false, change: -40, want: "no change (noise)"},
		{name: "tiny speedup", significant: true, change: -2, want: "tiny, ignore"},
		{name: "tiny slowdown", significant: true, change: 4.9, want: "tiny, ignore"},
		{name: "faster", significant: true, change: -85, want: "FASTER"},
		{name: "slower", significant: true, change: 20, want: "SLOWER"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := verdict(tt.significant, tt.change); got != tt.want {
				t.Errorf("verdict(%v, %v) = %q, want %q", tt.significant, tt.change, got, tt.want)
			}
		})
	}
}

// benchOutput строит вывод `go test -bench` с одинаковым временем во всех прогонах.
func benchOutput(lines ...string) string {
	return "goos: darwin\ngoarch: arm64\npkg: example\n" + strings.Join(lines, "\n") + "\nPASS\n"
}

func TestLoad(t *testing.T) {
	input := benchOutput(
		"BenchmarkScan/repos=10-14   10000   100000 ns/op   19672 B/op   153 allocs/op",
		"BenchmarkScan/repos=10-14   10000   120000 ns/op   19672 B/op   153 allocs/op",
		"BenchmarkParse-14   500000   2000 ns/op",
	)

	got, err := load(strings.NewReader(input), "bench.txt")
	if err != nil {
		t.Fatalf("load() error = %v", err)
	}

	if want := []string{"Scan/repos=10", "Parse"}; !slices.Equal(got.names, want) {
		t.Errorf("names = %v, want %v", got.names, want)
	}

	// ns/op приводится к секундам.
	if want := []float64{100e-6, 120e-6}; !slices.Equal(got.samples["Scan/repos=10"], want) {
		t.Errorf("samples = %v, want %v", got.samples["Scan/repos=10"], want)
	}
}

func TestLoadEmpty(t *testing.T) {
	if _, err := load(strings.NewReader("PASS\n"), "empty.txt"); err == nil {
		t.Error("load() of output without benchmarks succeeded, want error")
	}
}

func TestComparison(t *testing.T) {
	var before, after []string

	for range 6 {
		before = append(before, "BenchmarkStatus-14   100   10000000 ns/op")
		after = append(after, "BenchmarkStatus-14   100   1500000 ns/op")
	}

	beforeResults, err := load(strings.NewReader(benchOutput(before...)), "before.txt")
	if err != nil {
		t.Fatal(err)
	}

	afterResults, err := load(strings.NewReader(benchOutput(after...)), "after.txt")
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	printComparison(&out, beforeResults, afterResults)

	for _, want := range []string{"10.0000", "1.5000", "-85.00%", "FASTER"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("comparison output does not contain %q:\n%s", want, out.String())
		}
	}
}
