package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdVersion(version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(os.Stdout, FormatVersion(version, buildDate))
		},
	}

	return cmd
}

func FormatVersion(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	var dateStr string
	if buildDate != "" {
		dateStr = fmt.Sprintf(" (%s)", buildDate)
	}

	return fmt.Sprintf("kubectl-role-diff version %s%s\n", version, dateStr)
}
