.PHONY: run

run:
	go run cmd/app/main.go files/test2.txt

build:
	go build cmd/app/main.go

clean:
	rm ./main