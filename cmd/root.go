package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func Execute(version, commit, date, token string) error {
	app := cli.NewApp()
	app.Name = RenderGreen("dhub")
	app.Version = RenderGreen(version)
	app.Authors = createAuthors()
	app.Commands = createCommands()
	app.Metadata = map[string]interface{}{
		"token": token,
	}

	cli.VersionPrinter = versionPrinter(commit, date)
	return app.Run(os.Args)
}

func createCommands() []*cli.Command {
	return []*cli.Command{
		NewCommandList(),
	}
}

func createAuthors() []*cli.Author {
	return []*cli.Author{
		{
			Name:  "Janilton Maciel",
			Email: "janilton@gmail.com",
		},
	}
}

func versionPrinter(commit, date string) func(c *cli.Context) {
	return func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "version: %s\n", c.App.Version)
		fmt.Fprintf(c.App.Writer, "commit: %s\n", commit)
		fmt.Fprintf(c.App.Writer, "date: %s\n", date)
		fmt.Fprintf(c.App.Writer, "author: %s <%s>\n", c.App.Authors[0].Name, c.App.Authors[0].Email)
	}
}
