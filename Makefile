all: test badge

test:
	go test -v -shuffle=on -race -coverprofile cover.out  -covermode=atomic ./...
	go tool cover -func=cover.out
	go tool cover -html=cover.out -o coverage-report.html

badge:
	go install github.com/jpoles1/gopherbadger
	gopherbadger

