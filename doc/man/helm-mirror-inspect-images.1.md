% helm-mirror-inspect-images(1) # helm-mirror inspect-images - Extract all the container images listed in each chart.
# NAME
helm-mirror inspect-images - Extract all the container images listed in each chart.

# SYNOPSIS
**helm-mirror inspect-images** target
[**--help**|**-h**]

# DESCRIPTION
**helm-mirror inspect-images** Extract all the container images listed in each Helm Chart or
the Helm Charts in the folder provided. This command dumps the images on
**stdout** by default.

**helm-mirror inspect-images** Has different type of outputs for the images to make
it easier to interact with the sub-command, for more options check **output**
option.

# GLOBAL OPTIONS

**-v, --verbose**
  Verbose output

# OPTIONS

**-h, --help**
  Print usage statement.

**-i, --ignore-errors**
  Ignores errors while downloading or processing charts.

**-o, --output**
  choose an output for the list of images and specify the file name, if not specified 'images.out' will be the default.
  (file|json|skopeo|**stdout**|yaml)

# EXAMPLES
The following examples show different ways to interact with **mirror inspect-images**
command.

Inspect a folder and print to **stdout** (default for the **--output** option)
```
% helm-mirror inspect-images /tmp/helm
```

Inspect a chart file and print to **stdout**
```
% helm-mirror inspect-images /tmp/helm/chart.tgz
```

Inspect a folder and export to other formats.
```
% helm-mirror inspect-images /tmp/helm -o file=images.txt
% helm-mirror inspect-images /tmp/helm -o json=images.json
% helm-mirror inspect-images /tmp/helm -o yaml=images.yaml
```

Inspect a folder and ignore the errors while rendering the chart, this
errors are usually for missing required values in the charts.
```
% helm-mirror inspect-images /tmp/helm --ignore-errors
```

# SEE ALSO
**helm-mirror**(1),
**helm-mirror-help**(1),
**helm-mirror-version**(1)
