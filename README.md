# Rasa SDK - Golang

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=for-the-badge)](https://pkg.go.dev/go.scarlet.dev/rasa)
[![License](https://img.shields.io/github/license/scarlet-ai/rasa-sdk-go?style=for-the-badge)](https://https://www.mozilla.org/en-US/MPL/2.0/)

This package provides an SDK for Rasa chatbots written in Go.

## Rasa Features

This package implements the SDK based on the specifications for `Rasa 1.9.*`.

**Features:**

* Supports implementing custom action handlers at `/webhook`.
* Supports implementing a custom NLG endpoint at `/nlg`. (TODO)
* Supports the additional `/`, `/actions`, and `/health` endpoints.
* Exposes an API is similar to the python SDK.
* Configurable logging and server settings.
* Code generation CLI for boilerplate and constants, based on Rasa's
  `domain.yaml`.  

## Import

```bash
go get go.scarlet.dev/rasa
```

## Godoc

[https://pkg.go.dev/go.scarlet.dev/rasa](https://pkg.go.dev/go.scarlet.dev/rasa).

## License

Unless otherwise specified, code present in this library is licensed under the
[Mozilla Public License Version v2.0](https://www.mozilla.org/en-US/MPL/2.0/ "MPL v2.0").

## Authors

* Eddy (@scarlet.dev)
