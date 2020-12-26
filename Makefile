build:
	@go build -o syl
install: 
	@go build -o "$(GOPATH)/bin/syl"