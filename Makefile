build:
	go build -o cli .
read:
	./cli --read true --file ./writed.txt 
write:
	./cli --write true
fresh:
	./cli --refresh true

test:
	go test . -v