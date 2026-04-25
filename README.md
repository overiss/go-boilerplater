# boilerplater

Go CLI tool that generates a default service structure.

## Install

From repository root:

```bash
go install ./cmd/boilerplater
```

This will install the `boilerplater` binary into your Go bin path (`$(go env GOPATH)/bin` by default).

## Usage

```bash
boilerplater make --module github.com/acme/my-service
```

By default, `cmd/<service>/main.go` uses current directory name as `<service>`.

Optional flags:

```bash
boilerplater make --module github.com/acme/my-service --service my-service
```

- `--module` sets Go module for generated service and is used for `go mod init`
- `--service` sets `<service>` in `cmd/<service>/main.go`
- if `--service` is omitted, current directory name is used and files are generated in current directory
- if `--service` is set, tool creates a folder with the same name and generates project there

The command also runs:

- `go mod init <module>`
- `go mod tidy`
