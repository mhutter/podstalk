NAME = podstalk
IMAGE = mhutter/$(NAME)
TAG = 2019

REPO = $(IMAGE):$(TAG)
TAG_AMD64 = $(REPO)-amd64
TAG_ARM64 = $(REPO)-arm64

test:
	go test -v -race -cover ./...

dockerhub: images push manifest

manifest:
	docker manifest create --amend \
		$(REPO) \
		$(REPO)-arm64 $(REPO)-amd64
	docker manifest annotate \
		--os linux \
		--arch arm64 \
		$(REPO) $(REPO)-arm64
	docker manifest push $(REPO)

images: image-arm64 image-amd64
image-arm64: Dockerfile.arm64
	docker build -t $(TAG_ARM64) -f Dockerfile.arm64 .
image-amd64: Dockerfile
	docker build -t $(TAG_AMD64) .

push: push-amd64 push-arm64
push-amd64:
	docker push $(TAG_AMD64)
push-arm64:
	docker push $(TAG_ARM64)

podstalk: *.go
	go build -o $(NAME) ./cmd/$(NAME)

clean:
	rm -rf $(NAME) gin-bin client/build/

Dockerfile.arm64: Dockerfile
	sed -e 's|amd64|arm64|' Dockerfile > Dockerfile.arm64

dev-server:
	gin --build cmd/podstalk --immediate --port 8080
dev-client:
	cd client && yarn start

client:
	cd client && yarn build

.PHONY: test image clean push dev-server dev-client client
