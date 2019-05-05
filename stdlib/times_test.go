package stdlib_test

import (
	"testing"
	"time"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
)

func TestTimes(t *testing.T) {
	time1 := time.Date(1982, 9, 28, 19, 21, 44, 999, time.Now().Location())
	time2 := time.Now()

	module(t, "times").call("sleep", mockRuntime{},  1).expect(objects.UndefinedValue)

	assert.True(t, module(t, "times").call("since", mockRuntime{},  time.Now().Add(-time.Hour)).o.(*objects.Int).Value > 3600000000000)
	assert.True(t, module(t, "times").call("until", mockRuntime{},  time.Now().Add(time.Hour)).o.(*objects.Int).Value < 3600000000000)

	module(t, "times").call("parse_duration", mockRuntime{},  "1ns").expect(1)
	module(t, "times").call("parse_duration", mockRuntime{},  "1ms").expect(1000000)
	module(t, "times").call("parse_duration", mockRuntime{},  "1h").expect(3600000000000)
	module(t, "times").call("duration_hours", mockRuntime{},  1800000000000).expect(0.5)
	module(t, "times").call("duration_minutes", mockRuntime{},  1800000000000).expect(30.0)
	module(t, "times").call("duration_nanoseconds", mockRuntime{},  100).expect(100)
	module(t, "times").call("duration_seconds", mockRuntime{},  1000000).expect(0.001)
	module(t, "times").call("duration_string", mockRuntime{},  1800000000000).expect("30m0s")

	module(t, "times").call("month_string", mockRuntime{},  1).expect("January")
	module(t, "times").call("month_string", mockRuntime{},  12).expect("December")

	module(t, "times").call("date", mockRuntime{},  1982, 9, 28, 19, 21, 44, 999).expect(time1)
	nowD := time.Until(module(t, "times").call("now", mockRuntime{}).o.(*objects.Time).Value).Nanoseconds()
	assert.True(t, 0 > nowD && nowD > -100000000) // within 100ms
	parsed, _ := time.Parse(time.RFC3339, "1982-09-28T19:21:44+07:00")
	module(t, "times").call("parse", mockRuntime{},  time.RFC3339, "1982-09-28T19:21:44+07:00").expect(parsed)
	module(t, "times").call("unix", mockRuntime{},  1234325, 94493).expect(time.Unix(1234325, 94493))

	module(t, "times").call("add", mockRuntime{},  time2, 3600000000000).expect(time2.Add(time.Duration(3600000000000)))
	module(t, "times").call("sub", mockRuntime{},  time2, time2.Add(-time.Hour)).expect(3600000000000)
	module(t, "times").call("add_date", mockRuntime{},  time2, 1, 2, 3).expect(time2.AddDate(1, 2, 3))
	module(t, "times").call("after", mockRuntime{},  time2, time2.Add(time.Hour)).expect(false)
	module(t, "times").call("after", mockRuntime{},  time2, time2.Add(-time.Hour)).expect(true)
	module(t, "times").call("before", mockRuntime{},  time2, time2.Add(time.Hour)).expect(true)
	module(t, "times").call("before", mockRuntime{},  time2, time2.Add(-time.Hour)).expect(false)

	module(t, "times").call("time_year", mockRuntime{},  time1).expect(time1.Year())
	module(t, "times").call("time_month", mockRuntime{},  time1).expect(int(time1.Month()))
	module(t, "times").call("time_day", mockRuntime{},  time1).expect(time1.Day())
	module(t, "times").call("time_hour", mockRuntime{},  time1).expect(time1.Hour())
	module(t, "times").call("time_minute", mockRuntime{},  time1).expect(time1.Minute())
	module(t, "times").call("time_second", mockRuntime{},  time1).expect(time1.Second())
	module(t, "times").call("time_nanosecond", mockRuntime{},  time1).expect(time1.Nanosecond())
	module(t, "times").call("time_unix", mockRuntime{},  time1).expect(time1.Unix())
	module(t, "times").call("time_unix_nano", mockRuntime{},  time1).expect(time1.UnixNano())
	module(t, "times").call("time_format", mockRuntime{},  time1, time.RFC3339).expect(time1.Format(time.RFC3339))
	module(t, "times").call("is_zero", mockRuntime{},  time1).expect(false)
	module(t, "times").call("is_zero", mockRuntime{},  time.Time{}).expect(true)
	module(t, "times").call("to_local", mockRuntime{},  time1).expect(time1.Local())
	module(t, "times").call("to_utc", mockRuntime{},  time1).expect(time1.UTC())
	module(t, "times").call("time_location", mockRuntime{},  time1).expect(time1.Location().String())
	module(t, "times").call("time_string", mockRuntime{},  time1).expect(time1.String())
}
