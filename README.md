# cciu - Check Container Images for Updates

[![GitHub Release](https://img.shields.io/github/v/release/mgumz/cciu.svg)](https://github.com/mgumz/cciu/releases/latest)
[![Status](https://github.com/mgumz/cciu/actions/workflows/actions.yaml/badge.svg)](https://github.com/mgumz/cciu/actions/workflows/actions.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mgumz/cciu)](https://goreportcard.com/report/github.com/mgumz/cciu)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mgumz/cciu.svg)](https://github.com/mgumz/cciu)

`cciu` checks container repositories for available updates for a given set of
container image names.

## Usage

    $> cciu [flags] <image:1> [<image:2> ...]

### Flags

    -exclude-beta-tags        - exclude 'beta' tags (and 'alpha', 'rc')
    -h                        - show help
    -json                     - print JSON
    -json-pretty              - print JSON, prettyfied
    -limit-per-registry       - n concurrent fetch operations per registry
    -show-old                 - show older tags
    -simple-markers           - use simple ascii markers
    -skip-non-semver          - skip non-semver tags
    -stats                    - show stats
    -strict-labels            - strict label matching
    -timeout                  - time out fetch operation after <dur>
    -keep ["major"|"minor"]   - keep major/minor version
    -version                  - show version

## Snippets

Check for some abitrary container image update:

    $> cciu alpine:3.11
    alpine:3.11
    ▲       alpine:3.13.5

How to read the output: the image to be checked: "alpine:3.11". An update for
this container image exists in the alpine repository: "alpine:3.13.5" is
available.

If no update is available, the output would look like this:

    $> cciu alpine:3.13.5
    =       alpine:3.13.5

To only check for updates of patch version keeping the minor version:

    $> cciu -keep minor alpine:3.11.4
    ▲       alpine:3.11.13

In addition, the output could be JSON to process it somewhere else:

    $> cciu -json-pretty alpine:3.11
    {
      "images": [
        {
          "requested": "alpine:3.11",
          "tags": [
            {
              "name": "alpine:3.13.5",
              "version": "3.13.5",
              "verdict": "ahead"
            }
          ],
          "verdict": "outdated"
        }
      ]
    }

Check container images for upstream updates:

    $> docker images --format '{{.Repository}}:{{.Tag}}' | tee /tmp/images.txt
    alpine:3.13
    traefik:v2.4.7
    restic/restic:0.12.0

    $> xargs cciu < /tmp/images.txt

Check container image updates for running containers:

    $> docker ps --format '{{.Image}} | xargs cciu

Check container image updates for running containers on remote machine:

    $> ssh example.com -- docker ps --format \"{{.Image}}\" | xargs cciu

Check container image updates for running contains within a k8s cluster:

    $> kubectl get pods -o json -n example-ns | \
        jq -r '.items[]|. as $item| [(.spec.containers[]| [ \
            .image,"@",\
            (.=$item|.metadata.namespace),"--",\
            (.=$item|.metadata.name),"--",\
            .name]|join("") )]|join("\n")' | tee /tmp/images.txt

    $> cat /tmp/images.txt
    freeradius/freeradius-server:3.0.21-alpine@example-ns--sample1-deployment-6d99c6fd44-q5x82--radius-proxy
    freeradius/freeradius-server:3.0.21-alpine@example-ns--sample2-deployment-7979486848-f29gq--radius-proxy
    influxdb:1.8.3-alpine@example-ns--influxdb-0--influxdb

    $> xargs cciu < /tmp/images.txt
    <output>

Here the name of the Pod and the k8s namespace are added to the "context" part
of the image name. This helps to identify the deployment which might benefit
from an upgrade of the container image.

## Installation

    $> go install -v github.com/mgumz/cciu/cmd/cciu@latest

## Building

    $> make cciu

## Roadmap

### Credential handling

Status: **NOT** available atm to check private registries.

Needed: yes.

### Handling Image References by digest

As an alternative to specify container images with a "tag" (like
"alpine:latest"), one can also specify an image via a digest like:

    alpine@sha256:8d99168167baa6a6a0d7851b9684625df9c1455116a9601835c2127df2aaa2f5

These container images have a release date as well. To compare dates is easy
so one could collect "newer" container image names.

Needed: Optional.

## Related projects

* [skopeo](https://github.com/containers/skopeo) is a set of CLI tools to
  directly work with container registries and images. `skopeo-list-tags`
  is able to list images/tags.


## License

See LICENSE.
