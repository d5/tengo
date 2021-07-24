package stdlib

import (
	"time"

	"github.com/d5/tengo/v2"
)

var timesModule = map[string]tengo.Object{
	"format_ansic":        &tengo.String{Value: time.ANSIC},
	"format_unix_date":    &tengo.String{Value: time.UnixDate},
	"format_ruby_date":    &tengo.String{Value: time.RubyDate},
	"format_rfc822":       &tengo.String{Value: time.RFC822},
	"format_rfc822z":      &tengo.String{Value: time.RFC822Z},
	"format_rfc850":       &tengo.String{Value: time.RFC850},
	"format_rfc1123":      &tengo.String{Value: time.RFC1123},
	"format_rfc1123z":     &tengo.String{Value: time.RFC1123Z},
	"format_rfc3339":      &tengo.String{Value: time.RFC3339},
	"format_rfc3339_nano": &tengo.String{Value: time.RFC3339Nano},
	"format_kitchen":      &tengo.String{Value: time.Kitchen},
	"format_stamp":        &tengo.String{Value: time.Stamp},
	"format_stamp_milli":  &tengo.String{Value: time.StampMilli},
	"format_stamp_micro":  &tengo.String{Value: time.StampMicro},
	"format_stamp_nano":   &tengo.String{Value: time.StampNano},
	"nanosecond":          &tengo.Int{Value: int64(time.Nanosecond)},
	"microsecond":         &tengo.Int{Value: int64(time.Microsecond)},
	"millisecond":         &tengo.Int{Value: int64(time.Millisecond)},
	"second":              &tengo.Int{Value: int64(time.Second)},
	"minute":              &tengo.Int{Value: int64(time.Minute)},
	"hour":                &tengo.Int{Value: int64(time.Hour)},
	"january":             &tengo.Int{Value: int64(time.January)},
	"february":            &tengo.Int{Value: int64(time.February)},
	"march":               &tengo.Int{Value: int64(time.March)},
	"april":               &tengo.Int{Value: int64(time.April)},
	"may":                 &tengo.Int{Value: int64(time.May)},
	"june":                &tengo.Int{Value: int64(time.June)},
	"july":                &tengo.Int{Value: int64(time.July)},
	"august":              &tengo.Int{Value: int64(time.August)},
	"september":           &tengo.Int{Value: int64(time.September)},
	"october":             &tengo.Int{Value: int64(time.October)},
	"november":            &tengo.Int{Value: int64(time.November)},
	"december":            &tengo.Int{Value: int64(time.December)},
	"sleep": &tengo.UserFunction{
		Name:  "sleep",
		Value: timesSleep,
	}, // sleep(int)
	"parse_duration": &tengo.UserFunction{
		Name:  "parse_duration",
		Value: timesParseDuration,
	}, // parse_duration(str) => int
	"since": &tengo.UserFunction{
		Name:  "since",
		Value: timesSince,
	}, // since(time) => int
	"until": &tengo.UserFunction{
		Name:  "until",
		Value: timesUntil,
	}, // until(time) => int
	"duration_hours": &tengo.UserFunction{
		Name:  "duration_hours",
		Value: timesDurationHours,
	}, // duration_hours(int) => float
	"duration_minutes": &tengo.UserFunction{
		Name:  "duration_minutes",
		Value: timesDurationMinutes,
	}, // duration_minutes(int) => float
	"duration_nanoseconds": &tengo.UserFunction{
		Name:  "duration_nanoseconds",
		Value: timesDurationNanoseconds,
	}, // duration_nanoseconds(int) => int
	"duration_seconds": &tengo.UserFunction{
		Name:  "duration_seconds",
		Value: timesDurationSeconds,
	}, // duration_seconds(int) => float
	"duration_string": &tengo.UserFunction{
		Name:  "duration_string",
		Value: timesDurationString,
	}, // duration_string(int) => string
	"month_string": &tengo.UserFunction{
		Name:  "month_string",
		Value: timesMonthString,
	}, // month_string(int) => string
	"date": &tengo.UserFunction{
		Name:  "date",
		Value: timesDate,
	}, // date(year, month, day, hour, min, sec, nsec) => time
	"now": &tengo.UserFunction{
		Name:  "now",
		Value: timesNow,
	}, // now() => time
	"parse": &tengo.UserFunction{
		Name:  "parse",
		Value: timesParse,
	}, // parse(format, str) => time
	"unix": &tengo.UserFunction{
		Name:  "unix",
		Value: timesUnix,
	}, // unix(sec, nsec) => time
	"add": &tengo.UserFunction{
		Name:  "add",
		Value: timesAdd,
	}, // add(time, int) => time
	"add_date": &tengo.UserFunction{
		Name:  "add_date",
		Value: timesAddDate,
	}, // add_date(time, years, months, days) => time
	"sub": &tengo.UserFunction{
		Name:  "sub",
		Value: timesSub,
	}, // sub(t time, u time) => int
	"after": &tengo.UserFunction{
		Name:  "after",
		Value: timesAfter,
	}, // after(t time, u time) => bool
	"before": &tengo.UserFunction{
		Name:  "before",
		Value: timesBefore,
	}, // before(t time, u time) => bool
	"time_year": &tengo.UserFunction{
		Name:  "time_year",
		Value: timesTimeYear,
	}, // time_year(time) => int
	"time_month": &tengo.UserFunction{
		Name:  "time_month",
		Value: timesTimeMonth,
	}, // time_month(time) => int
	"time_day": &tengo.UserFunction{
		Name:  "time_day",
		Value: timesTimeDay,
	}, // time_day(time) => int
	"time_weekday": &tengo.UserFunction{
		Name:  "time_weekday",
		Value: timesTimeWeekday,
	}, // time_weekday(time) => int
	"time_hour": &tengo.UserFunction{
		Name:  "time_hour",
		Value: timesTimeHour,
	}, // time_hour(time) => int
	"time_minute": &tengo.UserFunction{
		Name:  "time_minute",
		Value: timesTimeMinute,
	}, // time_minute(time) => int
	"time_second": &tengo.UserFunction{
		Name:  "time_second",
		Value: timesTimeSecond,
	}, // time_second(time) => int
	"time_nanosecond": &tengo.UserFunction{
		Name:  "time_nanosecond",
		Value: timesTimeNanosecond,
	}, // time_nanosecond(time) => int
	"time_unix": &tengo.UserFunction{
		Name:  "time_unix",
		Value: timesTimeUnix,
	}, // time_unix(time) => int
	"time_unix_nano": &tengo.UserFunction{
		Name:  "time_unix_nano",
		Value: timesTimeUnixNano,
	}, // time_unix_nano(time) => int
	"time_format": &tengo.UserFunction{
		Name:  "time_format",
		Value: timesTimeFormat,
	}, // time_format(time, format) => string
	"time_location": &tengo.UserFunction{
		Name:  "time_location",
		Value: timesTimeLocation,
	}, // time_location(time) => string
	"time_string": &tengo.UserFunction{
		Name:  "time_string",
		Value: timesTimeString,
	}, // time_string(time) => string
	"is_zero": &tengo.UserFunction{
		Name:  "is_zero",
		Value: timesIsZero,
	}, // is_zero(time) => bool
	"to_local": &tengo.UserFunction{
		Name:  "to_local",
		Value: timesToLocal,
	}, // to_local(time) => time
	"to_utc": &tengo.UserFunction{
		Name:  "to_utc",
		Value: timesToUTC,
	}, // to_utc(time) => time
}

var timesSleep = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	i1, err := tengo.ToInt64(0, args...)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(i1))
	return tengo.UndefinedValue, nil
}, 1)

var timesParseDuration = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	dur, err := time.ParseDuration(s1)
	if err != nil {
		return wrapError(err), nil
	}

	return &tengo.Int{Value: int64(dur)}, nil
}, 1)

var timesSince = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(time.Since(t1))}, nil
}, 1)

var timesUntil = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(time.Until(t1))}, nil
}, 1)

var timesDurationHours = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	i1, err := tengo.ToInt64(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Float{Value: time.Duration(i1).Hours()}, nil
}, 1)

var timesDurationMinutes = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	i1, err := tengo.ToInt64(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Float{Value: time.Duration(i1).Minutes()}, nil
}, 1)

var timesDurationNanoseconds = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	i1, err := tengo.ToInt64(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: time.Duration(i1).Nanoseconds()}, nil
}, 1)

var timesDurationSeconds = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	i1, err := tengo.ToInt64(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Float{Value: time.Duration(i1).Seconds()}, nil
}, 1)

var timesDurationString = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	i1, err := tengo.ToInt64(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.String{Value: time.Duration(i1).String()}, nil
}, 1)

var timesMonthString = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	i1, err := tengo.ToInt64(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.String{Value: time.Month(i1).String()}, nil
}, 1)

var timesDate = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {

	i1, err := tengo.ToInt(0, args...)
	if err != nil {
		return nil, err
	}
	i2, err := tengo.ToInt(1, args...)
	if err != nil {
		return nil, err
	}
	i3, err := tengo.ToInt(2, args...)
	if err != nil {
		return nil, err
	}
	i4, err := tengo.ToInt(3, args...)
	if err != nil {
		return nil, err
	}
	i5, err := tengo.ToInt(4, args...)
	if err != nil {
		return nil, err
	}
	i6, err := tengo.ToInt(5, args...)
	if err != nil {
		return nil, err
	}
	i7, err := tengo.ToInt(6, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Time{
		Value: time.Date(i1,
			time.Month(i2), i3, i4, i5, i6, i7, time.Now().Location()),
	}, nil
}, 7)

var timesNow = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	return &tengo.Time{Value: time.Now()}, nil
})

var timesParse = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	s2, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	parsed, err := time.Parse(s1, s2)
	if err != nil {
		return wrapError(err), nil
	}

	return &tengo.Time{Value: parsed}, nil
}, 2)

var timesUnix = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	i1, err := tengo.ToInt64(0, args...)
	if err != nil {
		return nil, err
	}

	i2, err := tengo.ToInt64(1, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Time{Value: time.Unix(i1, i2)}, nil
}, 2)

var timesAdd = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	i2, err := tengo.ToInt64(1, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Time{Value: t1.Add(time.Duration(i2))}, nil
}, 2)

var timesSub = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	t2, err := tengo.ToTime(1, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(t1.Sub(t2))}, nil
}, 2)

var timesAddDate = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	i2, err := tengo.ToInt(1, args...)
	if err != nil {
		return nil, err
	}

	i3, err := tengo.ToInt(2, args...)
	if err != nil {
		return nil, err
	}

	i4, err := tengo.ToInt(3, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Time{Value: t1.AddDate(i2, i3, i4)}, nil
}, 4)

var timesAfter = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	t2, err := tengo.ToTime(1, args...)
	if err != nil {
		return nil, err
	}

	if t1.After(t2) {
		return tengo.TrueValue, nil
	}
	return tengo.FalseValue, nil
}, 2)

var timesBefore = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	t2, err := tengo.ToTime(1, args...)
	if err != nil {
		return nil, err
	}

	if t1.Before(t2) {
		return tengo.TrueValue, nil
	}
	return tengo.FalseValue, nil
}, 2)

var timesTimeYear = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(t1.Year())}, nil
}, 1)

var timesTimeMonth = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(t1.Month())}, nil
}, 1)

var timesTimeDay = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(t1.Day())}, nil
}, 1)

var timesTimeWeekday = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(t1.Weekday())}, nil
}, 1)

var timesTimeHour = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(t1.Hour())}, nil
}, 1)

var timesTimeMinute = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(t1.Minute())}, nil
}, 1)

var timesTimeSecond = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(t1.Second())}, nil
}, 1)

var timesTimeNanosecond = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: int64(t1.Nanosecond())}, nil
}, 1)

var timesTimeUnix = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: t1.Unix()}, nil
}, 1)

var timesTimeUnixNano = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Int{Value: t1.UnixNano()}, nil
}, 1)

var timesTimeFormat = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	s2, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	s := t1.Format(s2)
	if len(s) > tengo.MaxStringLen {
		return nil, tengo.ErrStringLimit
	}

	return &tengo.String{Value: s}, nil
}, 2)

var timesIsZero = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	if t1.IsZero() {
		return tengo.TrueValue, nil
	}
	return tengo.FalseValue, nil
}, 1)

var timesToLocal = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Time{Value: t1.Local()}, nil
}, 1)

var timesToUTC = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.Time{Value: t1.UTC()}, nil
}, 1)

var timesTimeLocation = tengo.CheckAnyArgs(func(args ...tengo.Object) (
	tengo.Object,
	error,
) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.String{Value: t1.Location().String()}, nil
}, 1)

var timesTimeString = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	t1, err := tengo.ToTime(0, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.String{Value: t1.String()}, nil
}, 1)
