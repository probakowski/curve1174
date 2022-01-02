# curve1174
curve1174 is a Go library implementing operations on Curve1174. It's Edwards curve with equation `x^2+y^2 = 1-1174x^2y^2`
over finite field `Fp, p=2^251-9`. It was introduced by [Bernstein, Hamburg, Krasnova, and Lange](https://eprint.iacr.org/2013/325) in 2013

[![Build](https://github.com/probakowski/curve1174/actions/workflows/go.yml/badge.svg)](https://github.com/probakowski/curve1174/actions/workflows/go.yml)
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
Each point on curve is represented by `curve1174.Point` object. Base point is provided in `curve1174.Base`, 
identity element of the curve (`x=0, y=1`) is `curve1174.E`.

API is similar to `math/big` package. The receiver denotes result and the method arguments are operation's operands.
For instance, given three `*Point` values a,b and c, the invocation

```go
c.Add(a,b)
```

computes the sum a + b and stores the result in c, overwriting whatever value was held in c before.
Operations permit aliasing of parameters, so it is perfectly ok to write
```go
sum.Add(sum, x)
```
to accumulate values x in a sum.

(By always passing in a result value via the receiver, memory use can be much better controlled. Instead of having to 
allocate new memory for each result, an operation can reuse the space allocated for the result value, and overwrite 
that value with the new result in the process.)

Methods usually return the incoming receiver as well, to enable simple call chaining.

Operations on curve return point in extended coordinates. To get simple x/y value they have to be converted to affine 
coordinates with `(*Point).ToAffine` method. This call is expensive so be sure to avoid it for 
intermediate values if possible.

Finally, `*FieldElement` satisfy fmt package's Formatter interface for formatted printing