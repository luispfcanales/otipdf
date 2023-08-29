build:
	go build main.go
run: build
	sudo PORT=80 ./main
