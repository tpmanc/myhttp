![Build](https://github.com/tpmanc/myhttp/actions/workflows/go.yml/badge.svg)
![Go Report](https://goreportcard.com/badge/github.com/tpmanc/myhttp)
![Repository Top Language](https://img.shields.io/github/languages/top/tpmanc/myhttp)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tpmanc/myhttp)
![GitHub last commit](https://img.shields.io/github/last-commit/tpmanc/myhttp)

# MyHTTP
This tool makes HTTP requests to provided URLs and prints result in following format:
```
URL_1 hash_string_1
URL_2 hash_string_2
```
The tool performs the requests in parallel. By default, number of parallel requests equals to 10, but you can change it
with flag `-parallel [int]`. This number should be grater than 0. 

#### Error handling
If some invalid url provided, they will be printed to logs with error message.

## Install
```
git clone git@github.com:tpmanc/myhttp.git
```

## Requirements
* Golang v1.20+, [install](https://golang.org/doc/install)

## Build
You can build the tool with command `make build`

## Usage
```
./myhttp [option] [url_1 url_2 ...]
```

Options:
 * `-parallel [int]` Set number of parallel requests, should be grater than 0. By default, equals to 10.
 * `-help` Display help message.

Examples:
```
./myhttp http://facebook.com http://google.com
./myhttp -parallel 5 http://facebook.com http://google.com
./myhttp -help
```

## Tests
Use command `make test` to run unit tests

## Result examples
* Run with default `-parallel` value:
![run](./images/run.png)

* Run with `-parallel 1` value:
![parallel-1](./images/parallel-1.png)

* Run with invalid URL:
![parallel-1](./images/invalid-url.png)

* Run help command:
![help](./images/help.png)