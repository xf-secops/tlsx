package output

import (
	"bufio"
	"os"
)

// fileWriter is a concurrent file based output writer.
type fileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

// NewFileOutputWriter creates a new buffered writer for a file
func NewFileOutputWriter(output string) (*fileWriter, error) {
	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &fileWriter{file: file, writer: bufio.NewWriter(file)}, nil
}

// Write writes data to the underlying file as a single line
func (w *fileWriter) Write(data []byte) error {
	// Write the entire line (data + newline) in one go to the buffer
	// This helps with atomicity at the application level
	line := make([]byte, 0, len(data)+1)
	line = append(line, data...)
	line = append(line, '\n')

	if _, err := w.writer.Write(line); err != nil {
		return err
	}
	// Periodically flush to avoid keeping too much in memory
	// and to ensure output is visible during long scans
	return w.writer.Flush()
}

// Close closes the underlying writer flushing everything to disk
func (w *fileWriter) Close() error {
	var flushErr error
	if w.writer != nil {
		flushErr = w.writer.Flush()
	}
	//nolint:errcheck // we don't care whether sync failed or succeeded.
	w.file.Sync()
	err := w.file.Close()
	if flushErr != nil {
		return flushErr
	}
	return err
}
