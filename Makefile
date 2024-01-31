build:
	go build -o cli .
read:
	./cli --read true --file ./writed.txt 
write:
	./cli --write true