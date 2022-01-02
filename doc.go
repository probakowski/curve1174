/*
Package curve1174 implements operations on Curve1174. It's Edwards curve with equation `x^2+y^2 = 1-1174x^2y^2`
over finite field `Fp, p=2^251-9`. It was introduced by [Bernstein, Hamburg, Krasnova, and Lange](https://eprint.iacr.org/2013/325) in 2013

Each point on curve is represented by `curve1174.Point` object. Base point is provided in `curve1174.Base`,
identity element of the curve (`x=0, y=1`) is `curve1174.E`.

API is similar to `math/big` package. The receiver denotes result and the method arguments are operation's operands.
For instance, given three `*Point` values a,b and c, the invocation

	c.Add(a,b)

computes the sum a + b and stores the result in c, overwriting whatever value was held in c before.
Operations permit aliasing of parameters, so it is perfectly ok to write

	sum.Add(sum, x)

to accumulate values x in a sum.

(By always passing in a result value via the receiver, memory use can be much better controlled. Instead of having to
allocate new memory for each result, an operation can reuse the space allocated for the result value, and overwrite
that value with the new result in the process.)

Methods usually return the incoming receiver as well, to enable simple call chaining.

Operations on curve return point in extended coordinates. To get simple x/y value they have to be converted to affine
coordinates with `(*Point).ToAffine` method. This call is expensive so be sure to avoid it for
intermediate values if possible.

Finally, `*Point` and `*FieldElement` satisfy fmt package's Formatter interface for formatted printing.
*/
package curve1174
