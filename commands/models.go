package commands

import "time"

type Page struct {
	URI         string       `json:"uri"`
	Description *Description `json:"description"`
	PDFTables   []*PDFTable  `json:"pdfTable,omitempty"`
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

type PDFTable struct {
	Title string `json:"title"`
	File  string `json:"file"`
}

type Data struct {
	URL            string
	Filename       string
	Title          string
	Name           string
	Email          string
	Telephone      string
	ReleaseDate    time.Time
	ReleaseDateStr string
	LastModDate    time.Time
	LastModDateStr string
	IsPDFTable     bool
}
