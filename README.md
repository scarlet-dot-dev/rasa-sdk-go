# Rasa SDK - Golang (Work In Progress)

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/go.scarlet.dev/rasa)
[![License](https://img.shields.io/github/license/scarlet-dot-dev/rasa-sdk-go?style=flat-square)](https://www.mozilla.org/en-US/MPL/2.0/)

This package provides an SDK for Rasa chatbots, written in Go.

## Rasa SDK Features

This package implements an SDK based on the specifications for `Rasa 2.0.*`.

**Features:**

* Custom action handlers at `/webhook`.
* Custom NLG endpoint at `/nlg`. _(TODO)_
* Supports additional `/`, `/actions`, and `/health` endpoints.
* Exposes an API similar to the python SDK.
* Configurable logging and server settings.
* Code generation utility `rasagen` for boilerplate and constants, based on
  Rasa's `domain.yaml`. _(TODO - current version is outdated)_
* Clients for the `Rest` and `Callback` webhooks.
* `Callback` output channel support.

**Notes:**

* Support for `kwargs` is _experimental_.

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
