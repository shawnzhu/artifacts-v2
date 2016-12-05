[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/shawnzhu/artifacts-v2.svg?branch=master)](https://travis-ci.org/shawnzhu/artifacts-v2)
[![Coverage Status](https://coveralls.io/repos/github/shawnzhu/artifacts-v2/badge.svg?branch=master)](https://coveralls.io/github/shawnzhu/artifacts-v2?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/shawnzhu/artifacts-v2)](https://goreportcard.com/report/github.com/shawnzhu/artifacts-v2)

# Artifacts

work with artifacts over HTTP API.

## Requirements

* Go v1.7.3
* A running PostgreSQL v9
* credentials for saving object into AWS S3
* a piece signing key for JWT token verification

## Install

```
$ go get ./...
```

### Database setup

1. Create a database named `test_artifacts`
1. Run SQL from [store/ddl/1.sql](store/ddl/1.sql)

### AWS credentials

See [Configuring Credentials](https://github.com/aws/aws-sdk-go#configuring-credentials)

## Testing

```
$ go test -v ./...
```

Starting test server:

```
$ go run server.go
```
