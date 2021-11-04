package vigo

import (
	"reflect"
	"testing"
)

func TestSplitUint16(t *testing.T) {
	type args struct {
		val uint16
	}
	tests := []struct {
		name  string
		args  args
		want  uint8
		want1 uint8
	}{
		{name: "correct return", args: args{val: 0x1234}, want: 0x12, want1: 0x34},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitUint16(tt.args.val)
			if got != tt.want {
				t.Errorf("SplitUint16() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitUint16() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSplitUint8(t *testing.T) {
	type args struct {
		val uint8
	}
	tests := []struct {
		name  string
		args  args
		want  uint8
		want1 uint8
	}{
		{name: "correct return", args: args{val: 0x12}, want: 0x1, want1: 0x2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitUint8(tt.args.val)
			if got != tt.want {
				t.Errorf("SplitUint8() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitUint8() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_combine3u4(t *testing.T) {
	type args struct {
		i1 uint8
		i2 uint8
		i3 uint8
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{name: "correct", args: args{i1: 0xA, i2: 0xB, i3: 0xC}, want: 0xABC},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combine3u4(tt.args.i1, tt.args.i2, tt.args.i3); got != tt.want {
				t.Errorf("combine3u4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_u8ToBits(t *testing.T) {
	type args struct {
		v uint8
	}
	tests := []struct {
		name string
		args args
		want [8]bool
	}{
		{name: "test 1", args: args{v: 0b1010_0101}, want: [8]bool{true, false, true, false, false, true, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u8ToBits(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("u8ToBits() = %v, want %v", got, tt.want)
			}
		})
	}
}
