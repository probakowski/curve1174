# curve1174
curve1174 is a Go client library for accessing the [Airly API](https://airly.org/en/pricing/airly-api/)

[![Build](https://github.com/probakowski/curve1174/actions/workflows/build.yml/badge.svg)](https://github.com/probakowski/curve1174/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/probakowski/curve1174)](https://goreportcard.com/report/github.com/probakowski/curve1174)

## Installation
curve1174 is compatible with modern Go releases in module mode, with Go installed:

```bash
go get github.com/probakowski/curve1174
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/probakowski/curve1174"
```

and run `go get` without parameters.

Finally, to use the top-of-trunk version of this repo, use the following command:

```bash
go get github.com/probakowski/curve1174@master
```

## Usage ##