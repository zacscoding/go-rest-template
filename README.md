# Go Rest API Base Template

Golang REST API Template/Boilerplate with Gin/Gorm/Fx.

# Run with docker

```shell
$ make compose.local.up
```

Check .http files for tests in [sample.http](./scripts/http/sample.http)

# Tests

```shell
# Clean test cache, Run go tests and build
$ make tests
```

# Lint

```shell
# Run golangci-lint
$ make lint
```

# Build

```shell
$ make build
$ tree ./build/bin                                       
./build/bin
└── apiserver
```

# Release

```shell
$ make release
# $ make release.local  # for local tests
# $ make release.dry    # for dry run
```

# Relocate project
; TBD