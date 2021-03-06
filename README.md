# ALoLStats

[![Build status](https://git.abyle.org/hps/alolstats/badges/master/pipeline.svg)](https://git.abyle.org/hps/alolstats/commits/master)
[![Coverage Status](https://git.abyle.org/hps/alolstats/badges/master/coverage.svg)](https://git.abyle.org/hps/alolstats/commits/master)
[![Docker](https://img.shields.io/docker/pulls/hpsch/alolstats.svg)](https://hub.docker.com/r/hpsch/alolstats/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

## Short Description

ALoLstats, a League of Legends Statistics aggregation and calculation server

## Abstract

For players which want to be competitive or which are still learning it is essential to study various game metrics. Usually players use one of the various
websites which exist solely for the purpose to present statistics about League of Legends and which give you informations like the lanes a
certain champion should be played in or its win/loss and ban rates.

Most of those websites do not describe their methods how those informations were obtained, how the statistics is calculated and on which dataset
the calculations are based on. In every scientific setting this would be a strict no-go and if bad underlying data is used, or wrong assumptions on the data
are made, this can lead to misinterpretation of the results which may lead to wrong decision making.

In this project we aim to improve on that by implementing a completely open source (MIT license) solution to the calculation and preparation of statistics for
various aspects of the game.

Currently the focus of development lies in the calculation of Champion statistics, like lane/role, average KDA, win/loss ratios, ban ratios and various other aspects
related to Champions. In the near future Summoner data shall also be considered. The decisions on which features are tackled next are based on suggestions and requests by friends and colleagues and personal experiences while playing League of Legends.

The implementation is done in Go with the possibility to use R for complex statistical calculations. The data storage is interchangeable and currently
a MongoDB backend is supported.

The Riot API client is implemented from scratch for this project, including region support and rate limit handling, and is currently based on V4 of the Riot API.

In combination with this backend a frontend called [ALoLstats-web](https://git.abyle.org/hps/alolstats-web) is also beeing developed.

## How to run it

Probably the easiest way to try out ALoLstats is using Docker. Currently no official images are provided in a public registry, but it can be easily built using

```
make build-container-gitcommit
```

and then start it using

```
docker run --name ALoLstats -v /path/to/config/file.toml:/app/config/config.toml:ro hpsch/alolstats:latest
```

where _/path/to/config/file.toml_ has to be replaced with the path to your config file.

An example is provided in the cfg/ directory and it is enough to insert your API key to use this config.

To expose the http port for the REST API use
```
-p 8000:8000
```
where 8000 should be exchanged with the port set in the config file.

## API Reference (usually out of date and highly in flux)

The following endpoints are currently available. A detailed description will be provided at a later point when the API becomes more stable.

### Match related endpoints

* **/v1/match?id=matchID**: Returns informations of the match with id=matchID (e.g., 2585564744)

### Champion related endpoints

* **/v1/champions**: Returns a list of all currently available champions
* **/v1/champion-rotations**: Returns a list of the current free champion rotation

### Summoner related endpoints

* **/v1/summoner/byname?name=summonerName**: Returns information about a summoner specified by name=summonerName
* **/v1/summoner/bysummonerid?id=summonerID**: Returns information about a summoner specified by name=summonerID
* **/v1/summoner/byaccountid?id=accountID**: Returns information about a summoner specified by name=accountID

### Statistics related endpoints

* **/v1/stats/overview**: Temporary page which lists all available plots related to Champion statistics
* **/v1/stats/champion/byid?id=championId&gameversion=exactGameVersion**: Returns stats for the Champion with id=championId and the specified game version (e.g., 110 and 8.24)
* **/v1/stats/champion/byname?name=championName&gameversion=exactGameVersion**: Returns stats for the Champion with name=championName and the specified game version (e.g., Sivir and 8.24)

### ALoLStats related endpoints

* **/v1/storage/summary**: Returns information of the stored data in the storage or its backend
