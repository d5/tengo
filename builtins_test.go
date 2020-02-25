package tengo

import (
	"errors"
	"reflect"
	"testing"
)

func Test_builtinDelete(t *testing.T) {
	type args struct {
		args []Object
	}
	tests := []struct {
		name      string
		args      args
		want      Object
		wantErr   bool
		wantedErr error
		target    interface{}
	}{
		//Map
		{name: "invalid-arg", args: args{[]Object{&String{}, &String{}}}, wantErr: true,
			wantedErr: ErrInvalidArgumentType{Name: "first", Expected: "map", Found: "string"}},
		{name: "no-args", wantErr: true, wantedErr: ErrWrongNumArguments},
		{name: "empty-args", args: args{[]Object{}}, wantErr: true, wantedErr: ErrWrongNumArguments},
		{name: "3-args", args: args{[]Object{(*Map)(nil), (*String)(nil), (*String)(nil)}}, wantErr: true, wantedErr: ErrWrongNumArguments},
		{name: "nil-map-empty-key", args: args{[]Object{&Map{}, &String{}}}, want: UndefinedValue},
		{name: "nil-map-nonstr-key", args: args{[]Object{&Map{}, &Int{}}}, wantErr: true,
			wantedErr: ErrInvalidArgumentType{Name: "second", Expected: "string", Found: "int"}},
		{name: "nil-map-no-key", args: args{[]Object{&Map{}}}, wantErr: true,
			wantedErr: ErrWrongNumArguments},
		{name: "map-missing-key",
			args: args{
				[]Object{
					&Map{Value: map[string]Object{
						"key": &String{Value: "value"},
					}},
					&String{Value: "key1"},
				}},
			want:   UndefinedValue,
			target: &Map{Value: map[string]Object{"key": &String{Value: "value"}}},
		},
		{name: "map-emptied",
			args: args{
				[]Object{
					&Map{Value: map[string]Object{
						"key": &String{Value: "value"},
					}},
					&String{Value: "key"},
				}},
			want:   UndefinedValue,
			target: &Map{Value: map[string]Object{}},
		},
		{name: "map-multi-keys",
			args: args{
				[]Object{
					&Map{Value: map[string]Object{
						"key1": &String{Value: "value1"},
						"key2": &Int{Value: 10},
					}},
					&String{Value: "key1"},
				}},
			want:   UndefinedValue,
			target: &Map{Value: map[string]Object{"key2": &Int{Value: 10}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinDelete(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantedErr) {
				if err.Error() != tt.wantedErr.Error() {
					t.Errorf("builtinDelete() error = %v, wantedErr %v", err, tt.wantedErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("builtinDelete() = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && tt.target != nil {
				switch v := tt.args.args[0].(type) {
				case *Map, *Array:
					if !reflect.DeepEqual(tt.target, tt.args.args[0]) {
						t.Errorf("builtinDelete() objects are not equal got: %+v, want: %+v", tt.args.args[0], tt.target)
					}
				default:
					t.Errorf("builtinDelete() unsuporrted arg[0] type %s", v.TypeName())
					return
				}
			}
		})
	}
}

func Test_builtinSplice(t *testing.T) {
	tests := []struct {
		name      string
		args      []Object
		deleted   Object
		Array     *Array
		wantErr   bool
		wantedErr error
	}{
		{name: "no args", args: []Object{}, wantErr: true, wantedErr: ErrWrongNumArguments},
		{name: "invalid args", args: []Object{&Map{}},
			wantErr: true, wantedErr: ErrInvalidArgumentType{Name: "first", Expected: "array", Found: "map"}},
		{name: "invalid args", args: []Object{&Array{}, &String{}},
			wantErr: true, wantedErr: ErrInvalidArgumentType{Name: "second", Expected: "int", Found: "string"}},
		{name: "negative index", args: []Object{&Array{}, &Int{Value: -1}},
			wantErr: true, wantedErr: ErrInvalidArgumentType{Name: "second", Expected: "non-negative int", Found: "negative int"}},
		{name: "non int count", args: []Object{&Array{}, &Int{Value: 0}, &String{Value: ""}},
			wantErr: true, wantedErr: ErrInvalidArgumentType{Name: "third", Expected: "int", Found: "string"}},
		{name: "negative count",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 0},
				&Int{Value: -1},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
		},
		{name: "push with negative count",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 0},
				&Int{Value: -1},
				&String{Value: "a"},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&String{Value: "a"}, &Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
		},
		{name: "insert with zero count",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 0},
				&Int{Value: 0},
				&String{Value: "b"},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&String{Value: "b"}, &Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
		},
		{name: "insert",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 1},
				&Int{Value: 0},
				&String{Value: "c"},
				&String{Value: "d"},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &String{Value: "c"}, &String{Value: "d"}, &Int{Value: 1}, &Int{Value: 2}}},
		},
		{name: "insert with negative count",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 1},
				&Int{Value: -1},
				&String{Value: "c"},
				&String{Value: "d"},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &String{Value: "c"}, &String{Value: "d"}, &Int{Value: 1}, &Int{Value: 2}}},
		},
		{name: "insert with delete",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 1},
				&Int{Value: 1},
				&String{Value: "c"},
				&String{Value: "d"},
			},
			deleted: &Array{Value: []Object{&Int{Value: 1}}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &String{Value: "c"}, &String{Value: "d"}, &Int{Value: 2}}},
		},
		{name: "insert with delete multi",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 1},
				&Int{Value: 2},
				&String{Value: "c"},
				&String{Value: "d"},
			},
			deleted: &Array{Value: []Object{&Int{Value: 1}, &Int{Value: 2}}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &String{Value: "c"}, &String{Value: "d"}}},
		},
		{name: "append with big index",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 5},
				&Int{Value: 0},
				&String{Value: "d"},
				&String{Value: "e"},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}, &String{Value: "d"}, &String{Value: "e"}}},
		},
		{name: "append with big index negative count",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 5},
				&Int{Value: -1},
				&String{Value: "d"},
				&String{Value: "e"},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}, &String{Value: "d"}, &String{Value: "e"}}},
		},
		{name: "delete all with positive count",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 0},
				&Int{Value: 3},
			},
			deleted: &Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
			Array:   &Array{Value: []Object{}},
		},
		{name: "delete all with big count",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 0},
				&Int{Value: 5},
			},
			deleted: &Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
			Array:   &Array{Value: []Object{}},
		},
		{name: "nothing0",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 0},
				&Int{Value: -1},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
		},
		{name: "nothing1",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 1},
				&Int{Value: -1},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
		},
		{name: "nothing2",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
			},
			deleted: &Array{Value: []Object{}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
		},
		{name: "pop without count",
			args: []Object{
				&Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}, &Int{Value: 2}}},
				&Int{Value: 2},
			},
			deleted: &Array{Value: []Object{&Int{Value: 2}}},
			Array:   &Array{Value: []Object{&Int{Value: 0}, &Int{Value: 1}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinSplice(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinSplice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.deleted) {
				t.Errorf("builtinSplice() = %v, want %v", got, tt.deleted)
			}
			if tt.wantErr && tt.wantedErr.Error() != err.Error() {
				t.Errorf("builtinSplice() error = %v, wantedErr %v", err, tt.wantedErr)
			}
			if tt.Array != nil && !reflect.DeepEqual(tt.Array, tt.args[0]) {
				t.Errorf("builtinSplice() arrays are not equal expected %s, got %s", tt.Array, tt.args[0].(*Array))
			}
		})
	}
}
