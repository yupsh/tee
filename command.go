package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[string, flags]

func Tee(parameters ...any) yup.Command {
	return command(yup.Initialize[string, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Open all output files
		var files []*os.File

		for _, filename := range p.Positional {
			var f *os.File
			var err error

			if bool(p.Flags.Append) {
				f, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			} else {
				f, err = os.Create(filename)
			}

			if err != nil {
				// Log error but continue
				fmt.Fprintf(stderr, "tee: %s: %v\n", filename, err)
				continue
			}

			files = append(files, f)
		}

		// Ensure files are closed when done
		defer func() {
			for _, f := range files {
				f.Close()
			}
		}()

		// Read from stdin and write to stdout and all files
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			line := scanner.Text()

			// Write to stdout
			if _, err := fmt.Fprintln(stdout, line); err != nil {
				return err
			}

			// Write to all files
			for _, f := range files {
				if _, err := fmt.Fprintln(f, line); err != nil {
					fmt.Fprintf(stderr, "tee: write error: %v\n", err)
				}
			}
		}

		return scanner.Err()
	}
}
