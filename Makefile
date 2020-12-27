build:
	@go build -i ./cmd/syl
install: 
	@go build -o "$(GOPATH)/bin/syl" -i ./cmd/syl 