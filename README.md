# certrenewal
[![Go Report Card](https://goreportcard.com/badge/github.com/jduepmeier/certrenewal)](https://goreportcard.com/report/github.com/jduepmeier/certrenewal)
[![Go Reference](https://pkg.go.dev/badge/github.com/jduepmeier/certrenewal.svg)](https://pkg.go.dev/github.com/jduepmeier/certrenewal)

Tool to renew certificates from [hashicorp vault](https://www.vaultproject.io/).

## Build

```bash
make
```

## Usage

Copy the example config [config.sample.yaml](config.sample.yaml) to `config.yaml` and configure it.

Then run the tool:
```bash
./bin/cert-renewal --config config.yaml
```

Run with `--help` for all command line arguments.

## Configuration

See [config.sample.yaml](config.sample.yaml) for all valid config options.

## CHANGELOG

See [CHANGELOG.md](CHANGELOG.md).