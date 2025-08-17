package tee

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
	localopt "github.com/yupsh/tee/opt"
)


// Flags represents the configuration options for the tee command
type Flags = localopt.Flags
// Command implementation
type command opt.Inputs[string, Flags]

// Tee creates a new tee command with the given parameters
func Tee(parameters ...any) yup.Command {
	return command(opt.Args[string, Flags](parameters...))
}

func (c command) Execute(ctx context.Context, input io.Reader, output, stderr io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Open all output files
	files := make([]*os.File, 0, len(c.Positional))
	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()

	for _, filename := range c.Positional {
		// Check for cancellation before each file
		if err := yup.CheckContextCancellation(ctx); err != nil {
			return err
		}

		var file *os.File
		var err error

		if bool(c.Flags.Append) {
			file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		} else {
			file, err = os.Create(filename)
		}

		if err != nil {
			fmt.Fprintf(stderr, "tee: %s: %v\n", filename, err)
			continue
		}

		files = append(files, file)
	}

	// Create list of all writers (stdout + files)
	writers := make([]io.Writer, 0, len(files)+1)
	writers = append(writers, output) // Always write to stdout
	for _, file := range files {
		writers = append(writers, file)
	}

	// Use TeeWriter for efficient copying to multiple destinations
	multiWriter := &TeeWriter{writers: writers}

	// Copy input to all outputs with context cancellation support
	_, err := yup.CopyWithContext(ctx, multiWriter, input)
	if err != nil {
		fmt.Fprintf(stderr, "tee: %v\n", err)
		return err
	}

	return nil
}

// TeeWriter implements io.Writer and writes to multiple destinations
type TeeWriter struct {
	writers []io.Writer
	mu      sync.Mutex
}

func (t *TeeWriter) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Write to all destinations
	// If any write fails, we continue with others but return the error
	var firstErr error
	written := len(p)

	for _, w := range t.writers {
		if _, err := w.Write(p); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return written, firstErr
}

// TeeReader creates a reader that also writes everything to the provided writers
// This could be useful for pipeline scenarios where tee needs to be transparent
func TeeReader(reader io.Reader, writers ...io.Writer) io.Reader {
	return &teeReader{
		reader:  reader,
		writers: writers,
	}
}

type teeReader struct {
	reader  io.Reader
	writers []io.Writer
}

func (t *teeReader) Read(p []byte) (int, error) {
	n, err := t.reader.Read(p)
	if n > 0 {
		// Write the read data to all tee outputs
		for _, w := range t.writers {
			w.Write(p[:n])
		}
	}
	return n, err
}
