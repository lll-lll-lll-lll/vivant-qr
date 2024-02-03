build:
	go build -o cli .
read:
	make build &&  ./cli --read true
write:
	make build &&  ./cli --write true
fresh:
	./cli --refresh true

test:
	go test . -v