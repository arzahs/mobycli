# Simple Moby (docker) CLI

## Requirements

 - DockerAPI > v1.12
 - Go 1.10+
 - Go dep

## Install
```sh
$ go get -u github.com/golang/dep/...
$ make dep
$ make build
```

## Usage example

```sh
$ ./mobycli run postgres
$ ./mobycli ps
$ ./mobycli stop $CONTAINERID

```