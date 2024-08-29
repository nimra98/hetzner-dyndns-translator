.DEFAULT_GOAL := default

# If makefile is called using sudo or as root, ensure that this user is also logged into dockerhub

IMAGE ?= nimra98/hetzner-dyndns-translator
LATEST ?= false

ifndef VERSION
$(error VERSION is not set)
endif

# Aktualisiere die Versionsnummer in main.go
.PHONY: update-version
update-version:
	sed -i 's/const VERSION = ".*"/const VERSION = "$(VERSION)"/g' main.go

.PHONY: build # Build the container image
build: update-version
	@docker buildx create --use --name=crossplat --node=crossplat && \
	if [ "$(LATEST)" = "true" ]; then \
		docker buildx build \
		--output "type=docker,push=false" \
		--tag $(IMAGE):$(VERSION) \
		--tag $(IMAGE):latest \
		. ; \
	else \
		docker buildx build \
		--output "type=docker,push=false" \
		--tag $(IMAGE):$(VERSION) \
		. ; \
	fi	
	

.PHONY: release # Push the image to the remote registry
release: update-version
	@docker buildx create --use --name=crossplat --node=crossplat && \
	if [ "$(LATEST)" = "true" ]; then \
		docker buildx build \
		--platform linux/386,linux/amd64,linux/arm/v7,linux/arm64 \
		--push \
		--tag $(IMAGE):$(VERSION) \
		--tag $(IMAGE):latest \
		. ; \
	else \
		docker buildx build \
		--platform linux/386,linux/amd64,linux/arm/v7,linux/arm64 \
		--output "type=image,push=true" \
		--tag $(IMAGE):$(VERSION) \
		. ; \
	fi
