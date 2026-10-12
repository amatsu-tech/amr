// Command benchtable показывает результаты `go test -bench` таблицей в миллисекундах.
package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"text/tabwriter"

	"golang.org/x/perf/benchfmt"
	"golang.org/x/perf/benchmath"
)

// Изменение меньше этого порога считаем шумом, даже если тест назвал его значимым.
const minRelevantChange = 5.0

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "benchtable:", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	switch len(args) {
	case 1:
		results, err := loadFile(args[0])
		if err != nil {
			return err
		}

		printSummary(out, results)

		return nil
	case 2:
		before, err := loadFile(args[0])
		if err != nil {
			return err
		}

		after, err := loadFile(args[1])
		if err != nil {
			return err
		}

		printComparison(out, before, after)

		return nil
	default:
		return fmt.Errorf("usage: benchtable <bench.txt> | benchtable <before.txt> <after.txt>")
	}
}

// results — времена одного вызова в секундах по каждому бенчмарку, в порядке появления.
type results struct {
	names   []string
	samples map[string][]float64
}

// gomaxprocsSuffix — хвост вида "-14" в имени бенчмарка (число ядер).
var gomaxprocsSuffix = regexp.MustCompile(`-\d+$`)

func loadFile(path string) (results, error) {
	file, err := os.Open(path)
	if err != nil {
		return results{}, err
	}
	defer file.Close()

	return load(file, path)
}

func load(r io.Reader, name string) (results, error) {
	parsed := results{samples: map[string][]float64{}}
	reader := benchfmt.NewReader(r, name)

	for reader.Scan() {
		result, ok := reader.Result().(*benchfmt.Result)
		if !ok {
			continue
		}

		seconds, ok := result.Value("sec/op")
		if !ok {
			continue
		}

		benchmark := gomaxprocsSuffix.ReplaceAllString(string(result.Name.Full()), "")

		if _, seen := parsed.samples[benchmark]; !seen {
			parsed.names = append(parsed.names, benchmark)
		}

		parsed.samples[benchmark] = append(parsed.samples[benchmark], seconds)
	}

	if err := reader.Err(); err != nil {
		return results{}, err
	}

	if len(parsed.names) == 0 {
		return results{}, fmt.Errorf("%s: no benchmark results", name)
	}

	return parsed, nil
}

// median возвращает медиану выборки — так же, как её показывает benchstat.
func median(values []float64) float64 {
	sample := benchmath.NewSample(values, &benchmath.DefaultThresholds)

	return benchmath.AssumeNothing.Summary(sample, 0.95).Center
}

func milliseconds(seconds float64) string {
	return fmt.Sprintf("%.4f", seconds*1000)
}

func printSummary(out io.Writer, r results) {
	table := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	fmt.Fprintln(table, "Benchmark\tms per call")

	for _, name := range r.names {
		fmt.Fprintf(table, "%s\t%s\n", name, milliseconds(median(r.samples[name])))
	}

	table.Flush()
}

func printComparison(out io.Writer, before, after results) {
	table := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	fmt.Fprintln(table, "Benchmark\tbefore ms\tafter ms\tchange\tverdict")

	for _, name := range after.names {
		beforeValues, ok := before.samples[name]
		if !ok {
			fmt.Fprintf(table, "%s\t-\t%s\t\tnew benchmark\n", name, milliseconds(median(after.samples[name])))
			continue
		}

		afterValues := after.samples[name]

		comparison := benchmath.AssumeNothing.Compare(
			benchmath.NewSample(beforeValues, &benchmath.DefaultThresholds),
			benchmath.NewSample(afterValues, &benchmath.DefaultThresholds),
		)

		beforeMedian, afterMedian := median(beforeValues), median(afterValues)
		change := (afterMedian - beforeMedian) / beforeMedian * 100
		significant := comparison.P < comparison.Alpha

		fmt.Fprintf(table, "%s\t%s\t%s\t%+.2f%%\t%s\n",
			name, milliseconds(beforeMedian), milliseconds(afterMedian), change, verdict(significant, change))
	}

	table.Flush()
}

// verdict переводит результат сравнения в короткий вывод для человека.
func verdict(significant bool, change float64) string {
	switch {
	case !significant:
		return "no change (noise)"
	case change > -minRelevantChange && change < minRelevantChange:
		return "tiny, ignore"
	case change < 0:
		return "FASTER"
	default:
		return "SLOWER"
	}
}
