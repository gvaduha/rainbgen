# RainbGen
toy SHA-2 hashing microservice with MongoDB as backend storage.

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/json-iterator/go/master/LICENSE)
[![Build Status](https://travis-ci.org/gvaduha/rainbgen.svg?branch=master)](https://travis-ci.org/gvaduha/rainbgen)

## Endpoints
Requests and responses uses JSON payload in body 

### Send data for hashing use-case
1. User issues POST request with fields:
 a. payload - text to hash
 b. hash_rounds_cnt - rounds of hashing f(...f(payload))
2. Service creates job put in storage and return job id immideately.
3. Service starts processing job in parallel.

### Fetching result use-case
Используя полученный id пользователь может получить результат этой задачи с
полями:
1. id - job id
2. payload - original text
3. hash_rounds_cnt - number of hashing rounds performed
4. status - processing state
 a. “in progress” - not completed
 b. “finished” - hash calculated
5. hash - hash string of payload

## Notes
* Design allows to easily change underlying database
* Service could be deployed via docker-compose and tuned with corresponding .env file
