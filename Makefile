.PHONY: run
run:
	go build -o build cmd/app/main.go
	./build
