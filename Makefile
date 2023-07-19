.PHONY: run, dev

run :
	@go run .

build :
	@go build .

dev :
	@RUN_MODE=dev go run .