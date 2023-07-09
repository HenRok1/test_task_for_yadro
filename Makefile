.PHONY: run

run:
	go run cmd/app/main.go files/test.txt

build:
	go build cmd/app/main.go
	./main files/test.txt

clean:
	rm ./main