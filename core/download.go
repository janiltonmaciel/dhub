package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

var defaultClient = &http.Client{
	Timeout: 15 * time.Second, // nolint
}

const (
	hostAPIGithub = "https://api.github.com"
	fileNameEtag  = "data.etag"

	etagFM    = os.FileMode(0644) // nolint
	homeDirFM = os.FileMode(0755) // nolint
)

type Download struct {
	Fs      afero.Fs
	HomeDir string
	Token   string
}

func NewDownload(token string) (d *Download, err error) {
	d = &Download{
		Fs:    afero.NewOsFs(),
		Token: token,
	}

	if d.HomeDir, err = d.createHomeDir(); err != nil {
		return nil, err
	}

	return d, d.extract()
}

func (d *Download) extract() error {
	if !d.needUpdate() {
		return nil
	}

	req, err := d.createRequest()
	if err != nil {
		return err
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		d.writeEtag(resp.Header.Get("ETAG"))
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error dowloanding, url: %s - status: %d", resp.Request.URL.String(), resp.StatusCode) // nolint
	}

	if err = Untar(d.HomeDir, resp.Body); err == nil {
		d.writeEtag(resp.Header.Get("ETAG"))
	}
	return err
}

func (d *Download) ReadLibraries(o interface{}) (err error) {
	return d.readYml("libraries.yml", o)
}

func (d *Download) ReadDataVersion(libraryName, o interface{}) (err error) {
	fileName := fmt.Sprintf("data/%s", libraryName)
	return d.readYml(fileName, o)
}

func (d *Download) readYml(fileName string, o interface{}) (err error) {
	var data []byte
	if data, err = d.readFile(fileName); err != nil {
		return err
	}
	return yaml.Unmarshal(data, o)
}

func (d *Download) createRequest() (*http.Request, error) {
	url := d.makeURL(
		"janiltonmaciel",
		"version-gen",
		"data/data.tar.gz",
	)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// add headers
	{
		etag := d.readLastEtag()
		req.Header.Set("If-None-Match", fmt.Sprintf(`W/%s`, etag))
		req.Header.Set("Authorization", fmt.Sprintf("token %s", d.Token))
		req.Header.Set("Accept", "application/vnd.github.v3.raw")
	}

	return req, nil
}

func (d *Download) writeEtag(etag string) {
	_ = afero.WriteFile(
		d.Fs,
		path.Join(d.HomeDir, fileNameEtag),
		[]byte(etag),
		etagFM,
	)
}

func (d *Download) readLastEtag() string {
	data, err := d.readFile(fileNameEtag)
	if err != nil {
		return ""
	}
	return string(data)
}

func (d *Download) createHomeDir() (home string, err error) {
	if home, err = os.UserHomeDir(); err != nil {
		return "", err
	}

	dir := path.Join(home, ".dhub")
	if _, err := d.Fs.Stat(dir); err == nil {
		return dir, nil
	}

	err = d.Fs.Mkdir(dir, homeDirFM)
	return dir, err
}

func (d *Download) readFile(fileName string) ([]byte, error) {
	filePath := path.Join(d.HomeDir, fileName)
	if _, err := d.Fs.Stat(filePath); err != nil {
		return nil, err
	}

	return ioutil.ReadFile(filePath)
}

func (d *Download) makeURL(owner, repo, path string) string {
	return fmt.Sprintf("%s/repos/%s/%s/contents/%s", hostAPIGithub, owner, repo, path)
}

func (d *Download) needUpdate() bool {
	filePath := path.Join(d.HomeDir, fileNameEtag)
	fileInfo, err := d.Fs.Stat(filePath)
	if err != nil {
		return true
	}

	modifiedtime := fileInfo.ModTime().UTC()
	now := time.Now()
	return !(modifiedtime.Year() == now.Year() && modifiedtime.Month() == now.Month() && modifiedtime.Day() == now.Day())
}
