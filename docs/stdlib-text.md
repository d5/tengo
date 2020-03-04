# Module - "text"

```golang
text := import("text")
```

## Functions

- `re_match(pattern string, text string) => bool/error`: reports whether the
  string s contains any match of the regular expression pattern.
- `re_find(pattern string, text string, count int) => [[{text: string, begin: int, end: int}]]/undefined`:
  returns an array holding all matches, each of which is an array of map object
  that contains matching text, begin and end (exclusive) index.
- `re_replace(pattern string, text string, repl string) => string/error`:
  returns a copy of src, replacing matches of the pattern with the replacement
  string repl.
- `re_split(pattern string, text string, count int) => [string]/error`: slices
  s into substrings separated by the expression and returns a slice of the
  substrings between those expression matches.
- `re_compile(pattern string) => Regexp/error`: parses a regular expression and
  returns, if successful, a Regexp object that can be used to match against
  text.
- `compare(a string, b string) => int`: returns an integer comparing two
  strings lexicographically. The result will be 0 if a==b, -1 if a < b, and +1
  if a > b.
- `contains(s string, substr string) => bool`: reports whether substr is within
  s.
- `contains_any(s string, chars string) => bool`: reports whether any Unicode
  code points in chars are within s.
- `count(s string, substr string) => int`: counts the number of non-overlapping
  instances of substr in s.
- `equal_fold(s string, t string) => bool`: reports whether s and t,
  interpreted as UTF-8 strings,
- `fields(s string) => [string]`: splits the string s around each instance of
  one or more consecutive white space characters, as defined by unicode.IsSpace,
  returning a slice of substrings of s or an empty slice if s contains only
  white space.
- `has_prefix(s string, prefix string) => bool`: tests whether the string s
  begins with prefix.
- `has_suffix(s string, suffix string) => bool`: tests whether the string s
  ends with suffix.
- `index(s string, substr string) => int`: returns the index of the first
  instance of substr in s, or -1 if substr is not present in s.
- `index_any(s string, chars string) => int`: returns the index of the first
  instance of any Unicode code point from chars in s, or -1 if no Unicode code
  point from chars is present in s.
- `join(arr string, sep string) => string`: concatenates the elements of a to
  create a single string. The separator string sep is placed between elements
  in the resulting string.
- `last_index(s string, substr string) => int`: returns the index of the last
  instance of substr in s, or -1 if substr is not present in s.
- `last_index_any(s string, chars string) => int`: returns the index of the
  last instance of any Unicode code point from chars in s, or -1 if no Unicode
  code point from chars is present in s.
- `repeat(s string, count int) => string`: returns a new string consisting of
  count copies of the string s.
- `replace(s string, old string, new string, n int) => string`: returns a copy
  of the string s with the first n non-overlapping instances of old replaced by
  new.
- `substr(s string, lower int, upper int) => string => string`: returns a
  substring of the string s specified by the lower and upper parameters.
- `split(s string, sep string) => [string]`: slices s into all substrings
  separated by sep and returns a slice of the substrings between those
  separators.
- `split_after(s string, sep string) => [string]`: slices s into all substrings
  after each instance of sep and returns a slice of those substrings.
- `split_after_n(s string, sep string, n int) => [string]`: slices s into
  substrings after each instance of sep and returns a slice of those substrings.
- `split_n(s string, sep string, n int) => [string]`: slices s into substrings
  separated by sep and returns a slice of the substrings between those
  separators.
- `title(s string) => string`: returns a copy of the string s with all Unicode
  letters that begin words mapped to their title case.
- `to_lower(s string) => string`: returns a copy of the string s with all
  Unicode letters mapped to their lower case.
- `to_title(s string) => string`: returns a copy of the string s with all
  Unicode letters mapped to their title case.
- `to_upper(s string) => string`: returns a copy of the string s with all
  Unicode letters mapped to their upper case.
- `pad_left(s string, pad_len int, pad_with string) => string`: returns a copy
  of the string s padded on the left with the contents of the string pad_with
  to length pad_len. If pad_with is not specified, white space is used as the
  default padding.
- `pad_right(s string, pad_len int, pad_with string) => string`: returns a
  copy of the string s padded on the right with the contents of the string
  pad_with to length pad_len. If pad_with is not specified, white space is
  used as the default padding.
- `trim(s string, cutset string) => string`: returns a slice of the string s
  with all leading and trailing Unicode code points contained in cutset removed.
- `trim_left(s string, cutset string) => string`: returns a slice of the string
  s with all leading Unicode code points contained in cutset removed.
- `trim_prefix(s string, prefix string) => string`: returns s without the
  provided leading prefix string.
- `trim_right(s string, cutset string) => string`: returns a slice of the
  string s, with all trailing Unicode code points contained in cutset removed.
- `trim_space(s string) => string`: returns a slice of the string s, with all
  leading and trailing white space removed, as defined by Unicode.
- `trim_suffix(s string, suffix string) => string`: returns s without the
  provided trailing suffix string.
- `atoi(str string) => int/error`: returns the result of ParseInt(s, 10, 0)
  converted to type int.
- `format_bool(b bool) => string`: returns "true" or "false" according to the
  value of b.
- `format_float(f float, fmt string, prec int, bits int) => string`: converts
  the floating-point number f to a string, according to the format fmt and
  precision prec.
- `format_int(i int, base int) => string`: returns the string representation of
  i in the given base, for 2 <= base <= 36. The result uses the lower-case
  letters 'a' to 'z' for digit values >= 10.
- `itoa(i int) => string`: is shorthand for format_int(i, 10).
- `parse_bool(s string) => bool/error`: returns the boolean value represented
  by the string. It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false,
  False. Any other value returns an error.
- `parse_float(s string, bits int) => float/error`: converts the string s to a
  floating-point number with the precision specified by bitSize: 32 for float32,
  or 64 for float64. When bitSize=32, the result still has type float64, but it
  will be convertible to float32 without changing its value.
- `parse_int(s string, base int, bits int) => int/error`: interprets a string s
  in the given base (0, 2 to 36) and bit size (0 to 64) and returns the
  corresponding value i.
- `quote(s string) => string`: returns a double-quoted Go string literal
  representing s. The returned string uses Go escape sequences (\t, \n, \xFF,
  \u0100) for control characters and non-printable characters as defined by
  IsPrint.
- `unquote(s string) => string/error`: interprets s as a single-quoted,
  double-quoted, or backquoted Go string literal, returning the string value
  that s quotes.  (If s is single-quoted, it would be a Go character literal;
  Unquote returns the corresponding one-character string.)

## Regexp

- `match(text string) => bool`: reports whether the string s contains any match
  of the regular expression pattern.
- `find(text string, count int) => [[{text: string, begin: int, end: int}]]/undefined`:
  returns an array holding all matches, each of which is an array of map object
  that contains matching text, begin and end (exclusive) index.
- `replace(src string, repl string) => string`: returns a copy of src,
  replacing matches of the pattern with the replacement string repl.
- `split(text string, count int) => [string]`: slices s into substrings
  separated by the expression and returns a slice of the substrings between
  those expression matches.
