build:
	go build -o cli .
read:
	make build &&  ./cli --file ./output.png --read true 
write: 
	make build &&  ./cli --file ./output.png --write true 
fresh:
	./cli --refresh true

test:
	go test . -v