# GitHub(QL) Exporter [![Build Status](https://drone.github.matthiasloibl.com/api/badges/metalmatze/githubql_exporter/status.svg)](https://drone.github.matthiasloibl.com/metalmatze/githubql_exporter)

[![Docker Pulls](https://img.shields.io/docker/pulls/metalmatze/githubql_exporter.svg?maxAge=604800)](https://hub.docker.com/r/metalmatze/githubql_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/metalmatze/githubql_exporter)](https://goreportcard.com/report/github.com/metalmatze/githubql_exporter)

Prometheus exporter for various metrics about your [GitHub](https://github.com/) repositories, written in Go.

### Development

You obviously should get the code

```bash
go get -u github.com/metalmatze/githubql_exporter
```

This should already put a binary called `githubql_exporter` into `$GOPATH/bin`.

Make sure you copy the `.env.example` to `.env` and change this one to your preferences.

Now during development I always run:

```bash
make install && githubql_exporter
```

Use `make install` which uses `go install` in the background to build faster during development.
