package stdlib_test

import (
	"testing"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

func TestTimes(t *testing.T) {
	time1 := time.Date(1982, 9, 28, 19, 21, 44, 999, time.Now().Location())
	time2 := time.Now()

	module(t, "times").call("sleep", 1).expect(tengo.UndefinedValue)

	require.True(t, module(t, "times").
		call("since", time.Now().Add(-time.Hour)).
		o.(*tengo.Int).Value > 3600000000000)
	require.True(t, module(t, "times").
		call("until", time.Now().Add(time.Hour)).
		o.(*tengo.Int).Value < 3600000000000)

	module(t, "times").call("parse_duration", "1ns").expect(1)
	module(t, "times").call("parse_duration", "1ms").expect(1000000)
	module(t, "times").call("parse_duration", "1h").expect(3600000000000)
	module(t, "times").call("duration_hours", 1800000000000).expect(0.5)
	module(t, "times").call("duration_minutes", 1800000000000).expect(30.0)
	module(t, "times").call("duration_nanoseconds", 100).expect(100)
	module(t, "times").call("duration_seconds", 1000000).expect(0.001)
	module(t, "times").call("duration_string", 1800000000000).expect("30m0s")

	module(t, "times").call("month_string", 1).expect("January")
	module(t, "times").call("month_string", 12).expect("December")

	module(t, "times").call("date", 1982, 9, 28, 19, 21, 44, 999).
		expect(time1)
	nowD := time.Until(module(t, "times").call("now").
		o.(*tengo.Time).Value).Nanoseconds()
	require.True(t, 0 > nowD && nowD > -100000000) // within 100ms
	parsed, _ := time.Parse(time.RFC3339, "1982-09-28T19:21:44+07:00")
	module(t, "times").
		call("parse", time.RFC3339, "1982-09-28T19:21:44+07:00").
		expect(parsed)
	module(t, "times").
		call("unix", 1234325, 94493).
		expect(time.Unix(1234325, 94493))

	module(t, "times").call("add", time2, 3600000000000).
		expect(time2.Add(time.Duration(3600000000000)))
	module(t, "times").call("sub", time2, time2.Add(-time.Hour)).
		expect(3600000000000)
	module(t, "times").call("add_date", time2, 1, 2, 3).
		expect(time2.AddDate(1, 2, 3))
	module(t, "times").call("after", time2, time2.Add(time.Hour)).
		expect(false)
	module(t, "times").call("after", time2, time2.Add(-time.Hour)).
		expect(true)
	module(t, "times").call("before", time2, time2.Add(time.Hour)).
		expect(true)
	module(t, "times").call("before", time2, time2.Add(-time.Hour)).
		expect(false)

	module(t, "times").call("time_year", time1).expect(time1.Year())
	module(t, "times").call("time_month", time1).expect(int(time1.Month()))
	module(t, "times").call("time_day", time1).expect(time1.Day())
	module(t, "times").call("time_hour", time1).expect(time1.Hour())
	module(t, "times").call("time_minute", time1).expect(time1.Minute())
	module(t, "times").call("time_second", time1).expect(time1.Second())
	module(t, "times").call("time_nanosecond", time1).
		expect(time1.Nanosecond())
	module(t, "times").call("time_unix", time1).expect(time1.Unix())
	module(t, "times").call("time_unix_nano", time1).expect(time1.UnixNano())
	module(t, "times").call("time_format", time1, time.RFC3339).
		expect(time1.Format(time.RFC3339))
	module(t, "times").call("is_zero", time1).expect(false)
	module(t, "times").call("is_zero", time.Time{}).expect(true)
	module(t, "times").call("to_local", time1).expect(time1.Local())
	module(t, "times").call("to_utc", time1).expect(time1.UTC())
	module(t, "times").call("time_location", time1).
		expect(time1.Location().String())
	module(t, "times").call("time_string", time1).expect(time1.String())
}
