package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ONSdigital/dp-data-fix/out"
	"github.com/spf13/cobra"
)

var (
	zebedeeFlag    = "zebedee"
	domainFlag     = "domain"
	outputFile     = "pdftables.csv"
	masterDir      = "master"
	pdfExt         = ".pdf"
	dataJson       = "data.json"
	timeseries     = "/timeseries"
	pagePDF        = "page.pdf"
	zebedeeTimeFmt = "2006-01-02T15:04:05.000Z"
	cutoffDate     = time.Date(2018, 9, 1, 00, 00, 0, 0, time.UTC) //.Add(time.Duration(-1) * time.Millisecond)
	headerRow      = []string{"URL", "Filename", "Title", "Name", "Email", "Telephone", "Release Date", "Last Modified Date"}
)

type Page struct {
	URI         string       `json:"uri"`
	Description *Description `json:"description"`
}

type Description struct {
	ReleaseDate string   `json:"releaseDate"`
	Title       string   `json:"title"`
	Contact     *Contact `json:"contact"`
}

type Contact struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

type Row struct {
	URL            string
	Filename       string
	Title          string
	ReleaseDateStr string
	ReleaseDate    time.Time
	LastModDate    time.Time
	Name           string
	Email          string
	Telephone      string
}

func findPDFsCMD() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "pdfs",
		Short: "todo",
		RunE: func(cmd *cobra.Command, args []string) error {
			zebedeeDir, err := cmd.Flags().GetString(zebedeeFlag)
			if err != nil {
				return err
			}

			host, err := cmd.Flags().GetString(domainFlag)
			if err != nil {
				return err
			}

			return FindPDFs(zebedeeDir, host)
		},
	}

	cmd.Flags().StringP(zebedeeFlag, "z", "", "The path for the root Zebedee directory (Required)")
	cmd.Flags().StringP(domainFlag, "d", "http://www.ons.gov.uk", "The host of the instance being queried")
	if err := cmd.MarkFlagRequired(zebedeeFlag); err != nil {
		return nil, err
	}

	return cmd, nil
}

func FindPDFs(zebedeeDir, host string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	out.InfoF("Finding user generated PDFs under: %s", wd)

	masterDir := filepath.Join(zebedeeDir, masterDir)

	if !fileExists(masterDir) {
		return fmt.Errorf("file does not exist %s", masterDir)
	}

	csvF, err := createCSV(filepath.Join(zebedeeDir, outputFile))
	if err != nil {
		return err
	}

	defer csvF.Close()

	w := csv.NewWriter(csvF)
	if err = w.Write(headerRow); err != nil {
		return err
	}

	out.InfoF("searching for user generated PDFs in %s", masterDir)
	if err := filepath.Walk(masterDir, walkPDFs(w, host, masterDir)); err != nil {
		return err
	}

	w.Flush()

	out.InfoF("Generated results csv file: %s", outputFile)
	return err
}

func walkPDFs(w *csv.Writer, host, base string) filepath.WalkFunc {
	return func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip if matches any of the following criteria
		if info.IsDir() || strings.Contains(timeseries, p) || filepath.Ext(p) != pdfExt || pagePDF == info.Name() {
			return nil
		}

		uri, err := filepath.Rel(base, p)
		if err != nil {
			return err
		}

		r := &Row{
			URL:            fmt.Sprintf("%s/file?uri=%s", host, uri),
			Filename:       info.Name(),
			Title:          "",
			ReleaseDateStr: "",
			ReleaseDate:    time.Now(),
			LastModDate:    info.ModTime(), //.Format(time.RFC822),
			Name:           "",
			Email:          "",
			Telephone:      "",
		}

		dataJson := filepath.Join(filepath.Dir(p), dataJson)
		if err = extractPageData(dataJson, r); err != nil {
			return err
		}

		isBefore, err := IsBeforeCutoff(r)
		if err != nil {
			return err
		}

		if isBefore {
			return nil
		}

		out.InfoF("PDF found: %s: %+v", info.Name(), info.ModTime())

		releaseDate := ""
		if r.ReleaseDateStr != "" {
			releaseDate = r.ReleaseDate.Format(time.RFC1123)
		}

		lostModifiedDate := r.LastModDate.Format(time.RFC1123)

		rowData := []string{r.URL, r.Filename, r.Title, r.Name, r.Email, r.Telephone, releaseDate, lostModifiedDate}
		if err := w.Write(rowData); err != nil {
			return err
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
			r.ReleaseDateStr = p.Description.ReleaseDate

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

	out.InfoF("creating output csv file %s", p)
	return os.Create(p)
}

func IsBeforeCutoff(r *Row) (bool, error) {
	if r.ReleaseDateStr != "" {
		relDate, err := time.Parse(zebedeeTimeFmt, r.ReleaseDateStr)
		if err != nil {
			return false, err
		}

		r.ReleaseDate = relDate
		return relDate.Before(cutoffDate), nil
	}

	return r.LastModDate.Before(cutoffDate), nil
}
