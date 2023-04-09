package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
)

// CSVOutput represents CSV output writer.
type CSVOutput struct {
	filePath string
	writer   *csv.Writer
	mutex    sync.Mutex
}

// NewCSVOutput creates a new CSV output writer.
func NewCSVOutput(filePath string) (*CSVOutput, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	writer := csv.NewWriter(file)
	writer.Write([]string{"Time", "Method", "URL", "Payload", "Status", "Content-Type", "Content-Length"})
	return &CSVOutput{
		filePath: filePath,
		writer:   writer,
	}, nil
}

// Write writes the results in CSV format.
func (c *CSVOutput) Write(result Result) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	data := []string{
		time.Now().Format("2006-01-02 15:04:05"),
		result.Method,
		result.URL,
		result.Payload,
		fmt.Sprintf("%d", result.StatusCode),
		strings.Join(result.Headers.Values("Content-Type"), ","),
		fmt.Sprintf("%d", result.ContentLength),
	}
	err := c.writer.Write(data)
	if err != nil {
		return err
	}
	c.writer.Flush()
	return nil
}

// Close closes the CSV output writer.
func (c *CSVOutput) Close() error {
	return nil
}

func init() {
	RegisterOutput("csv", func(filePath string) (Output, error) {
		return NewCSVOutput(filePath)
	})
}

// PrintSummary prints a summary of the results to stdout.
func (c *CSVOutput) PrintSummary(summary Summary) {
	fmt.Println()
	color.Info.Tips("Summary:")
	color.Info.Tips("  Total requests............: %d", summary.Total)
	color.Info.Tips("  Successful requests.......: %d", summary.Successful)
	color.Info.Tips("  Failed requests...........: %d", summary.Failed)
	color.Info.Tips("  Percentage of successful..: %.2f%%", summary.SuccessRate())
	color.Info.Tips("  Total time................: %v", summary.TotalTime().Truncate(time.Millisecond))
	color.Info.Tips("  Average time..............: %v", summary.AverageTime().Truncate(time.Millisecond))
	color.Info.Tips("  Fastest time..............: %v", summary.FastestTime().Truncate(time.Millisecond))
	color.Info.Tips("  Slowest time..............: %v", summary.SlowestTime().Truncate(time.Millisecond))
}
