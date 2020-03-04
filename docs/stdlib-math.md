# Module - "math"

```golang
math := import("math")
```

## Constants

- `e`
- `pi`
- `phi`
- `sqrt2`
- `sqrtE`
- `sprtPi`
- `sqrtPhi`
- `ln2`
- `log2E`
- `ln10`
- `ln10E`

## Functions

- `abs(x float) => float`: returns the absolute value of x.
- `acos(x float) => float`: returns the arccosine, in radians, of x.
- `acosh(x float) => float`: returns the inverse hyperbolic cosine of x.
- `asin(x float) => float`: returns the arcsine, in radians, of x.
- `asinh(x float) => float`: returns the inverse hyperbolic sine of x.
- `atan(x float) => float`: returns the arctangent, in radians, of x.
- `atan2(y float, xfloat) => float`: returns the arc tangent of y/x, using the
  signs of the two to determine the quadrant of the return value.
- `atanh(x float) => float`: returns the inverse hyperbolic tangent of x.
- `cbrt(x float) => float`: returns the cube root of x.
- `ceil(x float) => float`: returns the least integer value greater than or
  equal to x.
- `copysign(x float, y float) => float`: returns a value with the magnitude of
  x and the sign of y.
- `cos(x float) => float`: returns the cosine of the radian argument x.
- `cosh(x float) => float`: returns the hyperbolic cosine of x.
- `dim(x float, y float) => float`: returns the maximum of x-y or 0.
- `erf(x float) => float`: returns the error function of x.
- `erfc(x float) => float`: returns the complementary error function of x.
- `exp(x float) => float`: returns e**x, the base-e exponential of x.
- `exp2(x float) => float`: returns 2**x, the base-2 exponential of x.
- `expm1(x float) => float`: returns e**x - 1, the base-e exponential of x
  minus 1. It is more accurate than Exp(x) - 1 when x is near zero.
- `floor(x float) => float`: returns the greatest integer value less than or
  equal to x.
- `gamma(x float) => float`: returns the Gamma function of x.
- `hypot(p float, q float) => float`: returns `Sqrt(p * p + q * q)`, taking care
  to avoid unnecessary overflow and underflow.
- `ilogb(x float) => float`: returns the binary exponent of x as an integer.
- `inf(sign int) => float`: returns positive infinity if sign >= 0, negative
  infinity if sign < 0.
- `is_inf(f float, sign int) => float`: reports whether f is an infinity,
  according to sign. If sign > 0, IsInf reports whether f is positive infinity.
  If sign < 0, IsInf reports whether f is negative infinity. If sign == 0,
  IsInf reports whether f is either infinity.
- `is_nan(f float) => float`: reports whether f is an IEEE 754 ``not-a-number''
  value.
- `j0(x float) => float`: returns the order-zero Bessel function of the first
  kind.
- `j1(x float) => float`: returns the order-one Bessel function of the first
  kind.
- `jn(n int, x float) => float`: returns the order-n Bessel function of the
  first kind.
- `ldexp(frac float, exp int) => float`: is the inverse of frexp. It returns
  frac × 2**exp.
- `log(x float) => float`: returns the natural logarithm of x.
- `log10(x float) => float`: returns the decimal logarithm of x.
- `log1p(x float) => float`: returns the natural logarithm of 1 plus its
  argument x. It is more accurate than Log(1 + x) when x is near zero.
- `log2(x float) => float`: returns the binary logarithm of x.
- `logb(x float) => float`: returns the binary exponent of x.
- `max(x float, y float) => float`: returns the larger of x or y.
- `min(x float, y float) => float`: returns the smaller of x or y.
- `mod(x float, y float) => float`: returns the floating-point remainder of x/y.
- `nan() => float`: returns an IEEE 754 ``not-a-number'' value.
- `nextafter(x float, y float) => float`: returns the next representable
  float64 value after x towards y.
- `pow(x float, y float) => float`: returns x**y, the base-x exponential of y.
- `pow10(n int) => float`: returns 10**n, the base-10 exponential of n.
- `remainder(x float, y float) => float`: returns the IEEE 754 floating-point
  remainder of x/y.
- `signbit(x float) => float`: returns true if x is negative or negative zero.
- `sin(x float) => float`: returns the sine of the radian argument x.
- `sinh(x float) => float`: returns the hyperbolic sine of x.
- `sqrt(x float) => float`: returns the square root of x.
- `tan(x float) => float`: returns the tangent of the radian argument x.
- `tanh(x float) => float`: returns the hyperbolic tangent of x.
- `trunc(x float) => float`: returns the integer value of x.
- `y0(x float) => float`: returns the order-zero Bessel function of the second
  kind.
- `y1(x float) => float`: returns the order-one Bessel function of the second
  kind.
- `yn(n int, x float) => float`: returns the order-n Bessel function of the
  second kind.
