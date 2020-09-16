package main

import (
	"os"

	"github.com/ONSdigital/dp-data-fix/commands"
	"github.com/ONSdigital/dp-data-fix/out"
)

func main() {
	if err := run(); err != nil {
		out.ErrF("cli error: %s", err.Error())
		out.ErrF("exiting\n")
		os.Exit(1)
	}

}

func run() error {
	root, err := commands.NewCli()
	if err != nil {
		return err
	}

	return root.Execute()
}
