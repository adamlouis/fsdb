
run:
	go run main.go index -v --hash

clean:
	rm -f ./*.db
	rm -f ./**/*.db
	rm -f ./*.db-journal
	rm -f ./**/*.db-journal
	rm -f fsdb

build:
	go build -o fsdb main.go

test:
	go run main.go index -v --hash --root=./test
