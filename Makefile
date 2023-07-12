.PHONY: run

run:
	go run cmd/app/main.go files/test2.txt

build:
	go build cmd/app/main.go

clean:
	rm ./main

docker:
	docker build -t test_task .

docker_run:
	docker run -p 8080:8080 test_task
