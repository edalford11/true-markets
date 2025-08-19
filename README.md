# True Markets

## Install Go
If you haven't already, download and install the [Go Programming Language](https://go.dev/).

## Local Development
Run the API server locally with the following command:
```shell
go run cmd/main.go api
```

## Run the test suite
```shell
go test ./... -race
```

## Example curl commands
```shell
curl -XGET 'http://localhost:8080/v1/prices'
curl -XGET 'http://localhost:8080/v1/price?symbols=BTCUSDT,ETHUSDT'
```