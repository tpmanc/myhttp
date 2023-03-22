package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_DoRequest(t1 *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		u *url.URL
	}

	serverOk := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte("response"))
	}))
	defer serverOk.Close()
	uOk, _ := url.Parse(serverOk.URL)

	serverErr := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverErr.Close()
	uErr, _ := url.Parse(serverErr.URL)

	uInvalid, _ := url.Parse("http://invalid.com")

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "without error",
			fields: fields{
				client: http.DefaultClient,
			},
			args: args{
				u: uOk,
			},
			want:    []byte("response"),
			wantErr: false,
		},

		{
			name: "with error: http request error",
			fields: fields{
				client: http.DefaultClient,
			},
			args: args{
				u: uInvalid,
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "with error: response code != ok",
			fields: fields{
				client: http.DefaultClient,
			},
			args: args{
				u: uErr,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Client{
				client: tt.fields.client,
			}
			got, err := t.DoRequest(tt.args.u)
			if (err != nil) != tt.wantErr {
				t1.Errorf("DoRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("DoRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		client *http.Client
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{
			name: "with default client",
			args: args{
				client: http.DefaultClient,
			},
			want: Client{
				client: http.DefaultClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
