.PHONY: test
test: 
	@go test ./...

.PHONY: bench
bench: 
	# run all benchmark tests 5 times, skip unit tests
	@go test ./... -bench=. -benchtime=100x -count=5 -run=^#