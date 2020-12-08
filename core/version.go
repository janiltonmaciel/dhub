package core

type (
	Distribution struct {
		Name            string   `yaml:"name"`
		ReleaseName     string   `yaml:"releaseName"`
		Release         float32  `yaml:"release"`
		Image           string   `yaml:"image"`
		Weight          int      `yaml:"weight"`
		Tags            []string `yaml:"tags"`
		URLDockerfile   string   `yaml:"urlDockerfile"`
		ImageRepository string   `yaml:"imageRepository"`

		Library       Library
		RepositoryURL string
	}

	Version struct {
		Version       string                    `yaml:"version"`
		// MajorVersion  string                    `yaml:"majorVersion"`
		// Latest        bool                      `yaml:"latest"`
		Prerelease    bool                      `yaml:"prerelease"`
		Date          string                    `yaml:"date"`
		// Current       bool                      `yaml:"current"`
		Distributions map[string][]Distribution `yaml:"distributions"`
	}

	LibraryVersion struct {
		Name          string    `yaml:"name"`
		HubURL        string    `yaml:"urlHub"`
		RepositoryURL string    `yaml:"urlRepository"`
		Versions      []Version `yaml:"version"`
	}
)

  - version: 4.7.7
    prerelease: false
    date: '2020-05-11T13:40:25'
    distributions:
      - alpine:


// func (lv *LibraryVersion) GetDistros() []Distribution {

// }
