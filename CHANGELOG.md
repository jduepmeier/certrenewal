# CHANGELOG

## 1.3.5 (2025-04-09)

### Fix

- **deps**: update module github.com/hashicorp/vault/api to v1.16.0
- **deps**: update module golang.org/x/crypto to v0.37.0

## 1.3.4 (2025-02-09)

### Fix

- **deps**: update module golang.org/x/crypto to v0.33.0

## 1.3.3 (2024-12-31)

### Fix

- **deps**: update module golang.org/x/crypto to v0.31.0 [security]
- **deps**: update module github.com/hashicorp/vault/api to v1.15.0
- **deps**: update module github.com/stretchr/testify to v1.10.0

## 1.3.2 (2024-06-29)

### Fix

- **deps**: update module github.com/hashicorp/vault/api to v1.14.0
- **deps**: update module golang.org/x/crypto to v0.24.0
- **deps**: update module github.com/jessevdk/go-flags to v1.6.1

## 1.3.1 (2024-04-03)

### Fix

- **deps**: update module github.com/hashicorp/vault/api to v1.12.2
- **deps**: update module golang.org/x/crypto to v0.21.0
- **deps**: update module github.com/stretchr/testify to v1.9.0
- **deps**: update module golang.org/x/crypto to v0.20.0
- **deps**: update module golang.org/x/crypto to v0.19.0
- **deps**: update module github.com/hashicorp/vault/api to v1.12.0
- **deps**: update module github.com/hashicorp/vault/api to v1.11.0
- **deps**: update module golang.org/x/crypto to v0.18.0
- **deps**: update module golang.org/x/crypto to v0.17.0 [security]
- **deps**: update module golang.org/x/crypto to v0.16.0

## 1.3.0 (2023-11-24)

### Feat

- **cert**: allow to overwrite the pki_path for certs

## 1.2.3 (2023-11-10)

### Fix

- **deps**: update module golang.org/x/crypto to v0.15.0
- **deps**: update module golang.org/x/crypto to v0.14.0
- **deps**: update module github.com/hashicorp/vault/api to v1.10.0
- **deps**: update module golang.org/x/crypto to v0.13.0

## 1.2.2 (2023-08-27)

### Fix

- **deps**: update module golang.org/x/crypto to v0.12.0
- **deps**: update module golang.org/x/crypto to v0.11.0
- **deps**: update module golang.org/x/crypto to v0.10.0
- **deps**: update module github.com/stretchr/testify to v1.8.4
- **deps**: update module github.com/sirupsen/logrus to v1.9.3
- **deps**: update module github.com/hashicorp/vault/api to v1.9.2
- **deps**: update module github.com/stretchr/testify to v1.8.3

## 1.2.1 (2023-05-18)

### Fix

- **deps**: update module golang.org/x/crypto to v0.9.0
- **deps**: update module github.com/sirupsen/logrus to v1.9.2

## 1.2.0 (2023-05-06)

### Feat

- **cert**: check ip sans in tls cert for changes

### Fix

- **deps**: update module github.com/hashicorp/vault/api to v1.9.1
- **deps**: update module golang.org/x/crypto to v0.8.0
- **deps**: update module golang.org/x/crypto to v0.7.0
- **deps**: update module github.com/stretchr/testify to v1.8.2

## 1.1.2 (2023-02-23)

### Fix

- **deps**: use yaml.v3 instead of yaml.v2 in code
- **deps**: update module gopkg.in/yaml.v2 to v3
- **deps**: update module golang.org/x/crypto to v0.6.0

## 1.1.1 (2023-02-23)

* Bump golang.org/x/net from 0.5.0 to 0.7.0

## 1.1.0 (2023-02-23)

* build(goreleaser): remove deprecated config
* fix(deps): update module github.com/stretchr/testify to v1.8.1
* Update module github.com/sirupsen/logrus to v1.9.0
* build: use jduepmeier/renovate-config
* Update module github.com/jessevdk/go-flags to v1.5.0
* Update module github.com/hashicorp/vault/api to v1.9.0
* Bump golang.org/x/text from 0.3.7 to 0.3.8
* Add renovate.json
* Add ssh and config tests
* Update dependencies
* fix -k, --insecure does not work

## 1.0.0 (2021-04-25)

* renew certificates
* auth with approle
* set log options from command line
* run hooks after renewal
* support ssh cert renewal
* add version flag
