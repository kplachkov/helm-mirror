package service

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"

	"github.com/kplachkov/helm-mirror/formatter"
)

// ImagesServiceInterface defines a Get service
type ImagesServiceInterface interface {
	Images() error
}

// ImagesService structure definition
type ImagesService struct {
	target         string
	formatter      formatter.Formatter
	verbose        bool
	ignoreErrors   bool
	exitWithErrors bool
	logger         *log.Logger
	buffer         bytes.Buffer
}

// NewImagesService return a new instance of ImagesService
func NewImagesService(target string, verbose bool, ignoreErrors bool, formatter formatter.Formatter, logger *log.Logger) ImagesServiceInterface {
	return &ImagesService{
		target:       target,
		formatter:    formatter,
		logger:       logger,
		verbose:      verbose,
		ignoreErrors: ignoreErrors,
	}
}

// Images extracts al the images in the Helm Charts downloaded by the get command
func (i *ImagesService) Images() error {
	fi, err := os.Stat(i.target)
	if err != nil {
		i.logger.Printf("error: cannot read target: %s", i.target)
		return err
	}

	if fi.IsDir() {
		err = i.processDirectory(i.target)
	} else {
		err = i.processTarget(i.target)
	}
	if err != nil {
		i.logger.Printf("error: processing target %s: %s", i.target, err)
		return err
	}
	err = i.formatter.Output(i.buffer)
	if err != nil {
		i.logger.Printf("writing output: %s", err)
		return err
	}
	return nil
}

func (i *ImagesService) processDirectory(target string) error {
	hasTgzCharts := false
	fi, err := os.Stat(i.target)
	if err != nil {
		i.logger.Printf("error: cannot read target: %s", i.target)
		return err
	}
	if !fi.IsDir() {
		return errors.New("error: inspectImages: processDirectory: target not a directory")
	}
	e := i.processTarget(target)
	if e != nil {
		err := filepath.Walk(target, func(dir string, info os.FileInfo, err error) error {
			if err != nil {
				i.logger.Printf("error: cannot access a dir %q: %v\n", dir, err)
				return err
			}
			if !info.IsDir() && strings.Contains(info.Name(), ".tgz") {
				hasTgzCharts = true
				err := i.processTarget(path.Join(target, info.Name()))
				if err != nil && i.ignoreErrors {
					i.exitWithErrors = true
				} else if err != nil {
					i.logger.Printf("error: cannot load chart: %s", err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			i.logger.Printf("error walking the path %q: %v\n", target, err)
			return err
		}
	}
	if e != nil && !hasTgzCharts {
		i.logger.Printf("error: cannot load chart: %s", e)
		return e
	}
	return nil
}

func (i *ImagesService) processTarget(target string) error {
	if i.verbose {
		i.logger.Printf("processing target: %s", target)
	}

	cht, err := loader.Load(target)
	if err != nil {
		return err
	}

	vals, err := chartutil.ToRenderValues(
		cht,
		nil,
		chartutil.ReleaseOptions{},
		chartutil.DefaultCapabilities,
	)
	if err != nil {
		i.logger.Printf("error: cannot render values: %s", err)
		return err
	}

	vals = cleanUp(vals)

	renderer := engine.Engine{
		LintMode: i.ignoreErrors,
	}

	rendered, err := renderer.Render(cht, vals)
	if err != nil {
		i.logger.Printf("error: cannot render chart: %s", err)
		return err
	}

	for _, t := range rendered {
		scanner := bufio.NewScanner(strings.NewReader(t))
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "image:") {
				im := sanitizeImageString(scanner.Text())
				i.buffer.WriteString(im + "\n")
			}
		}
	}
	return nil
}

func sanitizeImageString(s string) string {
	s = strings.Replace(s, "\"", "", 2)
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "-")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "image: ")
	s = strings.TrimSpace(s)
	return s
}

func cleanUp(i map[string]interface{}) map[string]interface{} {
	for n, v := range i {
		if reflect.TypeOf(v) == reflect.TypeOf(map[string]interface{}{}) {
			i[n] = cleanUp(v.(map[string]interface{}))
		} else if reflect.TypeOf(v) == reflect.TypeOf(chartutil.Values{}) {
			i[n] = cleanUp(v.(chartutil.Values))
		} else if v == nil {
			i[n] = ""
		} else {
			i[n] = v
		}
	}
	return i
}
