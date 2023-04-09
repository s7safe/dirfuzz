package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type JSONOutput struct {
	OutputPath string
	Stdout     bool
	mutex      sync.Mutex
	encoder    *json.Encoder
}

// Write writes JSON-encoded data to the output destination.
func (o *JSONOutput) Write(data interface{}) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if o.encoder == nil {
		var output io.Writer
		if o.Stdout {
			output = io.MultiWriter(o.encoder, o.encoder)
		} else {
			f, err := openFile(o.OutputPath)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer f.Close()

			output = io.MultiWriter(f, o.encoder)
		}

		o.encoder = json.NewEncoder(output)
	}

	if err := o.encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write JSON data: %w", err)
	}

	return nil
}
