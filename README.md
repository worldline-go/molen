# molen

[![License](https://img.shields.io/github/license/worldline-go/molen?color=red&style=flat-square)](https://raw.githubusercontent.com/worldline-go/molen/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/worldline-go_molen?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=worldline-go_molen)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/worldline-go/molen/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/worldline-go/molen/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/worldline-go/molen?style=flat-square)](https://goreportcard.com/report/github.com/worldline-go/molen)
[![Releases](https://img.shields.io/badge/download-releases-pink?style=flat-square&logo=github
)](https://github.com/worldline-go/molen/releases/latest)

Pub/sub API for message brokers.

## Usage

Configuration can give as `molen.[json|yaml|yml|toml]` in the /etc/ or near to binary or with `CONFIG_FILE` env value.  
It can also work with CONSUL and VAULT.

> For configuration check the [config.go](./internal/config/config.go) file.

### Development

Create kafka server and console use:

```sh
make env
```

Access to http://localhost:7071 for console.

Run molen with make file:

```sh
make run
```

Access to http://localhost:8080/swagger/index.html for API swagger page.
