package cmd_test

import (
	"testing"

	"github.com/44smkn/kubectl-role-diff/pkg/cmd"
)

func TestFormatVersion(t *testing.T) {
	expects := "kubectl-role-diff version 0.1.0 (2021-11-30)\n"
	if got := cmd.FormatVersion("0.1.0", "2021-11-30"); got != expects {
		t.Errorf("Format() = %q, wants %q", got, expects)
	}
}
