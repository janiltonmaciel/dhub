package core

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNotFound = errors.New("not found")

type Manager struct {
	Down *Download
}

func NewManager(token string) (*Manager, error) {
	down, err := NewDownload(token)
	if err != nil {
		return nil, err
	}

	return &Manager{
		Down: down,
	}, nil
}

func (m *Manager) GetLibraries() ([]Library, error) {
	var libraries []Library
	if err := m.Down.ReadLibraries(&libraries); err != nil {
		return nil, err
	}

	return libraries, nil
}

func (m *Manager) GetLibrary(libraryName string) (Library, error) {
	libraries, err := m.GetLibraries()
	if err != nil {
		return Library{}, err
	}

	for _, library := range libraries {
		if library.Name == libraryName {
			return library, nil
		}
	}

	return Library{}, fmt.Errorf("Library: %w", ErrNotFound)
}

func (m *Manager) FindLibrary(text string) ([]Library, error) {
	libraries, err := m.GetLibraries()
	if err != nil {
		return nil, err
	}

	data := make([]Library, 0)
	for _, library := range libraries {
		if strings.HasPrefix(library.Name, text) {
			data = append(data, library)
		}
	}

	return data, nil
}

func (m *Manager) GetLibraryVersion(library Library, text string, withPrerelease bool) (LibraryVersion, error) {
	var libraryVersion LibraryVersion
	if err := m.Down.ReadDataVersion(library.Name, &libraryVersion); err != nil {
		return LibraryVersion{}, err
	}

	tmp := libraryVersion.Versions[:0]
	for _, version := range libraryVersion.Versions {
		if text != "" && !strings.HasPrefix(version.Version, text) {
			continue
		}

		if withPrerelease {
			tmp = append(tmp, version)
		} else if !version.Prerelease {
			tmp = append(tmp, version)
		}
	}
	libraryVersion.Versions = tmp

	return libraryVersion, nil
}

// y := x[:0]
// for _, n := range x {
//     if n % 2 != 0 {
//         y = append(y, n)
//     }
// }

// func (l List) Render() error {
// 	fmt.Fprintln(l.config.Writer)
// 	fmt.Fprintf(l.config.Writer, "  %s  ", RenderYellow(strings.ToUpper(config.Library.Alias)))
// 	fmt.Fprintln(l.config.Writer)
// 	fmt.Fprintln(l.config.Writer)

// 	versions := FindVersions(config.Library, config.WithPrerelease, config.VersionName)
// 	if len(versions) == 0 {
// 		l.printNotFoundVersions(config)
// 		return nil
// 	}

// 	table := tablewriter.NewWriter(l.config.Writer)
// 	names := distributionNames(versions)

// 	var current string
// 	for _, version := range versions {
// 		current = " "
// 		// if version.Current {
// 		// 	current = "*"
// 		// }

// 		row := []string{fmt.Sprintf("%s%s", version.Version, current)}
// 		for _, name := range names {
// 			tags := []string{}
// 			for _, dist := range version.Distributions {
// 				if dist.Name == name {
// 					tags = append(tags, dist.ImageRepository)
// 				}
// 			}
// 			row = append(row, strings.Join(tags, "\n"))
// 		}
// 		if version.Prerelease {
// 			table.Rich(row, []tablewriter.Colors{tablewriter.Color(tablewriter.Bold, tablewriter.FgCyanColor)})
// 		} else {
// 			table.Append(row)
// 		}
// 	}

// 	colors := []tablewriter.Colors{
// 		tablewriter.Color(tablewriter.Bold, tablewriter.FgGreenColor),
// 	}
// 	headers := []string{"VERSION"}
// 	for _, name := range names {
// 		headers = append(headers, fmt.Sprintf("TAG (%s)", strings.ToUpper(name)))
// 		colors = append(colors, tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor})
// 	}

// 	if len(versions) > 10 {
// 		table.Rich(headers,
// 			[]tablewriter.Colors{
// 				tablewriter.Color(tablewriter.Normal, tablewriter.FgWhiteColor),
// 				tablewriter.Color(tablewriter.Normal, tablewriter.FgWhiteColor),
// 				tablewriter.Color(tablewriter.Normal, tablewriter.FgWhiteColor),
// 			},
// 		)
// 	}

// 	table.SetHeader(headers)
// 	table.SetColumnColor(colors...)
// 	table.SetAlignment(tablewriter.ALIGN_LEFT)
// 	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
// 	table.SetBorder(false)
// 	table.SetRowLine(true)
// 	table.Render()

// 	fmt.Fprintln(l.config.Writer)
// 	return nil
// }

// func (l List) printNotFoundVersions(config ListConfig) {
// 	fmt.Fprintf(l.config.Writer,
// 		"     %s `%s`\n",
// 		RenderRed("not found matching versions "),
// 		RenderYellow(config.VersionName))
// 	fmt.Fprintln(l.config.Writer)
// }

// func distributionNames(versions []Version) (names []string) {
// 	exist := make(map[string]bool)
// 	for _, version := range versions {
// 		for _, dist := range version.Distributions {
// 			if _, found := exist[dist.Name]; !found {
// 				exist[dist.Name] = true
// 				names = append(names, dist.Name)
// 			}
// 		}
// 	}
// 	sort.Strings(names)
// 	return names
// }
