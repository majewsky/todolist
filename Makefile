todolist: *.go
	@env GOPATH=$(CURDIR)/gopath go build -o $@

assets.go: static/*
	@go generate
