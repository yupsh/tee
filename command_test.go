package command_test

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/tee"
)

func TestTee_Stdout(t *testing.T) {
	result := run.Command(command.Tee()).
		WithStdinLines("a", "b", "c").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"a", "b", "c"})
}

func TestTee_Append(t *testing.T) {
	result := run.Command(command.Tee(command.Append)).
		WithStdinLines("test").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"test"})
}

func TestTee_EmptyInput(t *testing.T) {
	result := run.Quick(command.Tee())
	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

func TestTee_InputError(t *testing.T) {
	result := run.Command(command.Tee()).
		WithStdinError(errors.New("read failed")).Run()
	assertion.ErrorContains(t, result.Err, "read failed")
}

