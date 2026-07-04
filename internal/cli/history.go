package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"

	"timer/internal/store"
	"timer/internal/timeutil"
)

const historyPageSize = 20

var (
	historySearch string
	historyExport string
	historyPage   int
)

var historyCmd = &cobra.Command{
	Use:     "history",
	Aliases: []string{"h"},
	Short:   "Show past timer sessions, with optional search, pagination, and export",
	RunE: func(cmd *cobra.Command, args []string) error {
		var records []store.HistoryRecord
		var err error
		if historySearch != "" {
			records, err = store.SearchHistory(historySearch)
		} else {
			records, err = store.AllHistory()
		}
		if err != nil {
			return err
		}

		sort.Slice(records, func(i, j int) bool {
			return records[i].StartedAt.After(records[j].StartedAt)
		})

		switch historyExport {
		case "":
			return printHistoryTable(records, historyPage, historyPageSize)
		case "json":
			return exportHistoryJSON(records)
		case "csv":
			return exportHistoryCSV(records)
		default:
			return fmt.Errorf("unknown --export format %q (want json or csv)", historyExport)
		}
	},
}

func init() {
	historyCmd.Flags().StringVar(&historySearch, "search", "", "filter records by title substring")
	historyCmd.Flags().StringVar(&historyExport, "export", "", "export format: json or csv (writes to stdout, ignores pagination)")
	historyCmd.Flags().IntVar(&historyPage, "page", 1, "page of results to show (20 records per page)")
	rootCmd.AddCommand(historyCmd)
}

// paginate slices all into the requested page of pageSize records, newest
// first (the caller is expected to have already sorted all). Returns the
// page's records and the total number of pages. Returns an error if page
// is out of range.
func paginate(all []store.HistoryRecord, page, pageSize int) ([]store.HistoryRecord, int, error) {
	if len(all) == 0 {
		return nil, 0, nil
	}

	totalPages := (len(all) + pageSize - 1) / pageSize
	if page < 1 || page > totalPages {
		return nil, totalPages, fmt.Errorf("page %d does not exist (%d page(s) available)", page, totalPages)
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], totalPages, nil
}

func printHistoryTable(all []store.HistoryRecord, page, pageSize int) error {
	th := ResolveTheme()
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(th.AccentColor))
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(th.HintColor))
	cellStyle := lipgloss.NewStyle().Padding(0, 1)

	if len(all) == 0 {
		fmt.Println("No timer history yet.")
		return nil
	}

	records, totalPages, err := paginate(all, page, pageSize)
	if err != nil {
		return err
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(th.HintColor))).
		Headers("TITLE", "MODE", "ELAPSED", "STARTED", "DONE").
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return cellStyle
		})

	for _, r := range records {
		title := r.Title
		if title == "" {
			title = "(untitled)"
		}
		done := "yes"
		if !r.Completed {
			done = "no"
		}
		t.Row(title, r.Mode, timeutil.FormatHMS(r.Elapsed), r.StartedAt.Local().Format("2006-01-02 15:04"), done)
	}
	fmt.Println(t.Render())

	start := (page-1)*pageSize + 1
	end := start + len(records) - 1
	fmt.Println()
	fmt.Println(hintStyle.Render(fmt.Sprintf("Page %d of %d (records %d–%d of %d)", page, totalPages, start, end, len(all))))
	if page < totalPages {
		fmt.Println(hintStyle.Render(fmt.Sprintf("Use --page %d to see the next page.", page+1)))
	}

	fmt.Println()
	fmt.Println(headerStyle.Render("Totals by title:"))
	for title, stats := range aggregateStats(all) {
		if title == "" {
			title = "(untitled)"
		}
		fmt.Printf("  %s: %d session(s), %s total\n", title, stats.Count, timeutil.FormatHMS(stats.Total))
	}
	return nil
}

func aggregateStats(records []store.HistoryRecord) map[string]store.TitleStats {
	stats := make(map[string]store.TitleStats)
	for _, r := range records {
		s := stats[r.Title]
		s.Count++
		s.Total += r.Elapsed
		stats[r.Title] = s
	}
	return stats
}

func exportHistoryJSON(records []store.HistoryRecord) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func exportHistoryCSV(records []store.HistoryRecord) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	header := []string{"id", "title", "mode", "target_seconds", "elapsed_seconds", "started_at", "ended_at", "completed"}
	if err := w.Write(header); err != nil {
		return err
	}

	for _, r := range records {
		row := []string{
			r.ID,
			r.Title,
			r.Mode,
			strconv.FormatFloat(r.Target.Seconds(), 'f', -1, 64),
			strconv.FormatFloat(r.Elapsed.Seconds(), 'f', -1, 64),
			r.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
			r.EndedAt.Format("2006-01-02T15:04:05Z07:00"),
			strconv.FormatBool(r.Completed),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}
