NAME = podstalk
IMAGE = mhutter/$(NAME)

test:
	go test -v -race -cover ./...

image:
	docker build -t $(IMAGE) .

image-arm64: Dockerfile.arm64
	docker build -t $(IMAGE):arm64 -f Dockerfile.arm64 .

podstalk: *.go
	go build -o $(NAME) ./cmd/$(NAME)

clean:
	rm -f $(NAME) gin-bin

Dockerfile.arm64: Dockerfile
	sed -e 's|amd64|arm64|' Dockerfile > Dockerfile.arm64

.PHONY: test image clean
