# Operators

## Int

### Equality

- `(int) == (int) = (bool)`: equality
- `(int) != (int) = (bool)`: inequality

### Arithmetic Operators

- `(int) + (int) = (int)`: sum
- `(int) - (int) = (int)`: difference
- `(int) * (int) = (int)`: product
- `(int) / (int) = (int)`: quotient
- `(int) % (int) = (int)`: remainder
- `(int) + (float) = (float)`: sum
- `(int) - (float) = (float)`: difference
- `(int) * (float) = (float)`: product
- `(int) / (float) = (float)`: quotient
- `(int) + (char) = (char)`: sum
- `(int) - (char) = (char)`: difference

### Bitwise Operators

- `(int) & (int) = (int)`: bitwise AND
- `(int) | (int) = (int)`: bitwise OR
- `(int) ^ (int) = (int)`: bitwise XOR
- `(int) &^ (int) = (int)`: bitclear (AND NOT)
- `(int) << (int) = (int)`: left shift
- `(int) >> (int) = (int)`: right shift

### Comparison Operators

- `(int) < (int) = (bool)`: less than
- `(int) > (int) = (bool)`: greater than
- `(int) <= (int) = (bool)`: less than or equal to
- `(int) >= (int) = (bool)`: greater than or equal to
- `(int) < (float) = (bool)`: less than
- `(int) > (float) = (bool)`: greater than
- `(int) <= (float) = (bool)`: less than or equal to
- `(int) >= (float) = (bool)`: greater than or equal to
- `(int) < (char) = (bool)`: less than
- `(int) > (char) = (bool)`: greater than
- `(int) <= (char) = (bool)`: less than or equal to
- `(int) >= (char) = (bool)`: greater than or equal to

## Float

### Equality

- `(float) == (float) = (bool)`: equality
- `(float) != (float) = (bool)`: inequality

### Arithmetic Operators

- `(float) + (float) = (float)`: sum
- `(float) - (float) = (float)`: difference
- `(float) * (float) = (float)`: product
- `(float) / (float) = (float)`: quotient
- `(float) + (int) = (int)`: sum
- `(float) - (int) = (int)`: difference
- `(float) * (int) = (int)`: product
- `(float) / (int) = (int)`: quotient

### Comparison Operators

- `(float) < (float) = (bool)`: less than
- `(float) > (float) = (bool)`: greater than
- `(float) <= (float) = (bool)`: less than or equal to
- `(float) >= (float) = (bool)`: greater than or equal to
- `(float) < (int) = (bool)`: less than
- `(float) > (int) = (bool)`: greater than
- `(float) <= (int) = (bool)`: less than or equal to
- `(float) >= (int) = (bool)`: greater than or equal to

## String

### Equality

- `(string) == (string) = (bool)`: equality
- `(string) != (string) = (bool)`: inequality

### Concatenation

- `(string) + (string) = (string)`: concatenation
- `(string) + (other types) = (string)`: concatenation (after string-converted)

### Comparison Operators

- `(string) < (string) = (bool)`: less than
- `(string) > (string) = (bool)`: greater than
- `(string) <= (string) = (bool)`: less than or equal to
- `(string) >= (string) = (bool)`: greater than or equal to

## Char

### Equality

- `(char) == (char) = (bool)`: equality
- `(char) != (char) = (bool)`: inequality

### Arithmetic Operators

- `(char) + (char) = (char)`: sum
- `(char) - (char) = (char)`: difference
- `(char) + (int) = (char)`: sum
- `(char) - (int) = (char)`: difference

### Comparison Operators

- `(char) < (char) = (bool)`: less than
- `(char) > (char) = (bool)`: greater than
- `(char) <= (char) = (bool)`: less than or equal to
- `(char) >= (char) = (bool)`: greater than or equal to
- `(char) < (int) = (bool)`: less than
- `(char) > (int) = (bool)`: greater than
- `(char) <= (int) = (bool)`: less than or equal to
- `(char) >= (int) = (bool)`: greater than or equal to

## Bool

### Equality

- `(bool) == (bool) = (bool)`: equality
- `(bool) != (bool) = (bool)`: inequality

## Bytes

### Equality

Test whether two byte array contain the same data. Uses
[bytes.Compare](https://golang.org/pkg/bytes/#Compare) internally.

- `(bytes) == (bytes) = (bool)`: equality
- `(bytes) != (bytes) = (bool)`: inequality

## Time

### Equality

Tests whether two times represent the same time instance. Uses
[Time.Equal](https://golang.org/pkg/time/#Time.Equal) internally.

- `(time) == (time) = (bool)`: equality
- `(time) != (time) = (bool)`: inequality

### Arithmetic Operators

- `(time) - (time) = (int)`: difference in nanoseconds (duration)
- `(time) + (int) = (time)`: time + duration (nanoseconds)
- `(time) - (int) = (time)`: time - duration (nanoseconds)

### Comparison Operators

- `(time) < (time) = (bool)`: less than
- `(time) > (time) = (bool)`: greater than
- `(time) <= (time) = (bool)`: less than or equal to
- `(time) >= (time) = (bool)`: greater than or equal to

## Array and ImmutableArray

### Equality

Tests whether two _(immutable)_ arrays contain the same objects.

- `(array) == (array) = (bool)`: equality
- `(array) != (array) = (bool)`: inequality
- `(array) == (immutable-array) = (bool)`: equality
- `(array) != (immutable-array) = (bool)`: inequality
- `(immutable-array) == (immutable-array) = (bool)`: equality
- `(immutable-array) != (immutable-array) = (bool)`: inequality
- `(immutable-array) == (array) = (bool)`: equality
- `(immutable-array) != (array) = (bool)`: inequality

### Concatenation

- `(array) + (array)`: return a concatenated array  

## Map and ImmutableMap

### Equality

Tests whether two _(immutable)_ maps contain the same key-objects.

- `(map) == (map) = (bool)`: equality
- `(map) != (map) = (bool)`: inequality
- `(map) == (immutable-map) = (bool)`: equality
- `(map) != (immutable-map) = (bool)`: inequality
- `(immutable-map) == (immutable-map) = (bool)`: equality
- `(immutable-map) != (immutable-map) = (bool)`: inequality
- `(immutable-map) == (map) = (bool)`: equality
- `(immutable-map) != (map) = (bool)`: inequality
