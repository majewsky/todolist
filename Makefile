todolist: *.go
	@go build -o $@

assets.go: static/*
	@go generate
