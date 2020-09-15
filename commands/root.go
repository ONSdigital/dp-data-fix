package commands

import (
	"github.com/spf13/cobra"
)

func GetRoot() *cobra.Command {
	r := &cobra.Command{
		Use:   "dp-data-fix",
		Short: "Run the data-fix tool",
	}

	r.AddCommand(findPDFsCMD())
	return r
}

