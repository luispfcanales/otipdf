pull:
	git pull
build:
	go build main.go
run: pull build
	sudo PORT=80 ./main
