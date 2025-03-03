.PHONY: build up benchmark report clean

build:
	docker-compose build

up:
	docker-compose up -d

benchmark:
	mkdir results
	go run loadtester/main.go

report:
	cat results/*_report.txt > results/summary_report.txt
	cat results/summary_report.txt

clean:
	docker-compose down
	rm -rf results
