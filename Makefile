run-go: build-go
	./main.exe

build-go:
	go build -o main.exe go/main.go

clean:
	rm a.out
