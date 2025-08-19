# True Markets

## Install Go
If you haven't already, download and install the [Go Programming Language](https://go.dev/).

## Configuration
```shell
cp config/config.default.yml config/config.yml
```

You may make any changes to the config in that new file, which is not commited to version control.
Symbols can be added or removed in that config file or updated via env variable like so
```shell
export BINANCE_SYMBOLS=BTCUSDT,ETHUSDT
```

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

## Assumptions
We will not be letting the user know if a symbol is not found when using the /price endpoint. It will just not be
included in the response.