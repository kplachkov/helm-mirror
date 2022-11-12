package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is version ID for the source, read from VERSION in the source and
// populated on build by make.
var version = "unkwown"

// gitCommit is the commit hash that the binary was built from and will be
// populated on build by make.
var gitCommit = ""

// versionCmd represents the version command
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version of the helm mirror plugin",
		Run:   runVersion,
	}
}

func runVersion(*cobra.Command, []string) {
	v := ""
	if version != "" {
		v = version
	}
	if gitCommit != "" {
		v = fmt.Sprintf("%s~git%s", v, gitCommit)
	}
	fmt.Println(v)
}
