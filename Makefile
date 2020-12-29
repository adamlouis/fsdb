
run:
	go run main.go

run-index:
	go run main.go index -v

clean:
	rm -f ./*.db
	rm -f ./**/*.db
	rm -f ./*.db-journal
	rm -f ./**/*.db-journal
	rm -f fsdb

build:
	go build -o fsdb main.go

