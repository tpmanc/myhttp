build:
	go build -o myhttp

help:
	go run main.go -help

test:
	go test -v ./... -covermode=count -coverprofile=coverage.out

run:
	go run main.go http://www.adjust.com https://google.com \
					adjust.com google.com facebook.com yahoo.com yandex.com twitter.com\
                   reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com

run-single:
	go run main.go -parallel 1 http://www.adjust.com http://google.com facebook.com http://yandex.com http://twitter.com

run-invalid:
	go run main.go invalid http://www.adjust.com http://google.com
