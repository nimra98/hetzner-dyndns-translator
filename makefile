VERSION ?= 1.0.0

all: build build-docker

build:
	# Put VERSION in main.go by replacing the placeholder
	sed -i 's/const VERSION = ".*"/const VERSION = "$(VERSION)"/g' main.go
	# Build the binary
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./hetzner-dyndns-translator github.com/nimra98/hetzner-dyndns-translator

build-docker:
	docker build --rm -t nimra98/hetzner-dyndns-translator:latest .
	docker tag nimra98/hetzner-dyndns-translator:latest nimra98/hetzner-dyndns-translator:$(VERSION)

release: build build-docker
	docker push nimra98/hetzner-dyndns-translator:latest
	docker push nimra98/hetzner-dyndns-translator:$(VERSION)