package tee_test

import (
	"context"
	"os"
	"strings"

	"github.com/yupsh/tee"
	"github.com/yupsh/tee/opt"
)

func ExampleTee() {
	ctx := context.Background()
	input := strings.NewReader("Hello World\nSecond Line\n")

	cmd := tee.Tee("output.txt")
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output: Hello World
	// Second Line
	// (Also writes to output.txt)
}

func ExampleTee_append() {
	ctx := context.Background()
	input := strings.NewReader("Appended line\n")

	cmd := tee.Tee("log.txt", opt.Append)
	cmd.Execute(ctx, input, os.Stdout, os.Stderr)
	// Output: Appended line
	// (Also appends to log.txt)
}
