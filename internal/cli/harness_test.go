package cli

import (
	"bytes"
	"testing"
)

// execute runs the root command with args, capturing stdout/stderr
// separately, and returns the error Execute() produced (nil on success).
func execute(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	cmd := NewRootCmd()
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return outBuf.String(), errBuf.String(), err
}
