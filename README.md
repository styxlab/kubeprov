# kubeprov
Kubernetes Cluster Provisioning on Hetzner Cloud

## Disclaimer 

THIS IS WORK IN PROGRESS AND CURRENTLY NOT FUNCTIONAL

## Prerequisites

- `go version` >= `go1.11.2`
- `go env |grep GOPATH` is set
- `echo 'export HCLOUD_TOKEN=<HetznerCloudToken>\n' >> ~/.bashrc`

## Install

```
$ go get -u github.com/styxlab/kubeprov
$ sudo cp "$GOPATH/bin/kubeprov" /usr/local/bin/
```

## Usage

```
$ kubeprov
```