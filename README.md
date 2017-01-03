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
* a public key for verifying JWT token signed by a private key

## Install

```
$ make install
```

### Database setup

1. Create a database named `test_artifacts`
1. Run SQL from [store/ddl/1.sql](store/ddl/1.sql)

### AWS credentials

See [Configuring Credentials](https://github.com/aws/aws-sdk-go#configuring-credentials)

### Public key

Make sure the environment variable `JWT_PUBLIC_KEY` contains public key in PEM format.
E.g., `export JWT_PUBLIC_KEY="$(cat public_key.pem)"`

## Testing

```
$ make test
```

Starting test server:

```
$ go run cmd/travis-artifacts/main.go
```

## Build binary

```
make build
```

## Release

Push to docker hub:

```
make TAG=<tag-name> release
```

## Deploy to Kubernetes

create distributed app on ready Kubernetes cluster:

```
$ kubectl create -f k8s-app.yml
```
