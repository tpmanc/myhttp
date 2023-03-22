package hash

import (
	"reflect"
	"testing"
)

func TestMD5Hash_generate(t1 *testing.T) {
	type args struct {
		msg []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "generate hash",
			args: args{
				msg: []byte("test string"),
			},
			want: "6f8db599de986fab7a21625b7916589c",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := MD5Hash{}
			if got := t.generate(tt.args.msg); got != tt.want {
				t1.Errorf("generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMD5(t *testing.T) {
	tests := []struct {
		name string
		want MD5Hash
	}{
		{
			name: "create new MD5 generator",
			want: MD5Hash{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMD5(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMD5() = %v, want %v", got, tt.want)
			}
		})
	}
}
