.DEFAULT_GOAL := default

# If makefile is called using sudo or as root, ensure that this user is also logged into dockerhub

IMAGE ?= nimra98/hetzner-dyndns-translator
LATEST ?= false

ifndef VERSION
$(error VERSION is not set)
endif

.PHONY: build # Build the container image
build:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	if [ "$(LATEST)" = "true" ]; then \
		docker buildx build \
		--build-arg VERSION=$(VERSION) \
		--output "type=docker,push=false" \
		--tag $(IMAGE):$(VERSION) \
		--tag $(IMAGE):latest \
		. ; \
	else \
		docker buildx build \
		--build-arg VERSION=$(VERSION) \
		--output "type=docker,push=false" \
		--tag $(IMAGE):$(VERSION) \
		. ; \
	fi
	

.PHONY: release # Push the image to the remote registry
release:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	if [ "$(LATEST)" = "true" ]; then \
		docker buildx build \
		--build-arg VERSION=$(VERSION) \
		--platform linux/386,linux/amd64,linux/arm/v7,linux/arm64 \
		--push \
		--tag $(IMAGE):$(VERSION) \
		--tag $(IMAGE):latest \
		. ; \
	else \
		docker buildx build \
		--build-arg VERSION=$(VERSION) \
		--platform linux/386,linux/amd64,linux/arm/v7,linux/arm64 \
		--output "type=image,push=true" \
		--tag $(IMAGE):$(VERSION) \
		. ; \
	fi
