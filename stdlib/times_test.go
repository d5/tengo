package stdlib_test

import (
	"testing"
	"time"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
)

func TestTimes(t *testing.T) {
	time1 := time.Date(1982, 9, 28, 19, 21, 44, 999, time.Now().Location())
	time2 := time.Now()

	module(t, "times").call("sleep", mockInterop{}, 1).expect(tengo.UndefinedValue)

	assert.True(t, module(t, "times").call("since", mockInterop{}, time.Now().Add(-time.Hour)).o.(*tengo.Int).Value > 3600000000000)
	assert.True(t, module(t, "times").call("until", mockInterop{}, time.Now().Add(time.Hour)).o.(*tengo.Int).Value < 3600000000000)

	module(t, "times").call("parse_duration", mockInterop{}, "1ns").expect(1)
	module(t, "times").call("parse_duration", mockInterop{}, "1ms").expect(1000000)
	module(t, "times").call("parse_duration", mockInterop{}, "1h").expect(3600000000000)
	module(t, "times").call("duration_hours", mockInterop{}, 1800000000000).expect(0.5)
	module(t, "times").call("duration_minutes", mockInterop{}, 1800000000000).expect(30.0)
	module(t, "times").call("duration_nanoseconds", mockInterop{}, 100).expect(100)
	module(t, "times").call("duration_seconds", mockInterop{}, 1000000).expect(0.001)
	module(t, "times").call("duration_string", mockInterop{}, 1800000000000).expect("30m0s")

	module(t, "times").call("month_string", mockInterop{}, 1).expect("January")
	module(t, "times").call("month_string", mockInterop{}, 12).expect("December")

	module(t, "times").call("date", mockInterop{}, 1982, 9, 28, 19, 21, 44, 999).expect(time1)
	nowD := time.Until(module(t, "times").call("now", mockInterop{}).o.(*tengo.Time).Value).Nanoseconds()
	assert.True(t, 0 > nowD && nowD > -100000000) // within 100ms
	parsed, _ := time.Parse(time.RFC3339, "1982-09-28T19:21:44+07:00")
	module(t, "times").call("parse", mockInterop{}, time.RFC3339, "1982-09-28T19:21:44+07:00").expect(parsed)
	module(t, "times").call("unix", mockInterop{}, 1234325, 94493).expect(time.Unix(1234325, 94493))

	module(t, "times").call("add", mockInterop{}, time2, 3600000000000).expect(time2.Add(time.Duration(3600000000000)))
	module(t, "times").call("sub", mockInterop{}, time2, time2.Add(-time.Hour)).expect(3600000000000)
	module(t, "times").call("add_date", mockInterop{}, time2, 1, 2, 3).expect(time2.AddDate(1, 2, 3))
	module(t, "times").call("after", mockInterop{}, time2, time2.Add(time.Hour)).expect(false)
	module(t, "times").call("after", mockInterop{}, time2, time2.Add(-time.Hour)).expect(true)
	module(t, "times").call("before", mockInterop{}, time2, time2.Add(time.Hour)).expect(true)
	module(t, "times").call("before", mockInterop{}, time2, time2.Add(-time.Hour)).expect(false)

	module(t, "times").call("time_year", mockInterop{}, time1).expect(time1.Year())
	module(t, "times").call("time_month", mockInterop{}, time1).expect(int(time1.Month()))
	module(t, "times").call("time_day", mockInterop{}, time1).expect(time1.Day())
	module(t, "times").call("time_hour", mockInterop{}, time1).expect(time1.Hour())
	module(t, "times").call("time_minute", mockInterop{}, time1).expect(time1.Minute())
	module(t, "times").call("time_second", mockInterop{}, time1).expect(time1.Second())
	module(t, "times").call("time_nanosecond", mockInterop{}, time1).expect(time1.Nanosecond())
	module(t, "times").call("time_unix", mockInterop{}, time1).expect(time1.Unix())
	module(t, "times").call("time_unix_nano", mockInterop{}, time1).expect(time1.UnixNano())
	module(t, "times").call("time_format", mockInterop{}, time1, time.RFC3339).expect(time1.Format(time.RFC3339))
	module(t, "times").call("is_zero", mockInterop{}, time1).expect(false)
	module(t, "times").call("is_zero", mockInterop{}, time.Time{}).expect(true)
	module(t, "times").call("to_local", mockInterop{}, time1).expect(time1.Local())
	module(t, "times").call("to_utc", mockInterop{}, time1).expect(time1.UTC())
	module(t, "times").call("time_location", mockInterop{}, time1).expect(time1.Location().String())
	module(t, "times").call("time_string", mockInterop{}, time1).expect(time1.String())
}
