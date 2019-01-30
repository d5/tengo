package stdlib

import (
	"time"

	"github.com/d5/tengo/objects"
)

var timesModule = map[string]objects.Object{
	"format_ansic":         &objects.String{Value: time.ANSIC},
	"format_unix_date":     &objects.String{Value: time.UnixDate},
	"format_ruby_date":     &objects.String{Value: time.RubyDate},
	"format_rfc822":        &objects.String{Value: time.RFC822},
	"format_rfc822z":       &objects.String{Value: time.RFC822Z},
	"format_rfc850":        &objects.String{Value: time.RFC850},
	"format_rfc1123":       &objects.String{Value: time.RFC1123},
	"format_rfc1123z":      &objects.String{Value: time.RFC1123Z},
	"format_rfc3339":       &objects.String{Value: time.RFC3339},
	"format_rfc3339_nano":  &objects.String{Value: time.RFC3339Nano},
	"format_kitchen":       &objects.String{Value: time.Kitchen},
	"format_stamp":         &objects.String{Value: time.Stamp},
	"format_stamp_milli":   &objects.String{Value: time.StampMilli},
	"format_stamp_micro":   &objects.String{Value: time.StampMicro},
	"format_stamp_nano":    &objects.String{Value: time.StampNano},
	"nanosecond":           &objects.Int{Value: int64(time.Nanosecond)},
	"microsecond":          &objects.Int{Value: int64(time.Microsecond)},
	"millisecond":          &objects.Int{Value: int64(time.Millisecond)},
	"second":               &objects.Int{Value: int64(time.Second)},
	"minute":               &objects.Int{Value: int64(time.Minute)},
	"hour":                 &objects.Int{Value: int64(time.Hour)},
	"january":              &objects.Int{Value: int64(time.January)},
	"february":             &objects.Int{Value: int64(time.February)},
	"march":                &objects.Int{Value: int64(time.March)},
	"april":                &objects.Int{Value: int64(time.April)},
	"may":                  &objects.Int{Value: int64(time.May)},
	"june":                 &objects.Int{Value: int64(time.June)},
	"july":                 &objects.Int{Value: int64(time.July)},
	"august":               &objects.Int{Value: int64(time.August)},
	"september":            &objects.Int{Value: int64(time.September)},
	"october":              &objects.Int{Value: int64(time.October)},
	"november":             &objects.Int{Value: int64(time.November)},
	"december":             &objects.Int{Value: int64(time.December)},
	"sleep":                &objects.UserFunction{Value: timesSleep},               // sleep(int)
	"parse_duration":       &objects.UserFunction{Value: timesParseDuration},       // parse_duration(str) => int
	"since":                &objects.UserFunction{Value: timesSince},               // since(time) => int
	"until":                &objects.UserFunction{Value: timesUntil},               // until(time) => int
	"duration_hours":       &objects.UserFunction{Value: timesDurationHours},       // duration_hours(int) => float
	"duration_minutes":     &objects.UserFunction{Value: timesDurationMinutes},     // duration_minutes(int) => float
	"duration_nanoseconds": &objects.UserFunction{Value: timesDurationNanoseconds}, // duration_nanoseconds(int) => int
	"duration_seconds":     &objects.UserFunction{Value: timesDurationSeconds},     // duration_seconds(int) => float
	"duration_string":      &objects.UserFunction{Value: timesDurationString},      // duration_string(int) => string
	"month_string":         &objects.UserFunction{Value: timesMonthString},         // month_string(int) => string
	"date":                 &objects.UserFunction{Value: timesDate},                // date(year, month, day, hour, min, sec, nsec) => time
	"now":                  &objects.UserFunction{Value: timesNow},                 // now() => time
	"parse":                &objects.UserFunction{Value: timesParse},               // parse(format, str) => time
	"unix":                 &objects.UserFunction{Value: timesUnix},                // unix(sec, nsec) => time
	"add":                  &objects.UserFunction{Value: timesAdd},                 // add(time, int) => time
	"add_date":             &objects.UserFunction{Value: timesAddDate},             // add_date(time, years, months, days) => time
	"sub":                  &objects.UserFunction{Value: timesSub},                 // sub(t time, u time) => int
	"after":                &objects.UserFunction{Value: timesAfter},               // after(t time, u time) => bool
	"before":               &objects.UserFunction{Value: timesBefore},              // before(t time, u time) => bool
	"time_year":            &objects.UserFunction{Value: timesTimeYear},            // time_year(time) => int
	"time_month":           &objects.UserFunction{Value: timesTimeMonth},           // time_month(time) => int
	"time_day":             &objects.UserFunction{Value: timesTimeDay},             // time_day(time) => int
	"time_weekday":         &objects.UserFunction{Value: timesTimeWeekday},         // time_weekday(time) => int
	"time_hour":            &objects.UserFunction{Value: timesTimeHour},            // time_hour(time) => int
	"time_minute":          &objects.UserFunction{Value: timesTimeMinute},          // time_minute(time) => int
	"time_second":          &objects.UserFunction{Value: timesTimeSecond},          // time_second(time) => int
	"time_nanosecond":      &objects.UserFunction{Value: timesTimeNanosecond},      // time_nanosecond(time) => int
	"time_unix":            &objects.UserFunction{Value: timesTimeUnix},            // time_unix(time) => int
	"time_unix_nano":       &objects.UserFunction{Value: timesTimeUnixNano},        // time_unix_nano(time) => int
	"time_format":          &objects.UserFunction{Value: timesTimeFormat},          // time_format(time, format) => string
	"time_location":        &objects.UserFunction{Value: timesTimeLocation},        // time_location(time) => string
	"time_string":          &objects.UserFunction{Value: timesTimeString},          // time_string(time) => string
	"is_zero":              &objects.UserFunction{Value: timesIsZero},              // is_zero(time) => bool
	"to_local":             &objects.UserFunction{Value: timesToLocal},             // to_local(time) => time
	"to_utc":               &objects.UserFunction{Value: timesToUTC},               // to_utc(time) => time
}

func timesSleep(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	i1, ok := objects.ToInt64(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	time.Sleep(time.Duration(i1))
	ret = objects.UndefinedValue

	return
}

func timesParseDuration(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	s1, ok := objects.ToString(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	dur, err := time.ParseDuration(s1)
	if err != nil {
		ret = wrapError(err)
		return
	}

	ret = &objects.Int{Value: int64(dur)}

	return
}

func timesSince(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(time.Since(t1))}

	return
}

func timesUntil(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(time.Until(t1))}

	return
}

func timesDurationHours(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	i1, ok := objects.ToInt64(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Float{Value: time.Duration(i1).Hours()}

	return
}

func timesDurationMinutes(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	i1, ok := objects.ToInt64(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Float{Value: time.Duration(i1).Minutes()}

	return
}

func timesDurationNanoseconds(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	i1, ok := objects.ToInt64(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: time.Duration(i1).Nanoseconds()}

	return
}

func timesDurationSeconds(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	i1, ok := objects.ToInt64(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Float{Value: time.Duration(i1).Seconds()}

	return
}

func timesDurationString(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	i1, ok := objects.ToInt64(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.String{Value: time.Duration(i1).String()}

	return
}

func timesMonthString(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	i1, ok := objects.ToInt64(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.String{Value: time.Month(i1).String()}

	return
}

func timesDate(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 7 {
		err = objects.ErrWrongNumArguments
		return
	}

	i1, ok := objects.ToInt(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}
	i2, ok := objects.ToInt(args[1])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}
	i3, ok := objects.ToInt(args[2])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}
	i4, ok := objects.ToInt(args[3])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}
	i5, ok := objects.ToInt(args[4])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}
	i6, ok := objects.ToInt(args[5])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}
	i7, ok := objects.ToInt(args[6])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Time{Value: time.Date(i1, time.Month(i2), i3, i4, i5, i6, i7, time.Now().Location())}

	return
}

func timesNow(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 0 {
		err = objects.ErrWrongNumArguments
		return
	}

	ret = &objects.Time{Value: time.Now()}

	return
}

func timesParse(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = objects.ErrWrongNumArguments
		return
	}

	s1, ok := objects.ToString(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	s2, ok := objects.ToString(args[1])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	parsed, err := time.Parse(s1, s2)
	if err != nil {
		ret = wrapError(err)
		return
	}

	ret = &objects.Time{Value: parsed}

	return
}

func timesUnix(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = objects.ErrWrongNumArguments
		return
	}

	i1, ok := objects.ToInt64(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	i2, ok := objects.ToInt64(args[1])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Time{Value: time.Unix(i1, i2)}

	return
}

func timesAdd(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	i2, ok := objects.ToInt64(args[1])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Time{Value: t1.Add(time.Duration(i2))}

	return
}

func timesSub(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	t2, ok := objects.ToTime(args[1])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Sub(t2))}

	return
}

func timesAddDate(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 4 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	i2, ok := objects.ToInt(args[1])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	i3, ok := objects.ToInt(args[2])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	i4, ok := objects.ToInt(args[3])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Time{Value: t1.AddDate(i2, i3, i4)}

	return
}

func timesAfter(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	t2, ok := objects.ToTime(args[1])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	if t1.After(t2) {
		ret = objects.TrueValue
	} else {
		ret = objects.FalseValue
	}

	return
}

func timesBefore(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	t2, ok := objects.ToTime(args[1])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	if t1.Before(t2) {
		ret = objects.TrueValue
	} else {
		ret = objects.FalseValue
	}

	return
}

func timesTimeYear(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Year())}

	return
}

func timesTimeMonth(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Month())}

	return
}

func timesTimeDay(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Day())}

	return
}

func timesTimeWeekday(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Weekday())}

	return
}

func timesTimeHour(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Hour())}

	return
}

func timesTimeMinute(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Minute())}

	return
}

func timesTimeSecond(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Second())}

	return
}

func timesTimeNanosecond(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Nanosecond())}

	return
}

func timesTimeUnix(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.Unix())}

	return
}

func timesTimeUnixNano(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Int{Value: int64(t1.UnixNano())}

	return
}

func timesTimeFormat(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	s2, ok := objects.ToString(args[1])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.String{Value: t1.Format(s2)}

	return
}

func timesIsZero(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	if t1.IsZero() {
		ret = objects.TrueValue
	} else {
		ret = objects.FalseValue
	}

	return
}

func timesToLocal(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Time{Value: t1.Local()}

	return
}

func timesToUTC(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.Time{Value: t1.UTC()}

	return
}

func timesTimeLocation(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.String{Value: t1.Location().String()}

	return
}

func timesTimeString(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}

	t1, ok := objects.ToTime(args[0])
	if !ok {
		err = objects.ErrInvalidTypeConversion
		return
	}

	ret = &objects.String{Value: t1.String()}

	return
}
