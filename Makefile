.PHONY: run, dev

run :
	@go run .

build :
	@go build .

install:
	@go install github.com/nicewook/gopt

dev :
	@RUN_MODE=dev go run .