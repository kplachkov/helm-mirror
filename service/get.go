package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

const (
	indexFileName = "index.yaml"
)

// GetServiceInterface defines a Get service
type GetServiceInterface interface {
	Get() error
}

// GetService structure definition
type GetService struct {
	config        repo.Entry
	verbose       bool
	ignoreErrors  bool
	logger        *log.Logger
	newRootURL    string
	allVersions   bool
	chartName     string
	chartVersion  string
	indexFilePath string
}

// NewGetService return a new instance of GetService
func NewGetService(config repo.Entry, allVersions bool, verbose bool, ignoreErrors bool, logger *log.Logger, newRootURL string, chartName string, chartVersion string) GetServiceInterface {
	return &GetService{
		config:       config,
		verbose:      verbose,
		ignoreErrors: ignoreErrors,
		logger:       logger,
		newRootURL:   newRootURL,
		allVersions:  allVersions,
		chartName:    chartName,
		chartVersion: chartVersion,
	}
}

// Get methods downloads the index file and the Helm charts to the working directory.
func (g *GetService) Get() error {
	chartRepo, err := repo.NewChartRepository(&g.config, getter.All(&cli.EnvSettings{}))
	if err != nil {
		return err
	}

	g.indexFilePath, err = chartRepo.DownloadIndexFile()
	if err != nil {
		return err
	}

	chartRepo.IndexFile, err = repo.LoadIndexFile(g.indexFilePath)
	if err != nil {
		return err
	}

	index := search.NewIndex()
	index.AddRepo(chartRepo.Config.Name, chartRepo.IndexFile, g.allVersions || g.chartVersion != "")

	chartNameRegex := fmt.Sprintf("^.*%s.*", g.chartName)
	results, err := index.Search(chartNameRegex, 1, true)
	if err != nil {
		return err
	}

	for _, res := range results {
		if g.chartName != "" && res.Chart.Name != g.chartName {
			continue
		}
		if g.chartVersion != "" && res.Chart.Version != g.chartVersion {
			continue
		}

		for _, u := range res.Chart.URLs {
			b, err := chartRepo.Client.Get(u)
			if err != nil {
				if g.ignoreErrors {
					g.logger.Printf("WARNING: processing chart %s(%s) - %s", res.Name, res.Chart.Version, err)
					continue
				} else {
					return err
				}
			}

			chartFileName := fmt.Sprintf("%s-%s.tgz", res.Chart.Name, res.Chart.Version)
			chartPath := path.Join(g.config.Name, chartFileName)

			err = g.writeFile(chartPath, b.Bytes())
			if err != nil {
				return err
			}
		}
	}

	err = g.prepareIndexFile()
	return err
}

func (g *GetService) writeFile(name string, content []byte) error {
	err := os.WriteFile(name, content, 0666)
	if g.ignoreErrors {
		g.logger.Printf("cannot write files %s: %s", name, err)
	} else {
		return err
	}
	return nil
}

func (g *GetService) prepareIndexFile() error {
	indexPath := path.Join(g.config.Name, indexFileName)

	if g.newRootURL != "" {
		indexContent, err := os.ReadFile(g.indexFilePath)
		if err != nil {
			return err
		}

		content := bytes.Replace(indexContent, []byte(g.config.URL), []byte(g.newRootURL), -1)

		err = g.writeFile(g.indexFilePath, content)
		if err != nil {
			return err
		}
	}

	return moveFile(g.indexFilePath, indexPath)
}

func moveFile(src string, dst string) error {
	if src == dst {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return err
	}

	defer func() {
		os.Remove(src)
	}()

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.Chmod(dst, si.Mode())
	return err
}
