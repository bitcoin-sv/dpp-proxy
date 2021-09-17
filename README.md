# p4

[![Release](https://img.shields.io/github/release-pre/libsv/go-p4.svg?logo=github&style=flat&v=1)](https://github.com/libsv/p4/releases)
[![Build Status](https://img.shields.io/github/workflow/status/libsv/go-payment_protocol/run-go-tests?logo=github&v=3)](https://github.com/libsv/p4/actions)
[![Report](https://goreportcard.com/badge/github.com/libsv/p4?style=flat&v=1)](https://goreportcard.com/report/github.com/libsv/p4)
[![codecov](https://codecov.io/gh/libsv/go-bt/branch/master/graph/badge.svg?v=1)](https://codecov.io/gh/libsv/p4)
[![Go](https://img.shields.io/github/go-mod/go-version/libsv/p4?v=1)](https://golang.org/)
[![Sponsor](https://img.shields.io/badge/sponsor-libsv-181717.svg?logo=github&style=flat&v=3)](https://github.com/sponsors/libsv)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat&v=3)](https://gobitcoinsv.com/#sponsor)

P4 is a basic reference implementation of a Payment Protocol Server implementing the proposed BIP-270 payment flow.

This is written in go and integrates with a wallet running the Payment Protocol PayD Interface.

## Exploring Endpoints

To explore the endpoints and functionality, run the server using `go run cmd/rest-server/main.go` and navigate to [Swagger](http://localhost:8445/swagger/index.html) 
where the endpoints and their models are described in detail.

## Configuring P4

The server has a series of environment variables that allow you to configure the behaviours and integrations of the server.
Values can also be passed at build time to provide information such as build information, region, version etc.

### Server

| Key                    | Description                                                        | Default |
|------------------------|--------------------------------------------------------------------|---------|
| SERVER_PORT            | Port which this server should use                                  | :8445   |
| SERVER_HOST            | Host name under which this server is found                         | pptcl   |
| SERVER_SWAGGER_ENABLED | If set to true we will expose an endpoint hosting the Swagger docs | true    |

### Environment / Deployment Info

| Key                 | Description                                                                | Default          |
|---------------------|----------------------------------------------------------------------------|------------------|
| ENV_ENVIRONMENT     | What enviornment we are running in, for example 'production'               | dev              |
| ENV_REGION          | Region we are running in, for example 'eu-west-1'                          | local            |
| ENV_COMMIT          | Commit hash for the current build                                          | test             |
| ENV_VERSION         | Semver tag for the current build, for example v1.0.0                       | v0.0.0           |
| ENV_BUILDDATE       | Date the code was build                                                    | Current UTC time |
| ENV_BITCOIN_NETWORK | What bitcoin network we are connecting to (mainnet, testnet, stn, regtest) | regtest          |

### Logging

| Key       | Description                                                           | Default |
|-----------|-----------------------------------------------------------------------|---------|
| LOG_LEVEL | Level of logging we want within the server (debug, error, warn, info) | info    |

### PayD Wallet

| Key         | Description                                              | Default |
|-------------|----------------------------------------------------------|---------|
| PAYD_HOST   | Host for the wallet we are connecting to                 | payd    |
| PAYD_PORT   | Port the PayD wallet is listening on                     | :8443   |
| PAYD_SECURE | If true the P4 server will validate the wallet TLS certs | false   |

