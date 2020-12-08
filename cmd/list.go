package cmd

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/janiltonmaciel/dhub/core"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

type Config struct {
	LibraryName    string
	VersionName    string
	WithPrerelease bool
	WithTags       bool
	Writer         io.Writer
	Token          string
}

type list struct {
	conf Config
	m    *core.Manager
}

func NewCommandList() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List versions available for docker language",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "pre-release",
				Usage: "Show pre-release versions",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Show distributions release",
			},
		},
		Action: func(c *cli.Context) error {
			libraryName := strings.TrimSpace(c.Args().Get(0))
			if strings.ToLower(libraryName) == "help" {
				return showCommandHelp(c)
			}

			versionName := strings.TrimSpace(c.Args().Get(1))
			withPrerelease := c.Bool("pre-release")
			withTags := c.Bool("verbose")

			config := Config{
				LibraryName:    libraryName,
				VersionName:    versionName,
				WithPrerelease: withPrerelease,
				WithTags:       withTags,
				Writer:         c.App.Writer,
				Token:          c.App.Metadata["token"].(string),
			}

			list, err := NewList(config)
			if err != nil {
				return err
			}
			return list.Render()
		},
	}
}

func NewList(conf Config) (*list, error) {
	m, err := core.NewManager(conf.Token)
	if err != nil {
		return nil, err
	}

	l := &list{
		conf: conf,
		m:    m,
	}

	return l, nil
}

func (l *list) Render() error {
	if l.conf.LibraryName == "" {
		return l.renderLibraries()
	}

	library, err := l.m.GetLibrary(l.conf.LibraryName)
	if err != nil {
		return nil
	}

	return l.renderVersions(library)
}

func (l *list) renderVersions(library core.Library) error {
	// libraryVersion, err := l.m.GetLibraryVersion(library, l.conf.VersionName, l.conf.WithPrerelease)
	// if err != nil {
	// 	return err
	// }
	// libraryVersion.Versions

	return nil
}

// GetLibraryVersion

func (l *list) renderLibraries() error {
	libraries, err := l.m.GetLibraries()
	if err != nil {
		return err
	}

	sort.Slice(libraries, func(i, j int) bool {
		return libraries[i].PullCount < libraries[j].PullCount
	})

	rows := make([][]string, 0)
	for _, library := range libraries {
		row := []string{
			RenderGreen(library.Name),
			library.Description,
			library.LastUpdated.Format("2006-01-02"),
			FormatNumber.Sprint(library.PullCount),
		}
		rows = append(rows, row)
	}

	headers := []string{
		"Name",
		"Description",
		"Last Updated",
		"Pull Count",
	}
	l.tableRender(rows, headers, headers)

	return nil
}

func (l *list) tableRender(rows [][]string, header []string, footer []string) {
	fmt.Fprintln(l.conf.Writer)

	table := tablewriter.NewWriter(l.conf.Writer)
	table.SetHeader(header)
	table.SetFooter(footer)
	table.SetBorder(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.AppendBulk(rows)
	table.Render()

	fmt.Fprintln(l.conf.Writer)
}

func showCommandHelp(c *cli.Context) error {
	fmt.Fprintln(c.App.Writer, RenderRed("X Incorrect usage!"))
	fmt.Fprintln(c.App.Writer)
	return cli.ShowCommandHelp(c, c.Command.Name)
}

func showLibraryCommandHelp(c *cli.Context, libraryName string) error {
	msg := fmt.Sprintf("%s %s",
		RenderRed("X Libray not found:"),
		RenderYellow(libraryName),
	)
	fmt.Fprintln(c.App.Writer, msg)
	fmt.Fprintln(c.App.Writer)
	return cli.ShowCommandHelp(c, "libraries")
}

// dhub list -> exibe todas libraries
// dhub list redis -> exibe todas as versoes da library redis
// dhub list redis 5 -> exibe todas libraries
// dhub info redis 5454545

// dhub library -> exibe todas libraries
// dhub library redis -> exibe todas libraries
// dhub version redis -> exibe todas as versoes da library redis
// dhub version redis 5 -> exibe todas as versoes que iniciadl com 5 da library redis
// dhub show redis 5454545

// dhub list -> exibe todas libraries
// dhub list [text] -> exibe todas libraries
// dhub version <lib> -> exibe todas as versoes da library redis
// dhub version <lib> [version] -> exibe todas as versoes que iniciadl com 5 da library redis
// dhub show redis 5454545
