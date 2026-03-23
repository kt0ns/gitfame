package format

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	"gitlab.com/slon/shad-go/gitfame/internal/errors"
	"gitlab.com/slon/shad-go/gitfame/internal/statistics"
)

type StatRecord struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

type printerFunc func(io.Writer, []StatRecord) error

var formatPrinters = map[string]printerFunc{
	"":           printTabular,
	"tabular":    printTabular,
	"csv":        printCSV,
	"json":       printJSON,
	"json-lines": printJSONLines,
}

func PrintStats(out io.Writer, format string, stats []statistics.AuthorStats) error {
	records := make([]StatRecord, len(stats))
	for i, s := range stats {
		records[i] = StatRecord{
			Name:    s.Name,
			Lines:   s.Lines,
			Commits: len(s.Commits),
			Files:   len(s.Files),
		}
	}

	printer, ok := formatPrinters[format]
	if !ok {
		return fmt.Errorf(errors.MsgUnknownFormat, format)
	}

	return printer(out, records)
}

func printTabular(out io.Writer, records []StatRecord) error {
	w := tabwriter.NewWriter(out, 0, 0, 1, ' ', 0)

	if _, err := fmt.Fprint(w, "Name\tLines\tCommits\tFiles\r\n"); err != nil {
		return err
	}

	for _, rec := range records {
		if _, err := fmt.Fprintf(w, "%s\t%d\t%d\t%d\r\n", rec.Name, rec.Lines, rec.Commits, rec.Files); err != nil {
			return err
		}
	}

	return w.Flush()
}

func printCSV(out io.Writer, records []StatRecord) error {
	w := csv.NewWriter(out)
	w.UseCRLF = true

	if err := w.Write([]string{"Name", "Lines", "Commits", "Files"}); err != nil {
		return err
	}

	for _, rec := range records {
		row := []string{
			rec.Name,
			strconv.Itoa(rec.Lines),
			strconv.Itoa(rec.Commits),
			strconv.Itoa(rec.Files),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

func printJSON(out io.Writer, records []StatRecord) error {
	if records == nil {
		records = make([]StatRecord, 0)
	}

	enc := json.NewEncoder(out)
	return enc.Encode(records)
}

func printJSONLines(out io.Writer, records []StatRecord) error {
	enc := json.NewEncoder(out)
	for _, rec := range records {
		if err := enc.Encode(rec); err != nil {
			return err
		}
	}
	return nil
}
