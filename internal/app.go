package internal

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/tpmanc/myhttp/internal/models"
)

type hashService interface {
	GenerateHash(msg []byte) string
}

type httpClient interface {
	DoRequest(url *url.URL) ([]byte, error)
}

type App struct {
	logger      *log.Logger
	httpClient  httpClient
	hashService hashService
	limiter     chan struct{}
}

// NewApp returns a new instance of App
func NewApp(logger *log.Logger, parallelCount int, hashService hashService, httpClient httpClient) (App, error) {
	if logger == nil {
		return App{}, errors.New("logger mustn't be nil")
	}
	if parallelCount <= 0 {
		return App{}, errors.New("parallelCount must be grater than 0")
	}
	if hashService == nil {
		return App{}, errors.New("hashService mustn't be nil")
	}
	if httpClient == nil {
		return App{}, errors.New("httpClient mustn't be nil")
	}

	return App{
		logger:      logger,
		hashService: hashService,
		httpClient:  httpClient,
		limiter:     make(chan struct{}, parallelCount),
	}, nil
}

// Stop clears resources
func (a App) Stop() {
	close(a.limiter)
}

// Process handles provided URL slice and returns slice of models.Page
func (a App) Process(urls []string) []models.Page {
	wg := sync.WaitGroup{}
	wg.Add(len(urls))
	res := make([]models.Page, 0, len(urls))

	for _, urlStr := range urls {
		a.limiter <- struct{}{}

		go func(urlStr string) {
			defer func() {
				wg.Done()
				<-a.limiter
			}()

			u, err := prepareURL(urlStr)
			if err != nil {
				a.logger.Println(fmt.Errorf("unable to prepare url \"%s\": %v", urlStr, err))
				return
			}

			hash, err := a.processURL(u)
			if err != nil {
				a.logger.Println(fmt.Errorf("unable to process url %s: %v", urlStr, err))
				return
			}

			res = append(res, models.Page{
				URL:  u.String(),
				Hash: hash,
			})
		}(urlStr)
	}

	wg.Wait()

	return res
}

// processURL makes HTTP request for provided url and returns hash from its body
func (a App) processURL(u *url.URL) (string, error) {
	body, err := a.httpClient.DoRequest(u)
	if err != nil {
		return "", fmt.Errorf("unable to handle request: %v", err)
	}

	hash := a.hashService.GenerateHash(body)

	return hash, nil
}

// prepareURL validates URL and add scheme if doesn't have
func prepareURL(urlStr string) (*url.URL, error) {
	preparedUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %v", err)
	}

	if preparedUrl.Scheme == "" {
		preparedUrl.Scheme = "http"
	}

	return preparedUrl, nil
}
