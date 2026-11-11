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
	location, _ := time.LoadLocation("Pacific/Auckland")
	time3 := time.Date(1982, 9, 28, 19, 21, 44, 999, location)

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
	module(t, "times").call("date", 1982, 9, 28, 19, 21, 44, 999, "Pacific/Auckland").
		expect(time3)

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

	_, isoweek := time1.ISOWeek()
	module(t, "times").call("time_isoweek", time1).expect(isoweek)
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
	module(t, "times").call("in_location", time1, location.String()).expect(time1.In(location))

	for _, locName := range []string{"Asia/Shanghai", "UTC", "America/New_York"} {
		loc, err := time.LoadLocation(locName)
		require.NoError(t, err)
		module(t, "times").call("in_location", time1, locName).expect(time1.In(loc))
		module(t, "times").call("date", 1982, 9, 28, 19, 21, 44, 999, locName).
			expect(time.Date(1982, 9, 28, 19, 21, 44, 999, loc))
	}
	module(t, "times").call("in_location", time1, "Invalid/Location").expectError()
}

func TestTimesWeekday(t *testing.T) {
	module(t, "times").
		call("time_weekday", time.Date(1982, 9, 28, 19, 21, 44, 999, time.UTC)).
		expect(int(time.Tuesday))
	module(t, "times").
		call("time_weekday", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)).
		expect(0)
}

func TestTimesIntAsTime(t *testing.T) {
	module(t, "times").call("time_unix", 1234567890).expect(1234567890)
	module(t, "times").call("time_unix_nano", 2).expect(2000000000)
}

func TestTimesFormatLimit(t *testing.T) {
	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 10

	t1 := time.Date(1982, 9, 28, 0, 0, 0, 0, time.UTC)
	module(t, "times").call("time_format", t1, "2006-01-02").expect("1982-09-28")
	module(t, "times").call("time_format", t1, "2006-01-02 15").
		expectErr(tengo.ErrStringLimit)
}

func TestTimesArgErrors(t *testing.T) {
	now := time.Now()
	wrongNum := tengo.ErrWrongNumArguments
	invalidTime := func(name string) error {
		return invalidArg(name, "time(compatible)", "string")
	}
	invalidInt := func(name string) error {
		return invalidArg(name, "int(compatible)", "string")
	}
	invalidString := func(name string) error {
		return invalidArg(name, "string(compatible)", "undefined")
	}
	undefined := tengo.UndefinedValue

	var cases []argErrorCase
	for _, fn := range []string{
		"since", "until", "time_year", "time_month", "time_day",
		"time_weekday", "time_hour", "time_minute", "time_second",
		"time_nanosecond", "time_unix", "time_unix_nano", "time_location",
		"time_string", "is_zero", "to_local", "to_utc",
	} {
		cases = append(cases,
			argErrorCase{fn, ARR{}, wrongNum},
			argErrorCase{fn, ARR{now, now}, wrongNum},
			argErrorCase{fn, ARR{"x"}, invalidTime("first")},
		)
	}
	for _, fn := range []string{
		"sleep", "duration_hours", "duration_minutes",
		"duration_nanoseconds", "duration_seconds", "duration_string",
		"month_string",
	} {
		cases = append(cases,
			argErrorCase{fn, ARR{}, wrongNum},
			argErrorCase{fn, ARR{1, 2}, wrongNum},
			argErrorCase{fn, ARR{"x"}, invalidInt("first")},
		)
	}

	date := ARR{1982, 9, 28, 19, 21, 44, 999}
	names := []string{
		"first", "second", "third", "fourth", "fifth", "sixth", "seventh",
	}
	for i, name := range names {
		args := append(ARR{}, date...)
		args[i] = "x"
		cases = append(cases, argErrorCase{"date", args, invalidInt(name)})
	}
	cases = append(cases,
		argErrorCase{"date", date[:6], wrongNum},
		argErrorCase{"date", append(append(ARR{}, date...), "UTC", "UTC"),
			wrongNum},
		argErrorCase{"date", append(append(ARR{}, date...), undefined),
			invalidString("eighth")},
	)

	cases = append(cases, []argErrorCase{
		{"now", ARR{1}, wrongNum},
		{"parse_duration", ARR{}, wrongNum},
		{"parse_duration", ARR{undefined}, invalidString("first")},
		{"parse", ARR{time.RFC3339}, wrongNum},
		{"parse", ARR{undefined, "x"}, invalidString("first")},
		{"parse", ARR{time.RFC3339, undefined}, invalidString("second")},
		{"unix", ARR{1}, wrongNum},
		{"unix", ARR{"x", 1}, invalidInt("first")},
		{"unix", ARR{1, "x"}, invalidInt("second")},
		{"add", ARR{now}, wrongNum},
		{"add", ARR{"x", 1}, invalidTime("first")},
		{"add", ARR{now, "x"}, invalidInt("second")},
		{"sub", ARR{now}, wrongNum},
		{"sub", ARR{"x", now}, invalidTime("first")},
		{"sub", ARR{now, "x"}, invalidTime("second")},
		{"add_date", ARR{now, 1, 2}, wrongNum},
		{"add_date", ARR{"x", 1, 2, 3}, invalidTime("first")},
		{"add_date", ARR{now, "x", 2, 3}, invalidInt("second")},
		{"add_date", ARR{now, 1, "x", 3}, invalidInt("third")},
		{"add_date", ARR{now, 1, 2, "x"}, invalidInt("fourth")},
		{"after", ARR{now}, wrongNum},
		{"after", ARR{"x", now}, invalidTime("first")},
		{"after", ARR{now, "x"}, invalidTime("second")},
		{"before", ARR{now}, wrongNum},
		{"before", ARR{"x", now}, invalidTime("first")},
		{"time_format", ARR{now}, wrongNum},
		{"time_format", ARR{"x", time.RFC3339}, invalidTime("first")},
		{"time_format", ARR{now, undefined}, invalidString("second")},
		{"in_location", ARR{now}, wrongNum},
		{"in_location", ARR{"x", "UTC"}, invalidTime("first")},
		{"in_location", ARR{now, undefined}, invalidString("second")},
	}...)
	expectArgErrors(t, "times", cases)

	res := module(t, "times").call("before", now, "x")
	require.IsType(t, tengo.ErrInvalidArgumentType{}, res.e)
}
