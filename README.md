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

where _/path/to/config/file.toml_ has to be replaced with the path to your config file. An example is provided in the cfg/ directory and it is enough to insert your API key to use this config.

To expose the http port for the REST API use
```
-p 8000:8000
```
where 8000 should be exchanged with the port set in the config file.

## API Reference

The following endpoints are currently available. A detailed description will be provided at a later point when the API becomes more stable.

* **/v1/champions**: Returns a list of all currently available champions
* **/v1/champion-rotations**: Returns a list of the current free champion rotation
* **/v1/match?id=matchId**: Returns informations of the match with id=matchId (e.g., 2585564744)
* **/v1/matches?gameversion=exactGameVersion**: Returns all currently stored matches for the specified game version (e.g., 7.17.200.3955)
* **/v1/stats/champion?id=championId&gameversion=exactGameVersion**: Returns stats for the Champion with id=championId and the specified game version (e.g., 110 and 7.17.200.3955)
