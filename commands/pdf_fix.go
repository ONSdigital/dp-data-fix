package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ONSdigital/dp-data-fix/out"
	"github.com/spf13/cobra"
)

var (
	zebedeeFlag    = "zebedee"
	domainFlag     = "domain"
	filenameFlag   = "filename"
	masterDir      = "master"
	pdfExt         = ".pdf"
	dataJson       = "data.json"
	timeseries     = "/timeseries"
	pagePDF        = "page.pdf"
	zebedeeTimeFmt = "2006-01-02T15:04:05.000Z"
	cutoffDate     = time.Date(2018, 9, 1, 00, 00, 0, 0, time.UTC)
	headerRow      = []string{"URL", "Filename", "Title", "Name", "Email", "Telephone", "Release Date", "Last Modified Date", "PDF Table"}
)

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

			filename, err := cmd.Flags().GetString(filenameFlag)
			if err != nil {
				return err
			}

			return FindPDFs(zebedeeDir, host, filename)
		},
	}

	cmd.Flags().StringP(zebedeeFlag, "z", "", "The path for the root Zebedee directory (Required)")
	cmd.Flags().StringP(domainFlag, "d", "http://www.ons.gov.uk", "The host of the instance being queried")
	cmd.Flags().StringP(filenameFlag, "f", "user-generated-pdfs.csv", "The CVS file name")
	if err := cmd.MarkFlagRequired(zebedeeFlag); err != nil {
		return nil, err
	}

	return cmd, nil
}

func FindPDFs(zebedeeDir, host, filename string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	out.InfoF("Finding user generated PDFs under: %s", wd)

	masterDir := filepath.Join(zebedeeDir, masterDir)

	if !fileExists(masterDir) {
		return fmt.Errorf("file does not exist %s", masterDir)
	}

	csvF, err := createCSV(filepath.Join(zebedeeDir, filename))
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

	out.InfoF("Generated results csv file: %s", filename)
	return nil
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

		data := &Data{
			URL:            fmt.Sprintf("%s/file?uri=%s", host, uri),
			Filename:       info.Name(),
			Title:          "",
			Name:           "",
			Email:          "",
			Telephone:      "",
			ReleaseDateStr: "",
			ReleaseDate:    time.Now(),
			LastModDate:    info.ModTime(),
			LastModDateStr: info.ModTime().Format(time.RFC1123),
			IsPDFTable:     false,
		}

		dataJson := filepath.Join(filepath.Dir(p), dataJson)
		if err = extractPageData(dataJson, data); err != nil {
			return err
		}

		if data.IsBeforeCutoff() {
			return nil
		}

		out.InfoF("PDF found: %s: %+v", info.Name(), info.ModTime())

		rowData := []string{
			data.URL,
			data.Filename,
			data.Title,
			data.Name,
			data.Email,
			data.Telephone,
			data.ReleaseDateStr,
			data.LastModDateStr,
			strconv.FormatBool(data.IsPDFTable),
		}

		if err := w.Write(rowData); err != nil {
			return err
		}

		return nil
	}
}

func extractPageData(path string, data *Data) error {
	if fileExists(path) {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var page Page
		if err = json.Unmarshal(b, &page); err != nil {
			return err
		}

		if page.Description != nil {
			data.Title = page.Description.Title
			data.ReleaseDateStr = page.Description.ReleaseDate

			if data.ReleaseDateStr != "" {
				t, err := time.Parse(zebedeeTimeFmt, data.ReleaseDateStr)
				if err != nil {
					return err
				}

				data.ReleaseDateStr = t.Format(time.RFC1123)
				data.ReleaseDate = t
			}

			if page.Description.Contact != nil {
				data.Name = page.Description.Contact.Name
				data.Email = page.Description.Contact.Email
				data.Telephone = page.Description.Contact.Telephone

			}
		}

		if page.PDFTables != nil && len(page.PDFTables) >= 1 {
			for _, val := range page.PDFTables {
				if val.File == data.Filename {
					data.IsPDFTable = true
					break
				}
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

func (d *Data) IsBeforeCutoff() bool {
/*	if d.ReleaseDateStr != "" {
		return d.ReleaseDate.Before(cutoffDate)
	}*/

	return d.LastModDate.Before(cutoffDate)
}
