package stdlib

import (
	"time"

	"github.com/d5/tengo/objects"
)

var timesModule = map[string]objects.Object{
	// time format constants
	"format_ansic":        &objects.String{Value: time.ANSIC},
	"format_unix_date":    &objects.String{Value: time.UnixDate},
	"format_ruby_date":    &objects.String{Value: time.RubyDate},
	"format_rfc822":       &objects.String{Value: time.RFC822},
	"format_rfc822z":      &objects.String{Value: time.RFC822Z},
	"format_rfc850":       &objects.String{Value: time.RFC850},
	"format_rfc1123":      &objects.String{Value: time.RFC1123},
	"format_rfc1123z":     &objects.String{Value: time.RFC1123Z},
	"format_rfc3339":      &objects.String{Value: time.RFC3339},
	"format_rfc3339_nano": &objects.String{Value: time.RFC3339Nano},
	"format_kitchen":      &objects.String{Value: time.Kitchen},
	"format_stamp":        &objects.String{Value: time.Stamp},
	"format_stamp_milli":  &objects.String{Value: time.StampMilli},
	"format_stamp_micro":  &objects.String{Value: time.StampMicro},
	"format_stamp_nano":   &objects.String{Value: time.StampNano},
	// duration constants
	"nanosecond":  &objects.Int{Value: int64(time.Nanosecond)},
	"microsecond": &objects.Int{Value: int64(time.Microsecond)},
	"millisecond": &objects.Int{Value: int64(time.Millisecond)},
	"second":      &objects.Int{Value: int64(time.Second)},
	"minute":      &objects.Int{Value: int64(time.Minute)},
	"hour":        &objects.Int{Value: int64(time.Hour)},
	// month constants
	"january":   &objects.Int{Value: int64(time.January)},
	"february":  &objects.Int{Value: int64(time.February)},
	"march":     &objects.Int{Value: int64(time.March)},
	"april":     &objects.Int{Value: int64(time.April)},
	"may":       &objects.Int{Value: int64(time.May)},
	"june":      &objects.Int{Value: int64(time.June)},
	"july":      &objects.Int{Value: int64(time.July)},
	"august":    &objects.Int{Value: int64(time.August)},
	"september": &objects.Int{Value: int64(time.September)},
	"october":   &objects.Int{Value: int64(time.October)},
	"november":  &objects.Int{Value: int64(time.November)},
	"december":  &objects.Int{Value: int64(time.December)},
	// sleep(int)
	"sleep": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// parse_duration(str) => int
	"parse_duration": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// since(time) => int
	"since": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// until(time) => int
	"until": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// duration_hours(int) => float
	"duration_hours": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// duration_minutes(int) => float
	"duration_minutes": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// duration_nanoseconds(int) => int
	"duration_nanoseconds": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// duration_seconds(int) => float
	"duration_seconds": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// duration_string(int) => string
	"duration_string": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// month_string(int) => string
	"month_string": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// date(year, month, day, hour, min, sec, nsec) => time
	"date": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// now() => time
	"now": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				err = objects.ErrWrongNumArguments
				return
			}

			ret = &objects.Time{Value: time.Now()}

			return
		},
	},
	// parse(format, str) => time
	"parse": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// unix(sec, nsec) => time
	"unix": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// add(time, int) => time
	"add": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// add_date(time, years, months, days) => time
	"add_date": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// after(t time, u time) => bool
	"after": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// before(t time, u time) => bool
	"before": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_day(time) => int
	"time_day": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_hour(time) => int
	"time_hour": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_year(time) => int
	"time_year": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_month(time) => int
	"time_month": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_minute(time) => int
	"time_minute": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_second(time) => int
	"time_second": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_nanosecond(time) => int
	"time_nanosecond": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_unix(time) => int
	"time_unix": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_unix_nano(time) => int
	"time_unix_nano": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_format(time, format) => string
	"time_format": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// is_zero(time) => bool
	"is_zero": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// to_local(time) => time
	"to_local": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// to_utc(time) => time
	"to_utc": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// location(time) => string
	"time_location": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// time_string(time) => string
	"time_string": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
	// sub(t time, u time) => int
	"sub": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
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
		},
	},
}
