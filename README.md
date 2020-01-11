# Patch Hosts Service

This project is a reimplementation of https://patchbay.pub/ ([docs](https://patchbay.pub/docs/index.html)),
with a few changes:

## Requester/Responder

Headers look like this:
```
X-Phs-<index>-<header-name>
```

If the user sends a request with the headers:
```
foobar: barfoo_1
foobar: barfoo_2
```
the resulting headers on the other side look like this:
```
X-Phs-0-Foobar: barfoo_1
X-Phs-1-Foobar: barfoo_2
```

# Setup

This project has no third party dependencies, except `github.com/oxequa/realize` for development.

## dev
start hot reloading dev server

```
make dev
```

## build
build binary and save it to `build/patchbay-server`

```
make build
```

## test coverage
create test coverage and open it in browser

```
make cov
```

## test
run unit tests

```
make test
```

## clean
delete folder `build` and `coverage.html`

```
make clean
```