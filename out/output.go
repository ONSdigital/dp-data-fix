package out

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	infoC       = color.New(color.FgHiCyan, color.Italic)
	infoCItalic = color.New(color.FgHiCyan, color.Italic)
	errC        = color.New(color.FgHiRed, color.Bold)
	errCItalic  = color.New(color.FgHiRed, color.Italic)

	prefixC = color.New(color.FgHiGreen, color.Bold)
	prefix = "[dp-data-fix]"
)

func InfoF(fmt string, args ...interface{}) {
	cliPrefix(prefixC)
	highlight(infoC, fmt, args...)
}

func ErrF(fmt string, args ...interface{}) {
	cliPrefix(errC)
	highlight(errCItalic, fmt, args...)
}

func cliPrefix(c *color.Color) {
	c.Printf("%s ", prefix)
}

func highlight(c *color.Color, formattedMsg string, args ...interface{}) {
	var highlighted []interface{}

	for _, val := range args {
		highlighted = append(highlighted, c.SprintFunc()(val))
	}

	formattedMsg = fmt.Sprintf(formattedMsg, highlighted...)
	fmt.Printf("%s\n", formattedMsg)
}
