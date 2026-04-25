# boilerplater

`boilerplater` is a Go CLI tool for quickly bootstrapping a service with a
consistent project layout and стартовыми шаблонами файлов.

Tool helps to:

- create a standard folder structure for a new service
- generate basic `main/app/config/server` templates
- initialize Go module automatically (`go mod init` + `go mod tidy`)

## Install

Install latest version:

```bash
go install github.com/overiss/go-boilerplater@latest
```

This will install the `go-boilerplater` binary into your Go bin path (`$(go env GOPATH)/bin` by default).

Install exact version:

```bash
go install github.com/overiss/go-boilerplater@v0.1.1
```

## Usage

```bash
go-boilerplater make --module github.com/acme/my-service
```

By default:

- service name is taken from current directory name
- generated files are created in current directory

Optional flags:

```bash
go-boilerplater make --module github.com/acme/my-service --service my-service
```

- `--module` sets target Go module for generated service
- `--service` sets `<service>` for `cmd/<service>/main.go`
- if `--service` is set, tool creates a same-name folder and generates project inside it

The command also runs:

- `go mod init <module>`
- `go mod tidy`

## Generated structure

Tool generates the following base layout:

```text
cmd/<service>/main.go
internal/
  app/
  behavior/
  config/
  model/
    request.go
    response.go
    dto/
    dao/
  provider/
  repository/
    domain/
    integrations/
  server/
    container.go
    http/
      server.go
      handler/
      middleware/
  service/
  vars/
pkg/
  utils/
deploy/
docs/
```

## Important note about logging

Generated templates intentionally use standard Go `log` package as a neutral
default. After bootstrap, replace it with your project logger (zap/logrus/custom
wrapper/etc.) in all key entry points (`main`, `app`, `server`, `config`).
