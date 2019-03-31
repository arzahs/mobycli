# Simple Moby (docker) CLI

## Requirements

 - Docker daemon with Docker API > v1.18
 - Go 1.10+
 - Make
 - Go dep

## Install
```bash
go get -u github.com/golang/dep/cmd/dep # installing dep
make dep
make build
```

## Usage example

```sh
./mobycli run postgres
./mobycli ps
./mobycli stop $CONTAINERID

```