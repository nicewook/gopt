.PHONY: run, dev

run :
# nodemon --exec go run . --signal SIGTERM
	@go run .
dev :
# nodemon --exec go run . --signal SIGTERM
	@RUN_MODE=dev go run .