# helm-mirror plugin

![License: Apache 2.0](https://img.shields.io/github/license/kplachkov/helm-mirror.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/kplachkov/helm-mirror)](https://goreportcard.com/report/github.com/kplachkov/helm-mirror)

Helm plugin used to mirror Chart repositories

## Usage

Mirror Helm Charts from an index file into a local folder.

For example:

```shell
helm-mirror https://yourorg.com/charts /yourorg/charts
```

This will download the index file and the latest version of the charts
into the folder indicated.

The index file is a yaml that contains a list of
charts in this format. Example:

```yaml
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
```

This will download these charts:

- https://kubernetes-charts.yourorganization.com/chart-1.0.0.tgz

- https://kubernetes-charts.yourorganization.com/chart2-1.0.0.tgz

into your destination folder.

Usage:

```
  helm-mirror [Repo URL] [Destination Folder] [flags]
  helm-mirror [command]
```

Available Commands:

```
  help           Help about any command
  inspect-images Extract all the images of the Helm Charts.
  version        Show version of the helm-mirror plugin
```

Flags:

```
  -a, --all-versions                                   gets all the versions of the charts in the chart repository
      --ca-file string                                 verify certificates of HTTPS-enabled servers using this CA bundle
      --cert-file string                               identify HTTPS client using this SSL certificate file
      --chart-name string                              name of the chart that gets mirrored
      --chart-version string                           specific version of the chart that is going to be mirrored
  -h, --help                                           help for mirror
  -i, --ignore-errors                                  ignores errors while downloading or processing charts
      --key-file string                                identify HTTPS client using this SSL key file
      --new-root-url https://mirror.local.lan/charts   New root url of the chart repository (eg: https://mirror.local.lan/charts)
      --password string                                chart repository password
      --username string                                chart repository username
  -v, --verbose                                        verbose output
```

### Getting all charts

```shell
helm-mirror https://yourorg.com/charts /yourorg/charts --all-charts
```

This will download all the charts and all the available versions
of the charts.

### Getting one specific chart

```shell
helm-mirror https://yourorg.com/charts /yourorg/charts --chart-name nginx
```

This will download the latest version of the chart `nginx`.

### Getting one specific chart with specific version

```shell
helm-mirror https://yourorg.com/charts /yourorg/charts --chart-name nginx --chart-version 2.14.3
```

This will download the version `2.14.3` of the chart `nginx`.

Use `helm-mirror [command] --help` for more information about a command.

## Commands

### inspect-images

Extract all the container images listed in each Helm Chart or
the Helm Charts in the folder provided. This command dumps
the images on `stdout` by default, for more options check
`output flag`. Example:

- `helm-mirror inspect-images /tmp/helm`

- `helm-mirror inspect-images /tmp/helm/app.tgz`

The [folder|tgzfile] has to be a full path.

#### Usage

```
helm-mirror inspect-images [folder|tgzfile] [flags]
```

#### Flags

```
  -h, --help               help for inspect-images

  -i, --ignore-errors      ignores errors whiles processing charts. (Exit Code: 2)

  -o, --output string      choose an output for the list of images.(default "stdout")
```

- `file`: outputs all images to a file
- `json`: outputs all images to a file in JSON format
- `skopeo`: outputs all images to a file in YAML format
  to be used as source file with the `skopeo sync` command.
- `stdout`: prints all images to standard output
- `yaml`: outputs all images to a file in YAML format

```shell
helm-mirror inspect-images /tmp/helm --output stdout
helm-mirror inspect-images /tmp/helm -o stdout
helm-mirror inspect-images /tmp/helm -o file=filename
helm-mirror inspect-images /tmp/helm -o json=filename.json
helm-mirror inspect-images /tmp/helm -o yaml=filename.yaml
helm-mirror inspect-images /tmp/helm -o skopeo=filename.yaml
```

#### Global Flags

```
  -v, --verbose         verbose output
```

### version

Displays the current version of mirror.

## Install

```shell
go install github.com/kplachkov/helm-mirror@latest
```

## Test

### Running tests

To run test on this package simply run:

```shell
make test
```

#### Testing with Docker

```shell
make test.unit
```

## Building

Be sure you have all prerequisites, then build the binary by simply running:

```shell
make mirror
```

The binary will be stored under `bin` folder.
