package hash

import (
	"reflect"
	"testing"
)

func TestHash_GenerateHash(t *testing.T) {
	type fields struct {
		hashAlgo hashAlgo
	}
	type args struct {
		msg []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "md5 hash",
			fields: fields{
				hashAlgo: NewMD5(),
			},
			args: args{
				msg: []byte("test string"),
			},
			want: "6f8db599de986fab7a21625b7916589c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Hash{
				hashAlgo: tt.fields.hashAlgo,
			}
			if got := h.GenerateHash(tt.args.msg); got != tt.want {
				t.Errorf("GenerateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		hashAlgo hashAlgo
	}

	md5Algorithm := NewMD5()

	tests := []struct {
		name    string
		args    args
		want    Hash
		wantErr bool
	}{
		{
			name: "without error: with MD5 algorithm",
			args: args{
				hashAlgo: md5Algorithm,
			},
			want: Hash{
				hashAlgo: md5Algorithm,
			},
			wantErr: false,
		},

		{
			name: "with error: hash algorithm is nil",
			args: args{
				hashAlgo: nil,
			},
			want:    Hash{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.hashAlgo)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
