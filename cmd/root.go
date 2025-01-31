package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/kplachkov/helm-mirror/service"
)

var (
	// Verbose defines if the command is being run with verbose mode
	Verbose bool
	// IgnoreErrors ignores errors in processing charts
	IgnoreErrors bool
	// AllVersions gets all the versions of the charts when true, false by default
	AllVersions  bool
	chartName    string
	chartVersion string
	folder       string
	flags        = log.Ldate | log.Lmicroseconds | log.Lshortfile
	prefix       = "helm-mirror: "
	logger       *log.Logger
	username     string
	password     string
	caFile       string
	certFile     string
	keyFile      string
	newRootURL   string
)

const rootDesc = `Mirror Helm Charts from an index file into a local folder.

For example:

helm mirror https://yourorg.com/charts /yourorg/charts

This will download the index file and the charts into
the folder indicated.

The index file is a yaml that contains a list of
charts in this format. Example:

	apiVersion: v1
	entries:
	  chart:
	  - apiVersion: 1.0.0
	    created: 2018-08-08T00:00:00.00000000Z
	    description: A Helm chart for your application
	    digest: 3aa68d6cb66c14c1fcffc6dc6d0ad8a65b90b90c10f9f04125dc6fcaf8ef1b20
	    name: chart
	    urls:
	    - https://kubernetes-charts.yourorganization.com/chart-1.0.0.tgz
	  chart2:
	  - apiVersion: 1.0.0
	    created: 2018-08-08T00:00:00.00000000Z
	    description: A Helm chart for your application
	    digest: 7ae62d60b61c14c1fcffc6dc670e72e62b91b91c10f9f04125dc67cef2ef0b21
	    name: chart
	    urls:
	    - https://kubernetes-charts.yourorganization.com/chart2-1.0.0.tgz

This will download these charts

	https://kubernetes-charts.yourorganization.com/chart-1.0.0.tgz
	https://kubernetes-charts.yourorganization.com/chart2-1.0.0.tgz

into your destination folder.`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mirror [Repo URL] [Destination Folder]",
	Short: "Mirror Helm Charts from an index file into a local folder.",
	Long:  rootDesc,
	Args:  validateRootArgs,
	RunE:  runRoot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	logger = log.New(os.Stdout, prefix, flags)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&IgnoreErrors, "ignore-errors", "i", false, "ignores errors while downloading or processing charts")
	rootCmd.PersistentFlags().BoolVarP(&AllVersions, "all-versions", "a", false, "gets all the versions of the charts in the chart repository")
	rootCmd.Flags().StringVar(&chartName, "chart-name", "", "name of the chart that gets mirrored")
	rootCmd.Flags().StringVar(&chartVersion, "chart-version", "", "specific version of the chart that is going to be mirrored")
	rootCmd.Flags().StringVar(&username, "username", "", "chart repository username")
	rootCmd.Flags().StringVar(&password, "password", "", "chart repository password")
	rootCmd.Flags().StringVar(&caFile, "ca-file", "", "verify certificates of HTTPS-enabled servers using this CA bundle")
	rootCmd.Flags().StringVar(&certFile, "cert-file", "", "identify HTTPS client using this SSL certificate file")
	rootCmd.Flags().StringVar(&keyFile, "key-file", "", "identify HTTPS client using this SSL key file")
	rootCmd.Flags().StringVar(&newRootURL, "new-root-url", "", "New root url of the chart repository (eg: `https://mirror.local.lan/charts`)")
	rootCmd.AddCommand(newVersionCmd())
}

func validateRootArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		if len(args) == 1 && args[0] == "help" {
			return nil
		}
		logger.Printf("error: requires at least two args to execute")
		return errors.New("error: requires at least two args to execute")
	}

	repoURL, err := url.Parse(args[0])
	if err != nil {
		logger.Printf("error: not a valid URL for index file: %s", err)
		return err
	}

	if !strings.Contains(repoURL.Scheme, "http") {
		logger.Printf("error: not a valid URL protocol: `%s`", repoURL.Scheme)
		return errors.New("error: not a valid URL protocol")
	}
	if !path.IsAbs(args[1]) {
		logger.Printf("error: please provide a full path for destination folder: `%s`", args[1])
		return errors.New("error: please provide a full path for destination folder")
	}
	return nil
}

func runRoot(cmd *cobra.Command, args []string) error {
	repoURL, err := url.Parse(args[0])
	if err != nil {
		logger.Printf("error: not a valid URL for index file: %s", err)
		return err
	}

	folder = args[1]
	err = os.MkdirAll(folder, 0744)
	if err != nil {
		logger.Printf("error: cannot create destination folder: %s", err)
		return err
	}

	rootURL := &url.URL{}
	if newRootURL != "" {
		rootURL, err = url.Parse(newRootURL)
		if err != nil {
			logger.Printf("error: new-root-url not a valid URL: %s", err)
			return err
		}

		if !strings.Contains(rootURL.Scheme, "http") {
			logger.Printf("error: new-root-url not a valid URL protocol: `%s`", rootURL.Scheme)
			return errors.New("error: new-root-url not a valid URL protocol")
		}
	}

	if chartVersion != "" && chartName == "" {
		logger.Printf("error: chart Version depends on a chart name, please specify one")
		return errors.New("error: chart Version depends on a chart name, please specify one")
	}

	config := repo.Entry{
		Name:     folder,
		URL:      repoURL.String(),
		Username: username,
		Password: password,
		CAFile:   caFile,
		CertFile: certFile,
		KeyFile:  keyFile,
	}

	getService := service.NewGetService(config, AllVersions, Verbose, IgnoreErrors, logger, rootURL.String(), chartName, chartVersion)
	err = getService.Get()
	return err
}
