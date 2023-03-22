package internal

import (
	"errors"
	"log"
	"net/url"
	"reflect"
	"testing"

	"github.com/tpmanc/myhttp/internal/hash"
	"github.com/tpmanc/myhttp/internal/http"
	"github.com/tpmanc/myhttp/internal/models"
)

func TestApp_Process(t *testing.T) {
	type fields struct {
		logger      *log.Logger
		httpClient  httpClient
		hashService hashService
		limiter     chan struct{}
	}
	type args struct {
		urls []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []models.Page
	}{
		{
			name: "without error",
			fields: fields{
				logger:      log.Default(),
				hashService: hash.Mock{},
				httpClient: http.Mock{
					DoRequestFn: func(u *url.URL) ([]byte, error) {
						return []byte("response"), nil
					},
				},
				limiter: make(chan struct{}, 1),
			},
			args: args{
				urls: []string{"http://url.com", "url2.com", "127.0.0.1:80"},
			},
			want: []models.Page{
				{URL: "http://url.com", Hash: "hash string"},
				{URL: "http://url2.com", Hash: "hash string"},
			},
		},

		{
			name: "with error: processURL()",
			fields: fields{
				logger:      log.Default(),
				hashService: hash.Mock{},
				httpClient: http.Mock{
					DoRequestFn: func(u *url.URL) ([]byte, error) {
						return nil, errors.New("test")
					},
				},
				limiter: make(chan struct{}, 3),
			},
			args: args{
				urls: []string{"http://url.com", "url2.com", "127.0.0.1:80"},
			},
			want: []models.Page{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				logger:      tt.fields.logger,
				httpClient:  tt.fields.httpClient,
				hashService: tt.fields.hashService,
				limiter:     tt.fields.limiter,
			}
			if got := a.Process(tt.args.urls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_processURL(t *testing.T) {
	type fields struct {
		logger      *log.Logger
		httpClient  httpClient
		hashService hashService
		limiter     chan struct{}
	}
	type args struct {
		u *url.URL
	}

	u, _ := url.Parse("http://google.com")

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "without error",
			fields: fields{
				hashService: hash.Mock{},
				httpClient: http.Mock{
					DoRequestFn: func(u *url.URL) ([]byte, error) {
						return []byte("response"), nil
					},
				},
			},
			args:    args{u: u},
			want:    "hash string",
			wantErr: false,
		},

		{
			name: "with error: httpClient.DoRequest()",
			fields: fields{
				hashService: hash.Mock{},
				httpClient: http.Mock{
					DoRequestFn: func(u *url.URL) ([]byte, error) {
						return nil, errors.New("error")
					},
				},
			},
			args:    args{u: u},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				logger:      tt.fields.logger,
				httpClient:  tt.fields.httpClient,
				hashService: tt.fields.hashService,
				limiter:     tt.fields.limiter,
			}
			got, err := a.processURL(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("processURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("processURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewApp(t *testing.T) {
	type args struct {
		logger        *log.Logger
		parallelCount int
		hashService   hashService
		httpClient    httpClient
	}

	logger := log.Default()
	hashServiceMock := hash.Mock{}
	httpClientMock := http.Mock{}

	tests := []struct {
		name    string
		args    args
		want    App
		wantErr bool
	}{
		{
			name: "with error: logger == nil",
			args: args{
				logger:        nil,
				parallelCount: 1,
				hashService:   hashServiceMock,
				httpClient:    httpClientMock,
			},
			want:    App{},
			wantErr: true,
		},

		{
			name: "with error: parallelCount <= 0",
			args: args{
				logger:        logger,
				parallelCount: 0,
				hashService:   hashServiceMock,
				httpClient:    httpClientMock,
			},
			want:    App{},
			wantErr: true,
		},

		{
			name: "with error: hashService == nil",
			args: args{
				logger:        logger,
				parallelCount: 1,
				hashService:   nil,
				httpClient:    httpClientMock,
			},
			want:    App{},
			wantErr: true,
		},

		{
			name: "with error: httpClient == nil",
			args: args{
				logger:        logger,
				parallelCount: 1,
				hashService:   hashServiceMock,
				httpClient:    nil,
			},
			want:    App{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewApp(tt.args.logger, tt.args.parallelCount, tt.args.hashService, tt.args.httpClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareURL(t *testing.T) {
	type args struct {
		urlStr string
	}

	u, _ := url.Parse("http://google.com")

	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "without error",
			args: args{
				urlStr: "http://google.com",
			},
			want:    u,
			wantErr: false,
		},

		{
			name: "without error: add scheme",
			args: args{
				urlStr: "google.com",
			},
			want:    u,
			wantErr: false,
		},

		{
			name: "with error: invalid url",
			args: args{
				urlStr: "127.0.0.1:80",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareURL(tt.args.urlStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil || tt.want != nil {
				if got.String() != tt.want.String() {
					t.Errorf("prepareURL() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
