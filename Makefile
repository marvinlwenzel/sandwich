RND=mytmpcontainernamethatisveryrandomandunlikelytoalreadyexist
IMG_NAME=ghcr.io/marvinlwenzel/sandwich
IMG_TAG=dev

container-fs: image
	podman create --name="$(RND)" $(IMG_NAME):$(IMG_TAG)
	podman export $(RND) > container.tar
	podman rm $(RND)
	mkdir container
	tar -xf container.tar -C container/

image: clean sandwich.stripped
	podman build -t $(IMG_NAME):$(IMG_TAG) .

sandwich.stripped:
	go mod tidy
	CGO_ENABLED=0 go build -o "sandwich.stripped" -ldflags="-s -w"

sandwich.full:
	go mod tidy
	CGO_ENABLED=0 go build -o "sandwich.full"


clean:
	rm -f sandwich.stripped
	rm -f sandwich.full
	rm -f container.tar
	rm -rf container