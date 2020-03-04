# Module - "times"

```golang
times := import("times")
```

## Constants

- `format_ansic`: time format "Mon Jan _2 15:04:05 2006"
- `format_unix_date`: time format "Mon Jan _2 15:04:05 MST 2006"
- `format_ruby_date`: time format "Mon Jan 02 15:04:05 -0700 2006"
- `format_rfc822`: time format "02 Jan 06 15:04 MST"
- `format_rfc822z`: time format "02 Jan 06 15:04 -0700"
- `format_rfc850`: time format "Monday, 02-Jan-06 15:04:05 MST"
- `format_rfc1123`: time format "Mon, 02 Jan 2006 15:04:05 MST"
- `format_rfc1123z`: time format "Mon, 02 Jan 2006 15:04:05 -0700"
- `format_rfc3339`: time format "2006-01-02T15:04:05Z07:00"
- `format_rfc3339_nano`: time format "2006-01-02T15:04:05.999999999Z07:00"
- `format_kitchen`: time format "3:04PM"
- `format_stamp`: time format "Jan _2 15:04:05"
- `format_stamp_milli`: time format "Jan _2 15:04:05.000"
- `format_stamp_micro`: time format "Jan _2 15:04:05.000000"
- `format_stamp_nano`: time format "Jan _2 15:04:05.000000000"
- `nanosecond`
- `microsecond`
- `millisecond`
- `second`
- `minute`
- `hour`
- `january`
- `february`
- `march`
- `april`
- `may`
- `june`
- `july`
- `august`
- `september`
- `october`
- `november`
- `december`

## Functions

- `sleep(duration int)`: pauses the current goroutine for at least the duration
  d. A negative or zero duration causes Sleep to return immediately.
- `parse_duration(s string) => int`: parses a duration string. A duration
  string is a possibly signed sequence of decimal numbers, each with optional
  fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time
  units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
- `since(t time) => int`: returns the time elapsed since t.
- `until(t time) => int`: returns the duration until t.
- `duration_hours(duration int) => float`: returns the duration as a floating
  point number of hours.
- `duration_minutes(duration int) => float`: returns the duration as a floating
  point number of minutes.
- `duration_nanoseconds(duration int) => int`: returns the duration as an
  integer of nanoseconds.
- `duration_seconds(duration int) => float`: returns the duration as a floating
  point number of seconds.
- `duration_string(duration int) => string`: returns a string representation of
  duration.
- `month_string(month int) => string`:  returns the English name of the month
  ("January", "February", ...).
- `date(year int, month int, day int, hour int, min int, sec int, nsec int) => time`:
  returns the Time corresponding to "yyyy-mm-dd hh:mm:ss + nsec nanoseconds".
  Current location is used.
- `now() => time`: returns the current local time.
- `parse(format string, s string) => time`: parses a formatted string and
  returns the time value it represents. The layout defines the format by
  showing how the reference time, defined to be "Mon Jan 2 15:04:05 -0700 MST
  2006" would be interpreted if it were the value; it serves as an example of
  the input format. The same interpretation will then be made to the input
  string.
- `unix(sec int, nsec int) => time`: returns the local Time corresponding to
  the given Unix time, sec seconds and nsec nanoseconds since January 1,
  1970 UTC.
- `add(t time, duration int) => time`: returns the time t+d.
- `add_date(t time, years int, months int, days int) => time`: returns the time
  corresponding to adding the given number of years, months, and days to t. For
  example, AddDate(-1, 2, 3) applied to January 1, 2011 returns March 4, 2010.
- `sub(t time, u time) => int`: returns the duration t-u.
- `after(t time, u time) => bool`: reports whether the time instant t is after
  u.
- `before(t time, u time) => bool`: reports whether the time instant t is
  before u.
- `time_year(t time) => int`: returns the year in which t occurs.
- `time_month(t time) => int`: returns the month of the year specified by t.
- `time_day(t time) => int`: returns the day of the month specified by t.
- `time_weekday(t time) => int`: returns the day of the week specified by t.
- `time_hour(t time) => int`: returns the hour within the day specified by t,
  in the range [0, 23].
- `time_minute(t time) => int`: returns the minute offset within the hour
  specified by t, in the range [0, 59].
- `time_second(t time) => int`: returns the second offset within the minute
  specified by t, in the range [0, 59].
- `time_nanosecond(t time) => int`: returns the nanosecond offset within the
  second specified by t, in the range [0, 999999999].
- `time_unix(t time) => int`: returns t as a Unix time, the number of seconds
  elapsed since January 1, 1970 UTC. The result does not depend on the location
  associated with t.
- `time_unix_nano(t time) => int`: returns t as a Unix time, the number of
  nanoseconds elapsed since January 1, 1970 UTC. The result is undefined if the
  Unix time in nanoseconds cannot be represented by an int64 (a date before the
  year 1678 or after 2262). Note that this means the result of calling UnixNano
  on the zero Time is undefined. The result does not depend on the location
  associated with t.
- `time_format(t time, format) => string`: returns a textual representation of
  he time value formatted according to layout, which defines the format by
  showing how the reference time, defined to be "Mon Jan 2 15:04:05 -0700 MST
  2006" would be displayed if it were the value; it serves as an example of the
  desired output. The same display rules will then be applied to the time value.
- `time_location(t time) => string`: returns the time zone name associated with
  t.
- `time_string(t time) => string`: returns the time formatted using the format
  string "2006-01-02 15:04:05.999999999 -0700 MST".
- `is_zero(t time) => bool`: reports whether t represents the zero time
  instant, January 1, year 1, 00:00:00 UTC.
- `to_local(t time) => time`: returns t with the location set to local time.
- `to_utc(t time) => time`: returns t with the location set to UTC.
