# ALoLStats

[![Build Status](https://travis-ci.org/torlenor/alolstats.svg?branch=master)](https://travis-ci.org/torlenor/alolstats)
[![Coverage Status](https://coveralls.io/repos/github/torlenor/alolstats/badge.svg?branch=master)](https://coveralls.io/github/torlenor/alolstats?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/torlenor/alolstats)](https://goreportcard.com/report/github.com/torlenor/alolstats)
[![Docker](https://img.shields.io/docker/pulls/hpsch/alolstats.svg)](https://hub.docker.com/r/hpsch/alolstats/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

## Description

This is ALoLstats, a League of Legends Statistics aggregation and calculation server.

## How to run it

Probably the easiest way to try out ALoLstats is using Docker. To pull the latest version from DockerHub and start it just type

```
docker run --name ALoLstats -v /path/to/config/file.toml:/app/config/config.toml:ro hpsch/alolstats:latest
```

where _/path/to/config/file.toml_ has to be replaced with the path to your config file. Specify
```
-p 8000:8000
```
if you want to expose the port you specified in the config file for the REST API.
