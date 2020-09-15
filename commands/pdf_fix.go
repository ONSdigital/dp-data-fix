package commands

import (
	"fmt"
	"os"

	"github.com/ONSdigital/dp-data-fix/out"
	"github.com/spf13/cobra"
)

func findPDFsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pdfs",
		Short: "todo",
		RunE: func(cmd *cobra.Command, args []string) error {
			mPath, err := cmd.Flags().GetString("master")
			if err != nil {
				return err
			}

			return findPDFs(mPath)
		},
	}

	cmd.Flags().StringP("master", "m", "", "The path for Zebedee Master directory (Required)")
	cmd.MarkFlagRequired("master")

	return cmd
}

func findPDFs(mPath string) error {
	out.InfoF("Finding user generated PDFs under: %s", mPath)

	if !fileExists(mPath) {
		return fmt.Errorf("file does not exist %s", mPath)
	}

	return nil
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
