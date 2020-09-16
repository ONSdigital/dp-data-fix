package commands

import (
	"github.com/spf13/cobra"
)

func NewCli() (*cobra.Command, error) {
	r := &cobra.Command{
		Use:   "dp-data-fix",
		Short: "Run the data-fix tool",
	}

	fixPDFsCMD, err := findPDFsCMD()
	if err != nil {
		return nil, err
	}

	r.AddCommand(fixPDFsCMD)
	return r, nil
}

