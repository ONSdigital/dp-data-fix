package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ONSdigital/dp-data-fix/out"
	"github.com/spf13/cobra"
)

var (
	zebedeeFlag = "zebedee"
	targetFile  = "pdftables.pdf"
	outputFile  = "pdftables.csv"
	masterDir   = "master"
	pdfExt      = ".pdf"
	dataJson    = "data.json"
	cutoffDate  = time.Date(2018, 9, 1, 00, 00, 0, 0, time.UTC)
)

type Page struct {
	URI         string       `json:"uri"`
	Description *Description `json:"description"`
}

type Description struct {
	Title   string   `json:"title"`
	Contact *Contact `json:"contact"`
}

type Contact struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

type Row struct {
	URI       string
	Title     string
	Date      string
	Name      string
	Email     string
	Telephone string
}

func findPDFsCMD() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "pdfs",
		Short: "todo",
		RunE: func(cmd *cobra.Command, args []string) error {
			mPath, err := cmd.Flags().GetString(zebedeeFlag)
			if err != nil {
				return err
			}

			return FindPDFs(mPath)
		},
	}

	cmd.Flags().StringP(zebedeeFlag, "z", "", "The path for the root Zebedee directory (Required)")
	if err := cmd.MarkFlagRequired(zebedeeFlag); err != nil {
		return nil, err
	}

	return cmd, nil
}

func FindPDFs(zebedeeDir string) error {
	out.InfoF("Finding user generated PDFs under: %s", zebedeeDir)

	masterDir := filepath.Join(zebedeeDir, masterDir)

	if !fileExists(masterDir) {
		return fmt.Errorf("file does not exist %s", masterDir)
	}

	csvF, err := createCSV(outputFile)
	if err != nil {
		return err
	}

	defer csvF.Close()

	w := csv.NewWriter(csvF)
	if err = w.Write([]string{"URI", "Title", "Name", "Email", "Telephone", "Date"}); err != nil {
		return err
	}

	if err := filepath.Walk(masterDir, walkPDFs(w, masterDir)); err != nil {
		return err
	}

	w.Flush()

	out.InfoF("Generated results csv file: %s", outputFile)
	return err
}

func walkPDFs(w *csv.Writer, base string) filepath.WalkFunc {
	return func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(p) != pdfExt || info.Name() != targetFile {
			return nil
		}

		if info.ModTime().After(cutoffDate) {
			out.InfoF("PDF found: %s: %+v", info.Name(), info.ModTime())

			uri, err := filepath.Rel(base, filepath.Dir(p))
			if err != nil {
				return err
			}

			r := &Row{
				URI:       uri,
				Title:     "",
				Date:      info.ModTime().Format(time.RFC822),
				Name:      "",
				Email:     "",
				Telephone: "",
			}

			dataJson := filepath.Join(filepath.Dir(p), dataJson)
			if err = extractPageData(dataJson, r); err != nil {
				return err
			}

			r.Write(w)
		}

		return nil
	}
}

func extractPageData(p string, r *Row) error {
	if fileExists(p) {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		var p Page
		if err = json.Unmarshal(b, &p); err != nil {
			return err
		}

		if p.Description != nil {
			r.Title = p.Description.Title

			if p.Description.Contact != nil {
				r.Name = p.Description.Contact.Name
				r.Email = p.Description.Contact.Email
				r.Telephone = p.Description.Contact.Telephone
			}
		}
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

func createCSV(p string) (*os.File, error) {
	if fileExists(p) {
		if err := os.Remove(p); err != nil {
			return nil, err
		}
	}

	return os.Create(p)
}

func (r *Row) Write(w *csv.Writer) error {
	return w.Write([]string{r.URI, r.Title, r.Name, r.Email, r.Telephone, r.Date})
}
