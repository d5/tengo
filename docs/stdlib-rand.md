# Module - "rand"

```golang
rand := import("rand")
```

## Functions

- `seed(seed int)`: uses the provided seed value to initialize the default
  Source to a deterministic state.
- `exp_float() => float`:  returns an exponentially distributed float64 in the
  range (0, +math.MaxFloat64] with an exponential distribution whose rate
  parameter (lambda) is 1 and whose mean is 1/lambda (1) from the default
  Source.
- `float() => float`: returns, as a float64, a pseudo-random number in
  [0.0,1.0) from the default Source.
- `int() => int`: returns a non-negative pseudo-random 63-bit integer as an
  int64 from the default Source.
- `intn(n int) => int`: returns, as an int64, a non-negative pseudo-random
  number in [0,n) from the default Source. It panics if n <= 0.
- `norm_float) => float`: returns a normally distributed float64 in the range
  [-math.MaxFloat64, +math.MaxFloat64] with standard normal distribution
  (mean = 0, stddev = 1) from the default Source.
- `perm(n int) => [int]`: returns, as a slice of n ints, a pseudo-random
  permutation of the integers [0,n) from the default Source.
- `read(p bytes) => int/error`: generates len(p) random bytes from the default
  Source and writes them into p. It always returns len(p) and a nil error.
- `rand(src_seed int) => Rand`: returns a new Rand that uses random values from
  src to generate other random values.

## Rand

- `seed(seed int)`: uses the provided seed value to initialize the default
  Source to a deterministic state.
- `exp_float() => float`:  returns an exponentially distributed float64 in the
  range (0, +math.MaxFloat64] with an exponential distribution whose rate
  parameter (lambda) is 1 and whose mean is 1/lambda (1) from the default Source.
- `float() => float`: returns, as a float64, a pseudo-random number in
  [0.0,1.0) from the default Source.
- `int() => int`: returns a non-negative pseudo-random 63-bit integer as an
  int64 from the default Source.
- `intn(n int) => int`: returns, as an int64, a non-negative pseudo-random
  number in [0,n) from the default Source. It panics if n <= 0.
- `norm_float) => float`: returns a normally distributed float64 in the range
  [-math.MaxFloat64, +math.MaxFloat64] with standard normal distribution
  (mean = 0, stddev = 1) from the default Source.
- `perm(n int) => [int]`: returns, as a slice of n ints, a pseudo-random
  permutation of the integers [0,n) from the default Source.
- `read(p bytes) => int/error`: generates len(p) random bytes from the default
  Source and writes them into p. It always returns len(p) and a nil error.
