.PHONY: test
test: 
	@go test ./...

.PHONY: bench
bench: 
	@go test ./... -bench=. -benchtime=100x -count=5 -run=^#

.PHONY: bench-parser
bench-parser: 
	@go test ./... -bench=BenchmarkParsing -benchtime=100x -count=5 -run=^#

.PHONY: bench-tokenizer
bench-tokenizer: 
	@go test ./... -bench=BenchmarkTokenizer -benchtime=100x -count=5 -run=^#
