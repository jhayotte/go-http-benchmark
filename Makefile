.PHONY: build up benchmark report clean lint format

GO_LINTER := golangci-lint
GO_FORMATTER := gofmt

.DEFAULT_GOAL := build

build:
	docker-compose build

up:
	docker-compose up -d

benchmark:
	mkdir -p results
	go run loadtester/main.go

report:
	cat results/*_report.txt > results/summary_report.txt
	cat results/summary_report.txt

clean:
	docker-compose down
	rm -rf results

lint:
	$(GO_LINTER) run ./...

format:
	find . -type f -name '*.go' ! -path './vendor/*' -exec $(GO_FORMATTER) -s -w {} +

check: lint format
