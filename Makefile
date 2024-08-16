image: clean sandwich
	podman build -t ghcr.io/marvinlwenzel/sandwich:dev .

sandwich:
	go mod tidy
	go build

clean:
	rm -f sandwich