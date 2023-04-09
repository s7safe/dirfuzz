package output

import (
	"fmt"
	"io"
	"os"

	"github.com/projectdiscovery/gologger"
)

// Output defines an output instance to write output to
type Output struct {
	writer io.Writer
}

// New creates a new instance of output
func New(filename string) (*Output, error) {
	var writer io.Writer
	if filename == "" {
		writer = os.Stdout
	} else {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			gologger.Error().Msgf("Could not create output file: %s\n", err)
			return nil, err
		}
		writer = f
	}

	return &Output{writer: writer}, nil
}

// Println prints a string with a new line character
func (o *Output) Println(str string) {
	fmt.Fprintln(o.writer, str)
}

// Printf prints a formatted string
func (o *Output) Printf(format string, a ...interface{}) {
	fmt.Fprintf(o.writer, format, a...)
}
